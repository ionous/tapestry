package cmdweave

import (
	"bufio"
	"fmt"
	"io/fs"
	"path"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/support/flex"
)

const indexName = "_index"

// can return nil on success
func readIndexScene(fsys fs.FS, dir string, ents []fs.DirEntry) (ret *story.DefineScene, err error) {
	if idx := findIndexFile(ents); idx >= 0 {
		index := ents[idx]
		ret, err = readIndexFile(fsys, dir, index.Name())
	}
	return
}

// alt: we could just try to open the index file; but this way allows multiple extensions.
// find index ( if it exists ); return slice without it
// return a list containing the scene, and any dependencies
// scene can be nil if there is no script file; and that's okay
func findIndexFile(ents []fs.DirEntry) (ret int) {
	ret = -1 // provisionally
	for i, ent := range ents {
		name := ent.Name()
		if a, b := name[0], indexName[0]; a > b {
			break // entries are sorted; if we've passed the index, it doesnt exist.
		} else if a == b && !ent.IsDir() {
			if n, ext := files.SplitExt(name); ext.Story() && n == indexName {
				ret = i
				break
			}
		}
	}
	return
}

// open the passed file, error if there's anything more than a scene in it.
func readIndexFile(fsys fs.FS, dir, name string) (ret *story.DefineScene, err error) {
	fullpath := path.Join(dir, name)
	if fp, e := fsys.Open(fullpath); e != nil {
		err = e
	} else if els, e := flex.ReadStory(bufio.NewReader(fp)); e != nil {
		err = fmt.Errorf("couldn't read scene index %q because %s", fullpath, e)
	} else {
	Loop:
		for _, op := range els {
			switch op := op.(type) {
			case *story.Note:
				// continue looking
			case *story.DefineScene:
				if ret != nil {
					err = fmt.Errorf("expected scene index %q to contain a single scene", fullpath)
					break Loop
				} else if len(op.Statements) > 0 {
					err = fmt.Errorf("expected scene index %q to be empty", fullpath)
				} else {
					// record but check the next statement(s) too.
					ret = op
				}
			default:
				err = fmt.Errorf("expected scene index %q to contain a scene", fullpath)
				break Loop
			}
		}
	}

	return
}
