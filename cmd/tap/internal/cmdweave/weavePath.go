package cmdweave

import (
	"bufio"
	"fmt"
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
	"git.sr.ht/~ionous/tapestry/support/flex"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// Read all of the passed files and compile the output into
// ents NEW database at outFile. This will attempt to erase any existing outFile.
// uses WalkDir which doesn't follow symlinks of sub-directories.
func WeavePaths(outFile string, stories ...fs.FS) (err error) {
	if outFile, e := filepath.Abs(outFile); e != nil {
		err = e
	} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
		err = fmt.Errorf("couldn't clean output file %q because %s", outFile, e)
	} else {
		// 0755 -> readable by all but only writable by the user
		// 0700 -> read/writable by user
		// 0777 -> ModePerm ... read/writable by all
		os.MkdirAll(path.Dir(outFile), os.ModePerm)
		if db, e := tables.CreateBuildTime(outFile); e != nil {
			err = fmt.Errorf("couldn't create output file %q because %s", outFile, e)
		} else {
			defer db.Close()
			if run, e := qna.NewRuntime(db, decode.NewDecoder(story.AllSignatures)); e != nil {
				err = e
			} else {
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
			if e := importDir(cat, loggingFS{fsys}, []string{"."}); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// importDir recursively descends path
func importDir(cat *weave.Catalog, fsys fs.FS, dirs []string) (err error) {
	dir := path.Join(dirs...)
	if ents, e := fs.ReadDir(fsys, dir); e != nil {
		err = e
	} else if scene, e := readIndexScene(fsys, dir, ents); e != nil {
		err = e
	} else {
		// a scene index affects the whole directory
		if scene != nil {
			run := cat.GetRuntime()
			if scene, req, e := scene.GetSceneReqs(run); e != nil {
				err = e
			} else if inflect.Normalize(scene) != "tapestry" {
				// ^ hack for the shared library, so it doesnt try to depend on itself
				if e := cat.DomainStart(scene, req); e != nil {
					err = e
				} else {
					defer cat.DomainEnd() // called at end of function
				}
			}
		}
		// loop over the folders and files in the the directory:
		for i, cnt := 0, len(ents); i < cnt && err == nil; i++ {
			ent := ents[i]
			// filenames starting with a dot ( `.` ) or underscore ( `_` ) are ignored.
			if name := ent.Name(); name[0] == '_' || name[0] == '.' {
				continue
			} else if ent.IsDir() {
				err = importDir(cat, fsys, append(dirs, name))
			} else {
				if shortName, ext := files.SplitExt(name); ext.Story() {
					fullpath := path.Join(name)
					if fp, e := fsys.Open(fullpath); e != nil {
						err = e
					} else {
						cat.BeginFile(fullpath)
						if els, e := flex.ReadStory(bufio.NewReader(fp)); e != nil {
							err = fmt.Errorf("couldn't read %q because %s", fullpath, e)
						} else {
							// without a scene index we need to have one for the file itself.
							if scene == nil {
								els = wrapScene(shortName, els)
							}
							if e := story.Weave(cat, els); e != nil {
								err = fmt.Errorf("couldn't import %q because %s", fullpath, e)
							}
						}
						cat.EndFile()
						fp.Close()
					}
				}
			}
		}
	}
	return
}

// helper to log every open
type loggingFS struct {
	fs.FS
}

func (l loggingFS) Open(fullpath string) (fs.File, error) {
	if fullpath != "." { // avoid logging the root dir
		log.Println("reading", fullpath)
	}
	return l.FS.Open(fullpath)
}
