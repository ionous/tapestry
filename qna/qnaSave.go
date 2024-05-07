package qna

import (
	"errors"
	"fmt"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables"
)

// produces undefined results if called in the middle of a turn
// returns the saved file name
func (run *Runner) SaveGame(scene string) (ret string, err error) {
	if saveDir, e := getSaveDir(run); e != nil {
		err = e
	} else if e := writeValues(run.db, func(w writeCb) (err error) {
		if e := run.rand.writeRandomizeer(w); e != nil {
			err = e
		} else if e := run.dynamicVals.writeValues(w); e != nil {
			err = e
		}
		return
	}); e != nil {
		err = e
	} else {
		name := files.NameWithTime(scene, files.SaveFileExtension)
		outPath := filepath.Join(saveDir, name)
		if e := tables.SaveFile(outPath, false, run.db); e != nil {
			err = e
		} else {
			ret = outPath
		}
	}
	return
}

// produces undefined results if called in the middle of a turn
// returns the file that was loaded
func (run *Runner) LoadGame(scene string) (ret string, err error) {
	if saveDir, e := getSaveDir(run); e != nil {
		err = e
	} else if name, e := files.FindLatestNameWithTime(saveDir, scene, files.SaveFileExtension); e != nil {
		err = e
	} else if len(name) == 0 {
		err = fmt.Errorf("no save files found in %s", saveDir)
	} else {
		outPath := filepath.Join(saveDir, name)
		defer func() {
			if err != nil {
				err = fmt.Errorf("%w for %s", err, outPath)
			}
		}()
		if e := tables.LoadFile(run.db, outPath); e != nil {
			err = e
		} else if e := run.dynamicVals.readValues(run.db); e != nil {
			err = e
		} else if e := run.rand.readRandomizer(run.db); e != nil {
			err = e
		} else {
			ret = outPath
		}
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
