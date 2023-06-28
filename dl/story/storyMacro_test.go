package story_test

import (
	"database/sql"
	_ "embed"
	"log"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/kr/pretty"
)

// use a macro to define a relationship between nouns
func TestMacros(t *testing.T) {
	// ugh. this setup.
	dt := testweave.NewWeaverOptions(t.Name(), func(db *sql.DB) rt.Runtime {
		qx, e := qdb.NewQueryx(db)
		if e != nil {
			panic(e)
		}
		return qna.NewRuntimeOptions(
			log.Writer(),
			qx,
			decode.NewDecoder(story.AllSignatures),
			qna.NewOptions(),
		)
	}, false)
	defer dt.Close()
	cat := dt.Catalog()
	//
	if curr, e := story.CompactDecode(storyMacroData); e != nil {
		t.Fatal(e)
	} else if e := cat.AssertDomainStart("tapestry", nil); e != nil {
		t.Fatal(e)
	} else if e := addDefaultKinds(cat); e != nil {
		t.Fatal(e)
	} else if e := story.ImportStory(cat, t.Name(), &curr); e != nil {
		t.Fatal(e)
	} else if e := cat.AssertDomainEnd(); e != nil {
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

func addDefaultKinds(n assert.Assertions) (err error) {
	for _, k := range kindsOf.DefaultKinds {
		if e := n.AssertAncestor(k.String(), k.Parent().String()); e != nil {
			err = e
			break
		}
	}
	return
}
