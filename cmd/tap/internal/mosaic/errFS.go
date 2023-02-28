package mosaic

import (
	"io/fs"
	"time"
)

// wails serves urls *after* it serves files
// which is the opposite order than makese sense for tapestry.
// [ we know the url endpoints we want, and that's a *much* smaller set than actual assets ]
type ErrFS struct{}

func (f ErrFS) Open(name string) (fs.File, error) {
	// wails opens "." just to see if it can *sigh*
	if name == "." {
		return nil, nil
	} else {
		return nil, fs.ErrNotExist
	}
}

// wails uses Stat as part of NewAssetHandler's search to find an index.html
// it starts by just trying to stat it, and if that fails then uses walks using ReadDir and SubFS to pin that as the root.
// we cut it out as early as we can.
func (f ErrFS) Stat(name string) (ret fs.FileInfo, err error) {
	if name == "index.html" {
		ret = fakeIndex{}
	} else {
		err = fs.ErrNotExist
	}
	return
}

type fakeIndex struct{}

func (fakeIndex) Name() string                 { return "index.html" }
func (fakeIndex) IsDir() bool                  { return false }
func (fakeIndex) Type() fs.FileMode            { return fs.ModePerm }
func (i fakeIndex) Info() (fs.FileInfo, error) { return i, nil }

// extra fileInfo methods
func (fakeIndex) Size() int64        { return 0 }
func (fakeIndex) Mode() fs.FileMode  { return fs.ModePerm }
func (fakeIndex) ModTime() time.Time { return time.Now() }
func (fakeIndex) Sys() any           { return nil }
