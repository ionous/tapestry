// for wasm: dump a scene's worth of data to json.
// ignores checks since this is for wasm playback.
package dump

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
)

//go:embed sql/*.sql
var queries embed.FS

func DumpAll(db *sql.DB, dec decoder.Decoder, scene string) (ret raw.Data, err error) {
	if scenes, e := QueryScenes(db, scene); e != nil {
		err = fmt.Errorf("%w for scenes", e)
	} else if plurals, e := QueryPlurals(db, scene); e != nil {
		err = fmt.Errorf("%w for plurals", e)
	} else if kinds, e := QueryKinds(db, dec, scene); e != nil {
		err = fmt.Errorf("%w for kinds", e)
	} else if names, e := QueryNames(db, scene); e != nil {
		err = fmt.Errorf("%w for names", e)
	} else if nouns, e := QueryNouns(db, kindDecoder{dec, kinds}, scene); e != nil {
		err = fmt.Errorf("%w for nouns", e)
	} else if patterns, e := QueryPatterns(db, dec, scene); e != nil {
		err = fmt.Errorf("%w for patterns", e)
	} else if relatives, e := QueryRelatives(db, scene); e != nil {
		err = fmt.Errorf("%w for relatives", e)
	} else if grammar, e := QueryGrammar(db, dec, scene); e != nil {
		err = fmt.Errorf("%w for grammar", e)
	} else {
		ret = raw.Data{
			Scenes:    scenes,
			Plurals:   plurals,
			Kinds:     kinds,
			Names:     names,
			Nouns:     nouns,
			Patterns:  patterns,
			Relatives: relatives,
			Grammar:   grammar,
		}
	}
	return
}

type kindDecoder struct {
	decoder.Decoder
	ks []rt.Kind
}

func (q kindDecoder) GetKindByName(exactKind string) (*rt.Kind, error) {
	return raw.FindKind(q.ks, exactKind)
}

func QueryScenes(db *sql.DB, scene string) (ret []string, err error) {
	if scenes, e := tables.QueryStrings(db, must("scenes"), scene); e != nil {
		err = e
	} else if len(scenes) == 0 {
		err = fmt.Errorf("unknown scene %s", scene)
	} else {
		ret = scenes
	}
	return
}

func must(name string) (ret string) {
	if b, e := fs.ReadFile(queries, "sql/"+name+".sql"); e != nil {
		panic(e)
	} else {
		ret = string(b)
	}
	return
}
