package composer

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/web"
)

// Compose starts the composer server, this function doesnt return.
func Compose(cfg *Config) {
	// configure server
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/compose/index.html", http.StatusMovedPermanently)
	})
	http.Handle("/compose/", http.StripPrefix("/compose/", http.FileServer(http.Dir("./www"))))
	println(cfg.PathTo("ifspec"))
	http.Handle("/ifspec/", http.StripPrefix("/ifspec/", http.FileServer(specFsSystem{http.Dir(cfg.PathTo("ifspec"))})))
	http.HandleFunc("/stories/", web.HandleResourceWithContext(FilesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	}))

	log.Println("Composer using", cfg.Root)
	log.Println("Listening on port", strconv.Itoa(cfg.Port)+"...")
	if e := http.ListenAndServe(":3000", nil); e != nil {
		log.Fatal(e)
	}
}

type key int

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
