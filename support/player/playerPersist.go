package player

import (
	"errors"
	"fmt"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/files"
)

type Persistence struct {
	run   *qna.Runner
	query query.Query
}

func (p Persistence) LoadGame(scene string) (ret string, err error) {
	if saveDir, e := getSaveDir(p.run); e != nil {
		err = e
	} else if name, e := files.FindLatestNameWithTime(saveDir, scene, files.SaveFileExtension); e != nil {
		err = e
	} else if len(name) == 0 {
		err = fmt.Errorf("no save files found in %s", saveDir)
	} else {
		inPath := filepath.Join(saveDir, name)
		if d, e := p.query.LoadGame(inPath); e != nil {
			err = e
		} else {
			p.run.DynamicData().CacheMap = d
			ret = inPath
		}
	}
	return
}

func (p Persistence) SaveGame(scene string) (ret string, err error) {
	if saveDir, e := getSaveDir(p.run); e != nil {
		err = e
	} else {
		name := files.NameWithTime(scene, files.SaveFileExtension)
		outPath := filepath.Join(saveDir, name)
		if e := p.query.SaveGame(outPath, p.run.DynamicData().CacheMap); e != nil {
			err = e
		} else {
			ret = outPath
		}
	}
	return
}

func getSaveDir(run rt.Runtime) (ret string, err error) {
	if v, e := run.GetField(meta.Option, meta.SaveDir.String()); e != nil {
		err = fmt.Errorf("couldn't determine save directory, because %s", e)
	} else if str := v.String(); len(str) == 0 {
		err = errors.New("no save directory configured")
	} else {
		ret = str
	}
	return
}
