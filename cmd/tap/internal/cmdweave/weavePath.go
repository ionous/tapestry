package cmdweave

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// lets do this in the dumbest of ways for now.
func WeavePath(srcPath, outFile string) (err error) {
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
			// fix: why do we have to create qdb?
			if e := tables.CreateAll(db); e != nil {
				err = e
			} else if qx, e := qdb.NewQueries(db, false); e != nil {
				err = e
			} else {
				run := qna.NewRuntime(
					log.Writer(),
					qx,
					decode.NewDecoder(story.AllSignatures),
				)
				cat := weave.NewCatalogWithWarnings(db, run, nil)
				if e := cat.DomainStart("tapestry", nil); e != nil {
					err = e
				} else if e := addDefaultKinds(cat.Pin("tapestry", "default kinds")); e != nil {
					err = e
				} else if e := importStoryFiles(cat, srcPath); e != nil {
					err = e
				} else if e := cat.DomainEnd(); e != nil {
					err = e
				} else {
					err = cat.AssembleCatalog()
				}
			}
		}
	}
	return
}

func addDefaultKinds(pen *mdl.Pen) (err error) {
	for _, k := range kindsOf.DefaultKinds {
		if e := pen.AddKind(k.String(), k.Parent().String()); e != nil {
			err = e
			break
		}
	}
	return
}

// read a comma-separated list of files and directories
func importStoryFiles(cat *weave.Catalog, srcPath string) (err error) {
	recurse := true
	if e := files.ReadPaths(srcPath, recurse, storyExts, func(p string) error {
		return readOne(cat, p)
	}); e != nil {
		err = errutil.New("couldn't read file", srcPath, e)
	}
	return
}

var storyExts = []string{
	files.CompactExt.String(),
	files.TellStory.String()}

func readOne(cat *weave.Catalog, path string) (err error) {
	log.Println("reading", path)
	if script, e := decodeStory(path); e != nil {
		err = errutil.New("couldn't decode", path, "b/c", e)
	} else if e := story.ImportStory(cat, path, &script); e != nil {
		err = errutil.New("couldn't import", path, "b/c", e)
	}
	return
}

func decodeStory(path string) (ret story.StoryFile, err error) {
	switch ext := files.Ext(path); ext {
	case files.TellStory:
		var msg map[string]any
		if e := files.LoadTell(path, &msg); e != nil {
			err = e
		} else {
			err = story.Decode(&ret, msg)
		}
	case files.CompactExt:
		var msg map[string]any
		if b, e := files.ReadFile(path); e != nil {
			err = e
		} else if e := json.Unmarshal(b, &msg); e != nil {
			err = e
		} else {
			err = story.Decode(&ret, msg)
		}
	default:
		err = errutil.Fmt("unknown file type %q", ext)
	}
	return
}
