package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

// endpoint containing a local story file.
// handles getting the contents, and a subaction to post a check of the current contents.
//   - /stories/<path>/<to>/<file.ext>: get or put individual story files
//   - /stories/<path>/<to>/<file.ext>/check: post to test story tests
type storyFile struct {
	cfg  *Config
	path string
}

// String name of the file.
func (sf storyFile) String() string {
	return sf.path
}

// Find actions for individual files
// check tests a file.
func (sf storyFile) Find(sub string) (ret web.Resource) {
	switch sub {
	case "check":
		ret = &web.Wrapper{
			Posts: func(ctx context.Context, in io.Reader, out http.ResponseWriter) (err error) {
				if e := tempTest(ctx, sf.cfg, sf.path, in); e != nil {
					err = e
				}
				return
			},
		}
	}
	return
}

// files are stored in compact format
// we check that the file is valid ( by loading it ) before returning it.
func (sf storyFile) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	var msg map[string]any
	if e := files.FormattedLoad(sf.path, &msg); e != nil {
		err = e
	} else {
		// verify the story is valid by loading it.
		var file story.StoryFile
		if e := story.Decode(&file, msg); e != nil {
			err = e
		} else {
			w.Header().Set("Content-Type", "application/json")
			js := json.NewEncoder(w)
			err = js.Encode(msg)
		}
	}
	return
}

// files dont support posting; returns error
func (sf storyFile) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported post", sf)
}

// story files dont support putting; returns error
// ( mosaic puts block files )
func (sf storyFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported put", sf)
}
