// for wasm: dump a scene's worth of data to json.
// ignores checks since this is for wasm playback.
package dump

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

//go:embed sql/*.sql
var queries embed.FS

type AllData struct {
	Scenes    []string   // in order with root ( tapestry ) first
	Plurals   []Plural   // pairs of one, many; sorted by one
	Kinds     []KindData // sorted by .Kind
	Names     []NounName // sorted by .Name
	Nouns     []NounData
	Patterns  []PatternData
	Relatives []RelativeData // sorted by .Relation
}

type KindData struct {
	Id        int
	Domain    string
	Kind      string
	Ancestors []string // includes itself
	Fields    []query.FieldData
}

type NounName struct {
	Id   int
	Name string
}

type NounData struct {
	Id     int // mdl_noun
	Domain string
	Noun   string // unique id
	Kind   string // or would id be better?
	// Name list? or best name maybe?
	Values []query.ValueData
}

type PatternData struct {
	Id      int // mdl_pat
	Pattern string
	Labels  []string
	Result  string
	Rules   []query.RuleData
}

type RelativeData struct {
	Id       int    // mdl_rel
	Relation string // a kind
	Pairs    []Pair // can be empty
}

type Plural struct {
	One, Other string
}

type Pair struct {
	One, Other int // mdl_noun
}

func DumpAll(db *sql.DB, scene string) (ret AllData, err error) {
	if scenes, e := QueryScenes(db, scene); e != nil {
		err = fmt.Errorf("%w while querying scenes", e)
	} else if plurals, e := QueryPlurals(db, scene); e != nil {
		err = fmt.Errorf("%w while querying plurals", e)
	} else if kinds, e := QueryKinds(db, scene); e != nil {
		err = fmt.Errorf("%w while querying kinds", e)
	} else if names, e := QueryNames(db, scene); e != nil {
		err = fmt.Errorf("%w while querying names", e)
	} else if nouns, e := QueryNouns(db, scene); e != nil {
		err = fmt.Errorf("%w while querying nouns", e)
	} else if patterns, e := QueryPatterns(db, scene); e != nil {
		err = fmt.Errorf("%w while querying patterns", e)
	} else if relatives, e := QueryRelatives(db, scene); e != nil {
		err = fmt.Errorf("%w while querying relatives", e)
	} else {
		ret = AllData{
			Scenes:    scenes,
			Plurals:   plurals,
			Kinds:     kinds,
			Names:     names,
			Nouns:     nouns,
			Patterns:  patterns,
			Relatives: relatives,
		}
	}
	return
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
