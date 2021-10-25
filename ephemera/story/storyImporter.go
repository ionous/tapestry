package story

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/iffy/ident"
	"git.sr.ht/~ionous/iffy/tables"
)

// Importer helps read story specific json.
type Importer struct {
	*ephemera.Recorder
	// sometimes the importer needs to define a singleton like type or instance
	oneTime       map[string]bool
	autoCounter   ident.Counters
	entireGame    eph.Named
	env           StoryEnv
	activityDepth int
}

// low level
func NewImporter(db *sql.DB) *Importer {
	rec := ephemera.NewRecorder(db)
	k := &Importer{
		Recorder:    rec,
		oneTime:     make(map[string]bool),
		autoCounter: make(ident.Counters),
	}
	return k
}

func (k *Importer) ImportStory(path string, tgt jsn.Marshalee) (err error) {
	k.SetSource(path)
	return importStory(k, tgt)
}

func (k *Importer) NewName(name, category, ofs string) eph.Named {
	return k.NewDomainName(k.Env().Current.Domain, name, category, ofs)
}

func (k *Importer) Env() *StoryEnv {
	if !k.env.Game.Domain.IsValid() {
		k.env.Game.Domain = k.Recorder.NewName("entire_game", tables.NAMED_SCENE, "internal")
	}
	if !k.env.Current.Domain.IsValid() {
		k.env.Current.Domain = k.env.Game.Domain
	}
	return &k.env
}

// return true if m is the first time once has been called with the specified string.
func (k *Importer) Once(s string) (ret bool) {
	if !k.oneTime[s] {
		k.oneTime[s] = true
		ret = true
	}
	return
}

// a hopefully temporary hack
// fix? could replace with a push/pop of a state processor on execute
func (k *Importer) InProgram() bool {
	return k.activityDepth > 0
}

// NewImplicitAspect declares an assembler specified aspect and its traits
func (k *Importer) NewImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		domain := k.Env().Game.Domain
		kKind := k.NewDomainName(domain, kind, tables.NAMED_KINDS, src)
		kAspect := k.NewDomainName(domain, aspect, tables.NAMED_ASPECT, src)
		k.NewAspect(kAspect)
		k.NewField(kKind, kAspect, tables.PRIM_ASPECT, "")
		for i, trait := range traits {
			kTrait := k.NewDomainName(domain, trait, tables.NAMED_TRAIT, src)
			k.NewTrait(kTrait, kAspect, i)
		}
	}
}

func importStory(k *Importer, tgt jsn.Marshalee) error {
	// presumably, this will have to be fixed eventually
	// so that states can be composited not monolithic.
	ts := chart.MakeEncoder(nil)
	return ts.Marshal(tgt, Map(&ts, BlockMap{
		rt.Execute_Type: KeyMap{
			BlockStart: func(b jsn.Block, _ interface{}) (err error) {
				k.activityDepth++
				return
			},
			BlockEnd: func(b jsn.Block, _ interface{}) (err error) {
				k.activityDepth--
				return
			},
		},
		TestScene_Type: KeyMap{
			BlockStart: func(b jsn.Block, v interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); !ok {
					err = errutil.Fmt("trying to post import something other than a flow")
				} else if name, ok := flow.GetFlow().(*TestName); !ok {
					err = errutil.Fmt("trying to post import something other than a flow")
				} else {
					// unpack the name, resolving "CurrentTest" to the name of the current test
					if newDomain, e := NewTestName(k, *name); e != nil {
						err = e
					} else {
						// the most recent test might become the last popped test value
						// ( once domains and tests are stackable )
						k.env.Recent.Test = newDomain
						k.env.PushDomain(newDomain)
					}
				}
				return
			},
			BlockEnd: func(b jsn.Block, _ interface{}) (err error) {
				k.env.PopDomain()
				return
			},
		},
		OtherBlocks: KeyMap{
			BlockStart: func(b jsn.Block, v interface{}) (err error) {
				switch newBlock := b.(type) {
				case jsn.SlotBlock:
					if slat, ok := newBlock.GetSlot(); !ok {
						err = jsn.Missing
					} else {
						switch tgt := slat.(type) {
						case ImportStub:
							if rep, e := tgt.ImportStub(k); e != nil {
								err = errutil.New(e, "failed to create replacement")
							} else if !newBlock.SetSlot(rep) {
								err = errutil.New("failed to set replacement")
							}
						}
					}
				}
				return
			},
			BlockEnd: func(b jsn.Block, slot interface{}) (err error) {
				switch oldBlock := b.(type) {
				case jsn.FlowBlock:
					switch tgt := oldBlock.GetFlow().(type) {
					case StoryStatement:
						err = tgt.ImportPhrase(k)
					}
				}
				return
			},
		},
	}))
}

type ImportStub interface {
	ImportStub(*Importer) (interface{}, error)
}
