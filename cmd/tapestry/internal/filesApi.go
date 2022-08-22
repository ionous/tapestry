package tap

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

func FilesApi(cfg *Config) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			// by adding a trailing slash, walk will follow a symlink.
			path := cfg.PathTo("stories") + string(filepath.Separator)
			switch name {
			case "blocks":
				// generates blockly files out of .if files
				where := storyFolder{cfg, path}
				ret = blocksRoot{blocksFolder{where}}
			case "stories":
				// serves raw .if files ( compact json format )
				where := storyFolder{cfg, path}
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
				var dst story.Story // composer hands back old format stories ( because that's what we give it )
				if e := din.Decode(&dst, tapestry.Registry(), el.Story); e != nil {
					err = e
				} else {
					file := story.StoryFile{
						StoryLines: dst.Reformat(),
					}
					if data, e := story.Encode(&file); e != nil {
						err = errutil.Append(err, e)
					} else {
						err = writeOut(at, data)
					}
				}
			}
		}
	}
	if err != nil {
		log.Println("ERROR", err)
	}
	return
}
