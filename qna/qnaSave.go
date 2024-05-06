package qna

import (
	"errors"
	"fmt"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pack"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables"
)

// produces undefined results if called in the middle of a turn
func (run *Runner) SaveGame(scene string) (err error) {
	if saveDir, e := getSaveDir(run); e != nil {
		err = e
	} else {
		// to do: save randomizer
		err = writeValues(run.db, func(write writeCb) (err error) {
			if e := run.rand.writeRandomizeer(write); e != nil {
				err = e
			} else if e := run.count.writeCounters(write); e != nil {
				err = e
			} else if e := run.writeNouns(write); e != nil {
				err = e
			} else {
				name := files.NameWithTime(scene, files.SaveFileExtension)
				//
				outPath := filepath.Join(saveDir, name)
				err = tables.SaveFile(outPath, run.db)
			}
			return
		})

	}
	return
}

// produces undefined results if called in the middle of a turn
func (run *Runner) LoadGame(scene string) (err error) {
	if saveDir, e := getSaveDir(run); e != nil {
		err = e
	} else if name, e := files.FindLatest(saveDir, scene, files.SaveFileExtension); e != nil {
		err = e
	} else {
		outPath := filepath.Join(saveDir, name)
		err = tables.LoadFile(run.db, outPath)

	}
	return
}

func getSaveDir(run *Runner) (ret string, err error) {
	if v, e := run.GetField(meta.Option, meta.SaveDir.String()); e != nil {
		err = fmt.Errorf("couldn't determine save directory, because %s", e)
	} else if str := v.String(); len(str) == 0 {
		err = errors.New("no save directory configured")
	} else {
		ret = str
	}
	return
}

// write all dynamic values to the database using the prepared 'runValue' statement.
func (run *Runner) writeNouns(w writeCb) (err error) {
	for key, cached := range run.dynamicVals.store {
		// all stored values are variant values ( unless they are errors )
		if val, ok := cached.v.(rt.Value); ok {
			if str, e := pack.PackValue(val); e != nil {
				err = e
			} else if e := w(key.group, key.target, key.field, str); e != nil {
				err = e
				break
			}
		}
	}
	return
}
