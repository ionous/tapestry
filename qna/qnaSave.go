package qna

import (
	"errors"
	"fmt"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables"
)

// undefined results if called in the middle of a turn
func (run *Runner) SaveGame(scene string) (err error) {
	if saveDir, e := getSaveDir(run); e != nil {
		err = e
	} else {
		// SAVE THE RUN DOMAIN COUNTER
		// SAVE RANDOMIZE

		if e := run.writeCounters(run.db); e != nil {
			err = e
		} else {

			name := files.NameWithTime(scene, files.SaveFileExtension)
			//
			outPath := filepath.Join(saveDir, name)
			err = tables.SaveFile(outPath, run.db)
		}

	}
	return
}

// undefined results if called in the middle of a turn
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
