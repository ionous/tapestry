package imp

import (
	"log"

	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Importer helps read story specific json.
type Importer struct {
	// the importer uses a runtime so that it can handle macros.
	*qna.Runner
	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	macros      macroReg
	autoCounter Counters
	env         Environ
	assert.Assertions
}

// PostImport - happens at the end of a json block after all of its dependencies have been PostImported.
// Most story statements implement this interface.
type PostImport interface {
	// fix: the things that require importer directly should be moved to script
	// and then this could pass EphemeraWriter.
	PostImport(*Importer) error
}

// PreImport happens at the opening of a json block and it can transform the value into something completely new.
type PreImport interface {
	PreImport(*Importer) (interface{}, error)
}

// StoryStatement - import a single story statement.
// used during weave, and expects that the runtime is the importer's own runtime.
// ( as opposed to the story's playtime. )
func StoryStatement(run rt.Runtime, op PostImport) (err error) {
	if k, ok := run.(*Importer); !ok {
		err = errutil.Fmt("runtime %T doesn't support story statements", run)
	} else {
		err = op.PostImport(k)
	}
	return
}

func NewImporter(w assert.Assertions) *Importer {
	k := &Importer{
		macros:      make(macroReg),
		oneTime:     make(map[string]bool),
		autoCounter: make(Counters),
		Runner: qna.NewRuntimeOptions(
			log.Writer(),
			query.QueryNone("import doesn't support object queries"),
			qna.DecodeNone("import doesn't support the decoder"),
			qna.NewOptions()),
	}
	k.Assertions = w
	return k
}

func (k *Importer) SetSource(string) {
}

// Env - used for comments to determine if they should turn into log statements.
// todo: remove?
func (k *Importer) Env() *Environ {
	return &k.env
}

// NewCounter generates a unique string, and uses local markup to try to create a stable one.
// instead consider  "PreImport" could be used to write a key into the markup if one doesnt already exist.
// and a free function could also extract what it needs from any op's markup.
// ( then PostImport wouldn't need Importer for counters )
func (k *Importer) NewCounter(name string, markup map[string]any) (ret string) {
	// fix: use a special "id" marker instead?
	if at, ok := markup["comment"].(string); ok && len(at) > 0 {
		ret = at
	} else {
		ret = k.autoCounter.Next(name)
	}
	return
}
