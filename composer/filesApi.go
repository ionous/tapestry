package composer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/web"
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
		Path  string          `json:"path"`
		Story json.RawMessage `json:"story"`
	}
	dec := json.NewDecoder(r)
	if e := dec.Decode(&els); err != nil {
		err = e
	} else {
		root := d.String()
		for _, el := range els {
			if at := filepath.Join(root, el.Path); !strings.HasPrefix(at, root) {
				e := errutil.New("cant save to", at)
				err = errutil.Append(err, e)
			} else {
				var dst story.Story
				if e := din.Decode(&dst, iffy.Registry(), el.Story); e != nil {
					err = e
				} else if data, e := cout.Encode(&dst, story.CompactEncoder); e != nil {
					err = errutil.Append(err, e)
				} else {
					err = writeOut(at, data)
				}
			}
		}
	}
	if err != nil {
		log.Println("ERROR", err)
	}
	return
}
