package composer

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"strings"

	"git.sr.ht/~ionous/tapestry/web"
)

// Compose starts the composer server, this function doesnt return.
func Compose(cfg *web.Config, port int) {
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

	log.Println("Composer using", cfg.PathTo())
	log.Println("Listening on port", port, "...")
	if e := http.ListenAndServe(web.Endpoint(port), nil); e != nil {
		log.Fatal(e)
	}
}

func HandleBackend(cfg *web.Config, name string) {
	// in dev mode, we reflect all unhandled calls to port 3000
	if u, e := url.Parse(web.Endpoint(3000, "localhost")); e != nil {
		panic(e)
	} else {
		p := httputil.NewSingleHostReverseProxy(u)
		//
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if len(r.Method) == 0 || r.Method == "GET" {
				p.ServeHTTP(w, r)
			} else if r.Method == "POST" && r.URL.Path == "/" {
				// read a request from the client
				// see Mosaic.vue onPlay which sends {play:true}
				var cmd cmdFromClient
				dec := json.NewDecoder(r.Body)
				if e := dec.Decode(&cmd); e != nil {
					log.Println(e)
				} else if cmd.Play {
					cmd := cfg.Cmd("serve")
					args := []string{
						"-in", cfg.PathTo("stories", "shared"),
						"-out", cfg.Scratch("play.db"),
						"-open",
						// check
					}
					log.Println("playing", cmd, args)
					go func() {
						cmd := exec.Command(cmd, args...)
						if e := cmd.Run(); e != nil {
							log.Println(e)
						}
					}()
				}
			}
		})

		// tbd: basically i want everything but not any other apps
		// http.Handle("/", p)
		// http.Handle("/mosaic/", p)
		// http.Handle("/index.css", p)
		// http.Handle("/favicon.ico", p)
		// http.Handle("/assets/", p)
		// http.Handle("/lib/", p)
		// http.Handle("/node_modules/", p)
		// http.Handle("/@vite/", p)
		// http.Handle("/@id/", p)
	}
}

// starts the blockly editor, this function doesnt return.
func RunMosaic(cfg *web.Config, port int) {

	// configure server root
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/mosaic/index.html", http.StatusMovedPermanently)
	})

	if prod := cfg.Prod(); len(prod) == 0 {
		HandleBackend(cfg, "mosaic")
	} else {
		//	paths ending with / indicate subtrees ( otherwise it assumes a single resource. )
		// longer paths [ really path patterns ] take precedence.
		// so a bare slash means: all paths not otherwise matched.
		// 'Dir' is a typed string which implements an Fs/Filesystem restricted to the specified tree.
		http.Handle("/", http.FileServer(http.Dir(prod)))
	}

	// blockly blocks ( from .if )
	http.HandleFunc("/blocks/", web.HandleResourceWithContext(FilesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	}))

	// blockly shape files ( from .ifspecs )
	http.Handle("/shapes/", http.StripPrefix("/shapes/", web.HandleResourceWithContext(ShapesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	})))

	// blockly shape files ( from .ifspecs )
	http.Handle("/boxes/", http.StripPrefix("/boxes/", web.HandleResourceWithContext(BoxesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	})))

	// http.HandleFunc("/stories/", web.HandleResourceWithContext(FilesApi(cfg), func(ctx context.Context) context.Context {
	// 	return context.WithValue(ctx, configKey, cfg)
	// }))

	log.Println("Composer using", cfg.PathTo())
	log.Println("Listening on port", port, "...")
	if e := http.ListenAndServe(web.Endpoint(port), nil); e != nil {
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

type cmdFromClient struct {
	Play bool `json:"play"`
}
