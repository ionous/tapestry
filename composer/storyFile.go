package composer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn/dout"
	"git.sr.ht/~ionous/iffy/web"
	"github.com/ionous/errutil"
)

// path of a local .if file
type storyFile string

// String name of the file.
func (d storyFile) String() string {
	return string(d)
}

// Find actions for individual files
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

// Get the contents of this resource.
func (d storyFile) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	// files are stored in compact format, and we need to transform them to detailed format for the composer.
	var dst story.Story
	if b, e := readOne(string(d)); e != nil {
		err = e
	} else if e := story.Decode(&dst, b, iffy.AllSignatures); e != nil {
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

// Post a modification to this resource
func (d storyFile) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported post", d)
}

// Put new resource data in our place
func (d storyFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported put", d)
}

// -- saving uses storyFolder ( specifically rootFolder.Put )
// -- so what does this do?
// func (d storyFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
// 	// its okay to use Create because storyFolder.Get() ensures it already exists.
// 	if f, e := os.Create(string(d)); e != nil {
// 		err = e
// 	} else {
// 		defer f.Close()
// 		_, err = io.Copy(f, r)
// 	}
// 	return
// }

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
