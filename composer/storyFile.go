package composer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/dout"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

// path of a local .if file
// has subactions for "check" to test the file.
type storyFile string

// String name of the file.
func (d storyFile) String() string {
	return string(d)
}

// Find actions for individual files
// check tests a file.
func (d storyFile) Find(sub string) (ret web.Resource) {
	switch sub {
	case "check":
		ret = &web.Wrapper{
			Posts: func(ctx context.Context, in io.Reader, out http.ResponseWriter) (err error) {
				if e := tempTest(ctx, d.String(), in); e != nil {
					err = e
				}
				return
			},
		}
	}
	return
}

// files are stored in compact format,
// this transforms them to detailed format for the composer.
func (d storyFile) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	var dst story.Story
	if b, e := readOne(string(d)); e != nil {
		err = e
	} else if e := story.Decode(&dst, b, tapestry.AllSignatures); e != nil {
		err = e
	} else if data, e := dout.Encode(&dst); e != nil {
		err = e
	} else {
		w.Header().Set("Content-Type", "application/json")
		js := json.NewEncoder(w)
		err = js.Encode(data)
	}
	return
}

// files dont support posting; returns error
func (d storyFile) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported post", d)
}

// files dont support putting; returns error
// ( save is handled by putting stories into the folder )
func (d storyFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported put", d)
}

func writeOut(outPath string, data interface{}) (err error) {
	log.Println("writing", outPath)
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		js := json.NewEncoder(fp)
		js.SetIndent("", "  ")
		js.SetEscapeHTML(false)
		err = js.Encode(data)
	}
	return
}

func readOne(filePath string) (ret []byte, err error) {
	log.Println("reading", filePath)
	if fp, e := os.Open(filePath); e != nil {
		err = e
	} else {
		ret, err = io.ReadAll(fp)
		fp.Close()
	}
	return
}
