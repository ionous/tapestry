package imp

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/jsn"
)

// Importer helps read story specific json.
type Importer struct {
	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	autoCounter Counters
	env         Environ
	writer      WriterFun
	Marshal     MarshalFun
	queue       []eph.Ephemera
}

// Post happens at the end of a json block after all of its dependencies have been imported.
// This is generally the function most story statements implement.
type PostImport interface {
	PostImport(*Importer) error
}

// PreImport happens at the opening of a json block and it can transform the value into something completely new.
type PreImport interface {
	PreImport(*Importer) (interface{}, error)
}

// fix: add origin
type WriterFun func(eph eph.Ephemera)
type MarshalFun func(jsn.Marshalee) (string, error)

func NewImporter(writer WriterFun, marshal MarshalFun) *Importer {
	return &Importer{
		writer:      writer,
		Marshal:     marshal,
		oneTime:     make(map[string]bool),
		autoCounter: make(Counters),
	}
}

func (k *Importer) SetSource(path string) {
	//
}

func (k *Importer) Env() *Environ {
	return &k.env
}

func (k *Importer) WriteEphemera(op eph.Ephemera) {
	k.writer(op)
}

// put the passed ephemera into the global scope
// ( good for autogenerated or implicit data )
// ( fix: these should probably all move to .if files )
func (k *Importer) WriteOnce(op eph.Ephemera) {
	k.queue = append(k.queue, op)
}

// exposed for tests
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
func (k *Importer) Once(s string) (ret bool) {
	if !k.oneTime[s] {
		k.oneTime[s] = true
		ret = true
	}
	return
}

// AddImplicitAspect declares an assembler specified aspect and its traits
func (k *Importer) AddImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		k.WriteOnce(&eph.EphAspects{Aspects: aspect, Traits: traits})
		k.WriteOnce(&eph.EphKinds{Kinds: kind, Contain: []eph.EphParams{eph.AspectParam(aspect)}})
	}
}

// generate a unique name for the counter --
// for stability's sake, preferring an existing id in the source to an autogenerated id.
func (k *Importer) NewCounter(name string, markup map[string]any) (ret string) {
	// fix: use a special "id" marker instead?
	if at, ok := markup["comment"].(string); ok && len(at) > 0 {
		ret = at
	} else {
		ret = k.autoCounter.Next(name)
	}
	return
}
