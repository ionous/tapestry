package imp

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
	"log"
)

// Importer helps read story specific json.
type Importer struct {
	// the importer uses a runtime so that it can handle macros.
	rt.Runtime
	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	autoCounter Counters
	env         Environ
	writer      WriterFun
	//Marshal     MarshalFun
	queue []eph.Ephemera
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

type EphemeraWriter interface{ WriteEphemera(eph.Ephemera) }

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

// fix: add origin
type WriterFun func(eph eph.Ephemera)

//type MarshalFun func(jsn.Marshalee) (string, error)

func NewImporter(writer WriterFun) *Importer {
	return &Importer{
		writer: writer,
		//Marshal:     marshal,
		oneTime:     make(map[string]bool),
		autoCounter: make(Counters),
		Runtime: qna.NewRuntimeOptions(
			log.Writer(),
			query.QueryNone("import doesn't support object queries"),
			qna.DecodeNone("import doesn't support the decoder"),
			qna.NewOptions()),
	}
}

func (k *Importer) SetSource(string) {
}

// Env - used for comments to determine if they should turn into log statements.
// todo: remove?
func (k *Importer) Env() *Environ {
	return &k.env
}

// WriteEphemera - implements EphemeraWriter; the key part of importation.
func (k *Importer) WriteEphemera(op eph.Ephemera) {
	k.writer(op)
}

// WriteOnce elevates commands to the outermost domain.
// todo: all things using write once should live in if.script
func (k *Importer) WriteOnce(op eph.Ephemera) {
	k.queue = append(k.queue, op)
}

// Queued -> for tests, possibly could /should be removed after WriteOnce is removed.
func (k *Importer) Queued() []eph.Ephemera {
	return k.queue
}

func (k *Importer) Flush() {
	for i, q := range k.queue {
		k.writer(q)
		k.queue[i] = nil
	}
	k.queue = nil
}

// return true if m is the first time once has been called with the specified string.
// todo: move all calls into script.
func (k *Importer) Once(s string) (ret bool) {
	if !k.oneTime[s] {
		k.oneTime[s] = true
		ret = true
	}
	return
}

// AddImplicitAspect declares an assembler specified aspect and its traits
// todo: move all calls into script.
func (k *Importer) AddImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		k.WriteOnce(&eph.EphAspects{Aspects: aspect, Traits: traits})
		k.WriteOnce(&eph.EphKinds{Kind: kind, Contain: []eph.EphParams{eph.AspectParam(aspect)}})
	}
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
