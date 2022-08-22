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
	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

type blocksRoot struct {
	blocksFolder
}

// finds and gets .if files
type blocksFolder struct {
	storyFolder
}

// loads .if files and translates them into blocks
type blocksFile struct {
	cfg  *Config
	path string
}

// overrides story folder find to return blocks
func (d blocksFolder) Find(sub string) (ret web.Resource) {
	switch res := d.storyFolder.Find(sub).(type) {
	case storyFile:
		ret = blocksFile{d.cfg, res.path}
	case storyFolder:
		ret = blocksFolder{res}
	default:
		// fix: is there an error resource?
		// and/or why doesnt find return error?
		log.Println("unknown resource at", res, sub)
	}
	return
}

// save a bunch of files:
// the mosasic editor sends multiple files at once so this handlesave at the root folder level.
// because the same path written twice with the same data should have the same result... this uses put.
func (d blocksFolder) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	var els []struct {
		Path     string          `json:"path"`
		Contents json.RawMessage `json:"contents"`
	}
	dec := json.NewDecoder(r)
	if e := dec.Decode(&els); err != nil {
		err = e
	} else {
		root := d.String()
		for _, el := range els {
			var file story.StoryFile // mosaic hands back blocks
			if at := filepath.Join(root, el.Path); !strings.HasPrefix(at, root) {
				e := errutil.New("cant save to", at)
				err = errutil.Append(err, e)
			} else if e := unblock.Decode(&file, "story_file", tapestry.Registry(), el.Contents); e != nil {
				err = errutil.Append(err, e)
			} else if data, e := story.Encode(&file); e != nil {
				err = errutil.Append(err, e)
			} else if e := writeOut(at, data); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	if err != nil {
		log.Println("ERROR", err)
	}
	return
}

// Find actions for individual files ( none right now )
func (d blocksFile) Find(sub string) (none web.Resource) {
	return
}

// files dont support posting; returns error
// ( blocksFolder however support put )
func (d blocksFile) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported post", d)
}

// files dont support putting; returns error
// ( blocksFolder however support put )
func (d blocksFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported put", d)
}

// gets the contents of the story file, transforms it into blocks
func (d blocksFile) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	var file story.StoryFile
	if b, e := readOne(d.path); e != nil {
		err = e
	} else if e := story.Decode(&file, b, tapestry.AllSignatures); e != nil {
		err = e
	} else if str, e := block.Convert(&file); e != nil {
		err = e
	} else {
		w.Header().Set("Content-Type", "application/json")
		// fix? a small rewrite to js.Builder so that it wraps a custom js.Writer interface
		// which supports byte, rune, string, etc. and then a StreamWriter that implements those for a pure Writer
		// we could stream here.... the normal construction would still use strings.Builder directly
		_, err = io.WriteString(w, str)
	}
	return
}
