package composer

import (
	"context"
	"io"
	"log"
	"net/http"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/blockly/bgen"
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
type blocksFile string

// overrides story folder find to return blocks
func (d blocksFolder) Find(sub string) (ret web.Resource) {
	switch res := d.storyFolder.Find(sub).(type) {
	case storyFile:
		ret = blocksFile(res)
	case storyFolder:
		ret = blocksFolder{res}
	default:
		// fix: is there an error resource?
		// and/or why doesnt find return error?
		log.Println("unknown resource at", res, sub)
	}
	return
}

// Find actions for individual files ( none right now )
func (d blocksFile) Find(sub string) (none web.Resource) {
	return
}

// files dont support posting; returns error
func (d blocksFile) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported post", d)
}

// files dont support putting; returns error
func (d blocksFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported put", d)
}

// gets the contents of the story file, transforms it into blocks
func (d blocksFile) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	var dst story.Story
	if b, e := readOne(string(d)); e != nil {
		err = e
	} else if e := story.Decode(&dst, b, tapestry.AllSignatures); e != nil {
		err = e
	} else if str, e := bgen.Convert(&dst); e != nil {
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
