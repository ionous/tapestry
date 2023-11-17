package story_test

import (
	_ "embed"
	"encoding/json"
	"log"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/kr/pretty"
)

// use a macro to define a relationship between nouns
func TestMacros(t *testing.T) {
	// ugh. this setup.
	name := t.Name()
	db := testdb.Create(name)
	defer db.Close()

	if e := tables.CreateAll(db); e != nil {
		panic(e)
	}

	qx, e := qdb.NewQueries(db, false)
	if e != nil {
		panic(e)
	}
	run := qna.NewRuntime(
		log.Writer(),
		qx,
		decode.NewDecoder(story.AllSignatures),
	)
	cat := weave.NewCatalogWithWarnings(db, run, mdl.LogWarning)
	dt := testweave.NewWeaverCatalog(name, db, cat, true)
	//
	var msg map[string]any
	if e := json.Unmarshal(storyMacroData, &msg); e != nil {
		t.Fatal(e)
	} else if curr, e := story.CompactDecode(msg); e != nil {
		t.Fatal(e)
	} else if e := cat.DomainStart("tapestry", nil); e != nil {
		t.Fatal(e)
	} else if e := addDefaultKinds(cat.Pin("tapestry", "default kinds")); e != nil {
		t.Fatal(e)
	} else if e := story.ImportStory(cat, t.Name(), &curr); e != nil {
		t.Fatal(e)
	} else if e := cat.DomainEnd(); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
		t.Fatal(e)
	} else {
		if out, e := dt.ReadPairs(); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out, []string{
			"testing:whereabouts:hershel:scissors",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

//go:embed storyMacro_test.if
var storyMacroData []byte

func addDefaultKinds(pen *mdl.Pen) (err error) {
	for _, k := range kindsOf.DefaultKinds {
		if e := pen.AddKind(k.String(), k.Parent().String()); e != nil {
			err = e
			break
		}
	}
	return
}
