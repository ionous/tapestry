package composer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/iffy/web"
	"github.com/iancoleman/orderedmap"
	"github.com/ionous/errutil"
)

func FilesApi(cfg *Config) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			switch name {
			case "stories":
				// by adding a trailing slash, walk will follow a symlink.
				where := storyFolder(filepath.Join(cfg.Root, "stories") + "/")
				ret = rootFolder{where}
			}
			return
		},
	}
}

type rootFolder struct {
	storyFolder
}

func (d rootFolder) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	var els []struct {
		Path  string                `json:"path"`
		Story orderedmap.OrderedMap `json:"story"`
	}
	dec := json.NewDecoder(r)
	if e := dec.Decode(&els); err != nil {
		err = e
	} else {
		// fix: return status of some sort?
		root := d.String()
		for _, el := range els {
			where := filepath.Join(root, el.Path)
			if !strings.HasPrefix(where, root) {
				e := errutil.New("couldnt save", where)
				err = errutil.Append(err, e)
			} else if e := saveBytes(where, el.Story); e != nil {
				e := errutil.New("couldnt save", where)
				err = errutil.Append(err, e)
			} else {
				log.Println("saved", where)
			}
		}
	}
	if err != nil {
		log.Println("ERROR", err)
	}
	return
}

func saveBytes(where string, story orderedmap.OrderedMap) (err error) {
	if f, e := os.Create(where); e != nil {
		err = e
	} else {
		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		err = enc.Encode(story)
	}
	return
}
