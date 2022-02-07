package composer

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"git.sr.ht/~ionous/tapestry/web"
)

// Compose starts the composer server, this function doesnt return.
func Compose(cfg *Config) {
	// redirect from root.
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/compose/index.html", http.StatusMovedPermanently)
	})

	// vue app
	http.Handle("/compose/", http.StripPrefix("/compose/", http.FileServer(http.Dir("./www"))))

	// if specs
	println(cfg.PathTo("ifspec"))
	http.Handle("/ifspec/", http.StripPrefix("/ifspec/", http.FileServer(specFsSystem{http.Dir(cfg.PathTo("ifspec"))})))

	// story files
	http.HandleFunc("/stories/", web.HandleResourceWithContext(FilesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	}))

	log.Println("Composer using", cfg.Root)
	log.Println("Listening on port", cfg.PortString(), "...")
	if e := http.ListenAndServe(cfg.PortString(), nil); e != nil {
		log.Fatal(e)
	}
}

// starts the blockly editor, this function doesnt return.
func Mosaic(cfg *Config) {
	// configure server root
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/mosaic/index.html", http.StatusMovedPermanently)
	})

	// in dev mode, we reflect all calls to /mosaic/ to port 3000
	if m := cfg.Mosaic; !strings.HasPrefix(m, "http") {
		// maybe a directory?
		panic("not implemented")
		http.Handle("/mosaic/", http.StripPrefix("/mosaic/", http.FileServer(http.Dir("./www"))))
	} else if u, e := url.Parse(m); e != nil {
		panic(e)
	} else {
		http.Handle("/mosaic/", httputil.NewSingleHostReverseProxy(u))
	}

	// blockly blocks ( from .if )
	http.HandleFunc("/blocks/", web.HandleResourceWithContext(FilesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	}))

	// blockly shape files ( from .ifspecs )
	http.Handle("/shapes/", http.StripPrefix("/shapes/", web.HandleResourceWithContext(ShapesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	})))

	// http.HandleFunc("/stories/", web.HandleResourceWithContext(FilesApi(cfg), func(ctx context.Context) context.Context {
	// 	return context.WithValue(ctx, configKey, cfg)
	// }))

	log.Println("Composer using", cfg.Root)
	log.Println("Listening on port", cfg.PortString(), "...")
	if e := http.ListenAndServe(cfg.PortString(), nil); e != nil {
		log.Fatal(e)
	}
}

type key int

// passed to the http context to store a pointer to the composer Config.
var configKey key

// containsDotFile reports whether name contains a path element starting with a period.
// http.FileSystem guarantees the name has forward slash delimiting.
func containsDotFile(name string) (okay bool) {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			okay = true
			break
		}
	}
	return
}

// wrap http.File to filter for .ifspec files
type specFs struct{ http.File }

// Readdir filters files to only report .ifspec.
func (f specFs) Readdir(n int) (ret []fs.FileInfo, err error) {
	if files, e := f.File.Readdir(n); e != nil {
		err = e
	} else {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".ifspec") {
				ret = append(ret, file)
			}
		}
	}
	return
}

// an http.FileSystem to limit the files served
type specFsSystem struct{ http.FileSystem }

// serves a 403 permission error when has a requested file or dir starts with a dot.
func (fsys specFsSystem) Open(name string) (ret http.File, err error) {
	if containsDotFile(name) {
		err = fs.ErrPermission // 403
	} else if file, e := fsys.FileSystem.Open(name); e != nil {
		err = e
	} else {
		ret = specFs{file}
	}
	return
}
