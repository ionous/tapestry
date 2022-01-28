package files

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionous/errutil"
)

// read a comma-separated list of files and directories
// fix? the comma separation isnt used much -- maybe this could just be turned into an io/fs ?
func ReadPaths(filePaths string, exts []string, onFile func(string) error) (err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if srcPath, e := filepath.Abs(filePath); e != nil {
			err = e
		} else if info, e := os.Stat(srcPath); e != nil {
			err = errutil.Append(err, e)
		} else {
			which := readOne
			if info.IsDir() {
				which = readMany
			}
			if e := which(srcPath, exts, info, onFile); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

// read the complete contents of the passed file
// ( similar to fs.ReadFile(src, path) )
func ReadFile(path string) (ret []byte, err error) {
	if fp, e := os.Open(path); e != nil {
		err = e
	} else {
		defer fp.Close()
		if b, e := io.ReadAll(fp); e != nil {
			err = e
		} else {
			ret = b
		}
	}
	return
}

// exts: optional list of ".ext" to filter.
func readMany(path string, exts []string, _ os.FileInfo, onFile func(string) error) error {
	if !strings.HasSuffix(path, "/") {
		path += "/" // for opening symbolic directories
	}
	return filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() {
			err = readOne(path, exts, info, onFile)
		}
		return
	})
}

func readOne(path string, exts []string, info os.FileInfo, onFile func(string) error) (err error) {
	if ext := filepath.Ext(path); len(exts) == 0 || contains(ext, exts) {
		if e := onFile(path); e != nil {
			err = errutil.New("error reading", path, e)
		}
	}
	return
}

func contains(ext string, exts []string) (okay bool) {
	for _, x := range exts {
		if ext == x {
			okay = true
			break
		}
	}
	return
}
