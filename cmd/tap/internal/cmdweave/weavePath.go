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
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/support/flex"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type NamedFS struct {
	Name string
	fs.FS
}

// Read all of the passed files and compile the output into
// ents NEW database at outFile. This will attempt to erase any existing outFile.
// uses WalkDir which doesn't follow symlinks of sub-directories.
func WeavePaths(outFile string, stories ...NamedFS) (err error) {
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
			if q, e := qdb.NewQueries(db, query.NewDecoder(story.AllSignatures)); e != nil {
				err = e
			} else {
				run := qna.NewRuntime(q)
				cat := weave.NewCatalogWithWarnings(db, run, nil)
				d := cat.EnsureScene("tapestry")
				pos := compact.Source{
					File:    "default kinds",
					Line:    -1,
					Comment: "the built-in types",
				}
				if pen, e := cat.SceneBegin(d, pos, nil); e != nil {
					err = e
				} else {
					// mark it as a root scene:
					// fix? https://todo.sr.ht/~ionous/tapestry/60
					if e := pen.AddDependency(""); e != nil {
						err = e
					} else if e := addDefaultKinds(pen); e != nil {
						err = e
					} else if e := importAll(cat, stories...); e != nil {
						err = e
					} else {
						cat.SceneEnd() // doesn't pop on error; oh well.
						err = cat.AssembleCatalog()
					}
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

func importAll(cat *weave.Catalog, all ...NamedFS) (err error) {
	for _, fsys := range all {
		if e := importDir(cat, fsys, []string{"."}); e != nil {
			err = e
			break
		}
	}
	return
}

// importDir recursively descends path
func importDir(cat *weave.Catalog, fsys NamedFS, dirs []string) (err error) {
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
			} else {
				// hack for the shared library, so it doesnt try to depend on itself
				if d := cat.EnsureScene(scene); d.Name() != "tapestry" {
					// tbd: what's a helpful source path
					fullpath := path.Join(dir, indexName)
					pos := compact.Source{File: fullpath}
					if pen, e := cat.SceneBegin(d, pos, nil); e != nil {
						err = e
					} else {
						defer cat.SceneEnd() // called at end of function after collecting everything.
						err = pen.AddDependency(req...)
					}
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
					fullpath := path.Join(dir, name)
					if fp, e := fsys.Open(fullpath); e != nil {
						err = e
					} else {
						namedPath := path.Join(fsys.Name, fullpath)
						if els, e := flex.ReadStorySource(namedPath, bufio.NewReader(fp)); e != nil {
							err = fmt.Errorf("couldn't read %q because %s", namedPath, e)
						} else {
							// without a scene index we need to have one for the file itself.
							if scene == nil {
								els = wrapScene(shortName, namedPath, els)
							}
							if e := story.Weave(cat, els); e != nil {
								err = fmt.Errorf("couldn't import %q because %s", namedPath, e)
							}
						}
						fp.Close()
					}
				}
			}
		}
	}
	return
}

func (l NamedFS) Open(fullpath string) (fs.File, error) {
	if fullpath != "." { // avoid logging the root dir
		log.Println("reading", fullpath)
	}
	return l.FS.Open(fullpath)
}
