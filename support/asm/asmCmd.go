package asm

import (
	"database/sql"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/files"
	"github.com/ionous/errutil"
)

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

// lets do this in the dumbest of ways for now.
func AssembleFolder(srcPath, outFile string) (err error) {
	var cat eph.Catalog // fix: capture "Dilemmas" and LogWarning?
	var writeErr error  // fix: this seems less than ideal; maybe writer should return err.
	k := imp.NewImporter(collectEphemera(&cat, &writeErr), storyMarshaller)
	if e := importStoryFiles(k, srcPath); e != nil {
		err = e
	} else if writeErr != nil {
		err = writeErr
	} else if e := assembleCat(&cat, outFile); e != nil {
		err = e
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
			err = Assemble(cat, db)
		}
	}
	return
}

// fix: this is probably only core marshal here.
func storyMarshaller(m jsn.Marshalee) (string, error) {
	return cout.Marshal(m, story.CompactEncoder)
}

func collectEphemera(cat *eph.Catalog, out *error) imp.WriterFun {
	// fix: needs to be more clever eventually...
	if e := cat.AddEphemera(
		"asm",
		&eph.EphBeginDomain{Name: "tapestry"}); e != nil {
		panic(e)
	}
	// built in kinds -- see ephKinds.go
	for _, k := range kindsOf.DefaultKinds {
		pk := k.Parent()
		if e := cat.AddEphemera(
			"built in kinds",
			&eph.EphKinds{Kinds: k.String(), From: pk.String()}); e != nil {
			panic(e)
		}
	}
	var i int
	return func(el eph.Ephemera) {
		if e := cat.AddEphemera(strconv.Itoa(i), el); e != nil {
			*out = errutil.Append(*out, e)
		}
		i++ // temp
	}
}

// read a comma-separated list of files and directories
func importStoryFiles(k *imp.Importer, srcPath string) (err error) {
	recurse := true
	if e := files.ReadPaths(srcPath, recurse,
		[]string{CompactExt, DetailedExt}, func(p string) error {
			return readOne(k, p)
		}); e != nil {
		err = errutil.New("couldn't read file", srcPath, e)
	} else {
		k.Flush()
	}
	return
}

func readOne(k *imp.Importer, path string) (err error) {
	log.Println("reading", path)
	if b, e := files.ReadFile(path); e != nil {
		err = e
	} else if script, e := decodeStory(path, b); e != nil {
		err = errutil.New("couldn't decode", path, "b/c", e)
	} else if e := ImportStory(k, path, script); e != nil {
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
		var curr story.Story // fix: should also be story file...
		if e := din.Decode(&curr, tapestry.Registry(), b); e != nil {
			err = e
		} else {
			ret = &story.StoryFile{
				StoryLines: curr.Reformat(),
			}
		}
	default:
		err = errutil.Fmt("unknown file type %q", ext)
	}
	return
}
