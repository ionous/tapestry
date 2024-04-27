package files

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const SaveFileExtension = ".tapestry"

// golang doesn't provide file creation time in a cross-platform way
// so we add the current time to the name ( base 36 )
// future: could separate user file name with hashes or something
func NameWithTime(prefix, ext string) (ret string) {
	stamp := strconv.FormatInt(time.Now().UnixMilli(), 32)
	return prefix + "-" + stamp + ext
}

func FindLatest(dir, prefix, ext string) (ret string, err error) {
	if files, e := os.ReadDir(dir); e != nil {
		err = e
	} else {
		// files are in lexical order
		for i := len(files) - 1; i >= 0; i-- {
			file := files[i]
			if n := file.Name(); !file.IsDir() &&
				strings.HasSuffix(n, ext) &&
				strings.HasPrefix(n, prefix) {
				ret = filepath.Join(dir, n)
				break
			}
		}
	}
	return
}

// start := time.Now().Format(time.DateTime))
// "cloak 2006-01-02 15:04:05"
// "cloak-of-darkness-662BEDA4"
// s16 := strconv.FormatUint(v, 16)
// strconv.FormatInt(time.Now().UnixMilli(), 36)
