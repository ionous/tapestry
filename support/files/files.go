// grab bag of file utility functions
package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionous/errutil"
)

// read a comma-separated list of files and directories
// for directories, ext ( a list of file extensions ) optionally filters the files.
// fix? maybe filepaths could be turned into an io.fs?
func ReadPaths(filePaths string, recusive bool, exts []string, onFile func(string) error) (err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if srcPath, e := filepath.Abs(filePath); e != nil {
			err = e
		} else if info, e := os.Stat(srcPath); e != nil {
			err = errutil.Append(err, e)
		} else if !info.IsDir() {
			err = errutil.Append(err, readOne(srcPath, onFile))
		} else {
			err = errutil.Append(err, readMany(srcPath, recusive, exts, onFile))
		}
	}
	return
}

// read the complete contents of the passed file
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// exts: optional list of ".ext" to filter.
func readMany(root string, recusive bool, exts []string, onFile func(string) error) error {
	if !strings.HasSuffix(root, "/") {
		root += "/" // for opening symbolic directories
	}
	outErr := filepath.WalkDir(root, func(path string, info fs.DirEntry, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() {
			if len(exts) == 0 || IsValidExtension(path, exts) {
				err = readOne(path, onFile)
			}
		} else if !recusive && path != root {
			return filepath.SkipDir
		}
		return
	})
	return outErr
}

func readOne(path string, onFile func(string) error) (err error) {
	if e := onFile(path); e != nil {
		err = errutil.New("error reading", path, e)
	}
	return
}

// is the extension of the passed path one of the specified extensions?
func IsValidExtension(path string, exts []string) (okay bool) {
	ext := filepath.Ext(path)
	for _, x := range exts {
		if ext == x {
			okay = true
			break
		}
	}
	return
}