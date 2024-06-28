package qna

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt/meta"
)

// produces undefined results if called in the middle of a turn
// returns the saved file name
func (run *Runner) SaveGame(scene string) (ret string, err error) {
	// FIX! shouldnt depend on file here --
	// maybe part of qdb --
	// or a separate helper interface or something
	err = errors.New("cant save")
	// if saveDir, e := getSaveDir(run); e != nil {
	// 	err = e
	// } else {
	// 	name := files.NameWithTime(scene, files.SaveFileExtension)
	// 	outPath := filepath.Join(saveDir, name)
	// 	err = run.query.SaveGame(outPath, run.dynamicVals.CacheMap)
	// }
	return
}

// produces undefined results if called in the middle of a turn
// returns the file that was loaded
func (run *Runner) LoadGame(scene string) (ret string, err error) {
	// FIX! where's the code to load gone to?
	err = errors.New("cant load")
	// if saveDir, e := getSaveDir(run); e != nil {
	// 	err = e
	// } else if name, e := files.FindLatestNameWithTime(saveDir, scene, files.SaveFileExtension); e != nil {
	// 	err = e
	// } else if len(name) == 0 {
	// 	err = fmt.Errorf("no save files found in %s", saveDir)
	// } else {
	// 	outPath := filepath.Join(saveDir, name)
	// 	run.Query.LoadGame(outPath)
	// 	defer func() {
	// 		if err != nil {
	// 			err = fmt.Errorf("%w for %s", err, outPath)
	// 		}
	// 	}()
	// }
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
