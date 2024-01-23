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

	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

type blocksRoot struct {
	blocksFolder
}

// finds and gets story files
type blocksFolder struct {
	storyFolder
}

// loads story files and translates them into blocks
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
// the mosasic editor sends multiple files at once so this handle save at the root folder level.
// because the same path written twice with the same data should have the same result... this uses put.
func (d blocksFolder) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	var els []struct {
		Path     string          `json:"path"`
		Contents json.RawMessage `json:"contents"`
	}
	dec := json.NewDecoder(r)
	if e := dec.Decode(&els); e != nil {
		err = e
	} else {
		root := d.String()
		for _, el := range els {
			var file story.StoryFile // mosaic hands back blocks
			if at := filepath.Join(root, el.Path); !strings.HasPrefix(at, root) {
				e := errutil.New("cant save to", at)
				err = errutil.Append(err, e)
			} else if e := unblock.Decode(&file, "story_file", story.Registry(), el.Contents); e != nil {
				err = errutil.Append(err, e)
			} else if data, e := story.Encode(&file); e != nil {
				err = errutil.Append(err, e)
			} else if fp, e := os.OpenFile(d.path, os.O_WRONLY|os.O_TRUNC, 0); e != nil {
				err = e // ^ why cant this use create? ( and call FormattedSave?
			} else {
				if e := files.FormattedWrite(fp, data, files.Ext(d.path), true); e != nil {
					err = errutil.Append(err, e)
				}
				fp.Close()
			}
		}
	}
	if err != nil {
		log.Println("Error putting files", err)
	}
	return
}

// Find sub actions for individual files ( none right now )
func (d blocksFile) Find(sub string) (none web.Resource) {
	return
}

// files dont support posting; returns error
func (d blocksFile) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported post", d)
}

// save an individual file
func (d blocksFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	var file story.StoryFile // mosaic hands back blocks
	if raw, e := io.ReadAll(r); e != nil {
		err = e
	} else if e := unblock.Decode(&file, "story_file", story.Registry(), raw); e != nil {
		err = e
	} else if data, e := story.Encode(&file); e != nil {
		err = e
	} else if fp, e := os.OpenFile(d.path, os.O_WRONLY|os.O_TRUNC, 0); e != nil {
		err = e
	} else {
		if e := files.FormattedWrite(fp, data, files.Ext(d.path), true); e != nil {
			err = e
		}
		fp.Close()
	}
	if err != nil {
		log.Println("Error putting file", err)
	}
	return
}

// gets the contents of the story file, transforms it into blocks
func (d blocksFile) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	if b, e := os.ReadFile(d.path); e != nil {
		err = e
	} else {
		var msg map[string]any
		var file story.StoryFile
		if e := json.Unmarshal(b, &msg); e != nil {
			err = e
		} else if e := story.Decode(&file, msg); e != nil {
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
	}
	return
}
