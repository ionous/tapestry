package cmdweave

import (
	"database/sql"
	"io/fs"
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

// Read all of the passed files and compile the output into
// a NEW database at outFile. This will attempt to erase any existing outFile.
// uses WalkDir which doesn't follow symlinks of sub-directories.
func WeavePaths(outFile string, stories ...fs.FS) (err error) {
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
				} else if e := importAll(cat, stories...); e != nil {
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

func importAll(cat *weave.Catalog, all ...fs.FS) (err error) {
	for _, fsys := range all {
		if fsys != nil {
			if e := importStoryFiles(cat, fsys); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// read a comma-separated list of files and directories
func importStoryFiles(cat *weave.Catalog, fsys fs.FS) (err error) {
	// note: walk does not expand symlinks in sub-directories.
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil { // e contains errors arising from visiting the file
			err = e
		} else if !d.IsDir() { // dirs are entered unless SkipDir is returned
			if ext := files.Ext(path); ext.Story() {
				log.Println("reading", path)
				if fp, e := fsys.Open(path); e != nil {
					err = e
				} else {
					var m map[string]any       // we "normalize" story files into json-like maps
					var script story.StoryFile // and decode from those maps

					if e := files.FormattedRead(fp, ext, &m); e != nil {
						err = errutil.New("couldn't read", path, "b/c", e)
					} else if e := story.Decode(&script, m); e != nil {
						err = errutil.New("couldn't decode", path, "b/c", e)
					} else if e := story.ImportStory(cat, path, &script); e != nil {
						err = errutil.New("couldn't import", path, "b/c", e)
					}
				}
			}
		}
		return
	})
}
