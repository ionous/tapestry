package cmdweave

import (
	"database/sql"
	"log"
	"os"
	"path"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

// lets do this in the dumbest of ways for now.
func WeavePath(srcPath, outFile string) (err error) {
	var cat eph.Catalog // fix: capture "Dilemmas" and LogWarning?
	k := imp.NewImporter(cat.Weaver())
	if e := k.BeginDomain("tapestry", nil); e != nil {
		err = e
	} else if e := addDefaultKinds(k); e != nil {
		err = e
	} else if e := importStoryFiles(k, srcPath); e != nil {
		err = e
	} else if len(cat.Errors) > 0 {
		err = errutil.New(cat.Errors)
	} else if e := assembleCat(cat.Weaver(), outFile); e != nil {
		err = e
	}
	return
}

func addDefaultKinds(n assert.Assertions) (err error) {
	for _, k := range kindsOf.DefaultKinds {
		if e := n.AssertAncestor(k.String(), k.Parent().String()); e != nil {
			err = e
			break
		}
	}
	return
}

func assembleCat(cat *eph.Catalog, outFile string) (err error) {
	if outFile, e := filepath.Abs(outFile); e != nil {
		err = e
	} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
		err = errutil.New("couldn't clean output file", outFile, e)
	} else {
		// 0755 -> readable by all but only writable by the user
		// 0700 -> read/writable by user
		// 0777 -> ModePerm ... read/writable by all
		os.MkdirAll(path.Dir(outFile), os.ModePerm)
		if db, e := sql.Open(tables.DefaultDriver, outFile); e != nil {
			err = errutil.New("couldn't create output file", outFile, e)
		} else {
			defer db.Close()
			err = Weave(cat, db)
		}
	}
	return
}

// read a comma-separated list of files and directories
func importStoryFiles(k *imp.Importer, srcPath string) (err error) {
	recurse := true
	if e := files.ReadPaths(srcPath, recurse,
		[]string{CompactExt, DetailedExt}, func(p string) error {
			return readOne(k, p)
		}); e != nil {
		err = errutil.New("couldn't read file", srcPath, e)
	}
	return
}

func readOne(k *imp.Importer, path string) (err error) {
	log.Println("reading", path)
	if b, e := files.ReadFile(path); e != nil {
		err = e
	} else if script, e := decodeStory(path, b); e != nil {
		err = errutil.New("couldn't decode", path, "b/c", e)
	} else if e := story.ImportStory(k, path, script); e != nil {
		err = errutil.New("couldn't import", path, "b/c", e)
	}
	return
}

func decodeStory(path string, b []byte) (ret *story.StoryFile, err error) {
	switch ext := filepath.Ext(path); ext {
	case CompactExt:
		var curr story.StoryFile
		if e := story.Decode(&curr, b, tapestry.AllSignatures); e != nil {
			err = e
		} else {
			ret = &curr
		}
	case DetailedExt:
		var curr story.StoryFile
		if e := din.Decode(&curr, tapestry.Registry(), b); e != nil {
			err = e
		} else {
			ret = &curr
		}
	default:
		err = errutil.Fmt("unknown file type %q", ext)
	}
	return
}
