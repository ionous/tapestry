package files

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const SaveFileExtension = ".tap"

// golang doesn't provide file creation time in a cross-platform way
// so we add the current time to the name ( base 36 )
// future: include a user specified tag ( separated by a hash or something. )
func NameWithTime(prefix, ext string) (ret string) {
	stamp := strconv.FormatInt(time.Now().UnixMilli(), 32)
	return prefix + "-" + stamp + ext
}

// returns the name.ext ( no path )
// assumes files are formatted with NameWithTime
func FindLatestNameWithTime(dir, prefix, ext string) (ret string, err error) {
	if files, e := os.ReadDir(dir); e != nil {
		err = e
	} else {
		// files are in lexical order
		for i := len(files) - 1; i >= 0; i-- {
			file := files[i]
			if n := file.Name(); !file.IsDir() &&
				strings.HasSuffix(n, ext) &&
				strings.HasPrefix(n, prefix) {
				ret = n
				break
			}
		}
	}
	return
}
