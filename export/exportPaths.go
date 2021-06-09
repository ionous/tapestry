package export

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"github.com/ionous/errutil"
)

type File struct {
	Path string
	Data reader.Map
}

// read a comma-separated list of files and directories
func ReadPaths(filePaths string) (ret []File, err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if info, e := os.Stat(filePath); e != nil {
			err = e
		} else {
			if !info.IsDir() {
				if one, e := readOne(filePath); e != nil {
					err = e
				} else {
					ret = append(ret, one)
				}
			} else {
				if many, e := readMany(filePath); e != nil {
					err = e
				} else {
					ret = append(ret, many...)
				}
			}
		}
	}
	return
}

func readMany(path string) (ret []File, err error) {
	if !strings.HasSuffix(path, "/") {
		path += "/" // for opening symbolic directories
	}
	err = filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() && filepath.Ext(path) == ".if" {
			if one, e := readOne(path); e != nil {
				err = errutil.New("error reading", path, e)
			} else {
				ret = append(ret, one)
			}
		}
		return
	})
	return
}

func readOne(filePath string) (ret File, err error) {
	log.Println("reading", filePath)
	if f, e := os.Open(filePath); e != nil {
		err = e
	} else {
		defer f.Close()
		var one reader.Map
		if e := json.NewDecoder(f).Decode(&one); e != nil && e != io.EOF {
			err = e
		} else {
			ret = File{filePath, one}
		}
	}
	return
}
