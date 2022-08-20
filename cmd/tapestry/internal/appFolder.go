package tap

import (
	"os"
	"path/filepath"

	"github.com/ionous/errutil"
)

const DefaultFolder int = 8080

// implements the flag.Value interface for locating the author's stories.
type Folder string

func (f Folder) String() (ret string) {
	return string(f)
}

func (f *Folder) Set(s string) (err error) {
	*f = Folder(s) // validated on get.
	return
}

// turn the requested directory into the real thing
// if nothing specified, try the user's documents folder
func (f Folder) GetFolder() (ret string, err error) {
	dir := f.String()
	if len(dir) > 0 {
		dir, err = filepath.Abs(dir)
	} else if home, e := os.UserHomeDir(); e != nil {
		err = e
	} else {
		dir = filepath.Join(home, "Documents", "Tapestry")
	}
	if err == nil {
		if fi, e := os.Stat(dir); e != nil {
			err = errutil.New(dir, "is not a directory", e)
		} else if mode := fi.Mode(); !mode.IsDir() {
			err = errutil.New(dir, "is not a directory")
		} else {
			ret = dir
		}
	}
	return
}
