package mosaic

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

// a directory of story files
type storyFolder struct {
	cfg  *Config
	path string
}

// String name of the folder.
func (d storyFolder) String() string {
	return string(d.path)
}

// Find the named child resource.
func (d storyFolder) Find(sub string) (ret web.Resource) {
	base := string(d.path)
	path := filepath.Join(base, sub)
	// join cleans the elements; removing .. paths
	// it helps let us make sure we're still in the right spot.
	if strings.HasPrefix(path, base) {
		if info, e := os.Lstat(path); e != nil {
			// we could return an erroring resource if we really wanted i suppose...
			log.Println("ERROR: reading", d, sub, e)
		} else if info.IsDir() {
			ret = storyFolder{d.cfg, path}
		} else {
			ret = storyFile{d.cfg, path}
		}
	}
	return
}

// Get the contents of this resource.
func (d storyFolder) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	if files, e := listDirectory(d.path); e != nil {
		err = e
	} else if b, e := json.Marshal(files); e != nil {
		err = e
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
	return
}

// Post a modification to this resource
func (d storyFolder) Post(context.Context, io.Reader, http.ResponseWriter) error {
	return errutil.New("unsupported post", d)
}

// Put new resource data in our place
func (d storyFolder) Put(context.Context, io.Reader, http.ResponseWriter) error {
	return errutil.New("unsupported put", d)
}

// based on filepath.Walk
func listDirectory(path string) (ret []string, err error) {
	if f, e := os.Open(path); e != nil {
		err = e
	} else {
		defer f.Close()
		if names, e := f.Readdirnames(-1); e != nil {
			err = e
		} else {
			for _, name := range names {
				filename := filepath.Join(path, name)
				if info, e := os.Lstat(filename); e != nil {
					err = e
					break
				} else {
					isDir := info.IsDir()
					if isDir || files.Ext(name).Story() {
						if name[0] != '_' && name[0] != '.' {
							if isDir {
								name = "/" + name
							}
							ret = append(ret, name)
						}
					}
				}
			}
		}
	}
	return
}
