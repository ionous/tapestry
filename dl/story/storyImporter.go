package story

import (
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/iffy/ident"
)

// Importer helps read story specific json.
type Importer struct {
	// sometimes the importer needs to define a singleton like type or instance
	oneTime       map[string]bool
	autoCounter   ident.Counters
	entireGame    string
	env           StoryEnv
	activityDepth int
	Write         WriterFun
	Marshal       MarshalFun
}

// fix: add origin
type WriterFun func(eph eph.Ephemera)
type MarshalFun func(jsn.Marshalee) (string, error)

func NewImporter(writer WriterFun, marshal MarshalFun) *Importer {
	return &Importer{
		Write:       writer,
		Marshal:     marshal,
		oneTime:     make(map[string]bool),
		autoCounter: make(ident.Counters),
	}
}

func (k *Importer) ImportStory(path string, tgt jsn.Marshalee) (err error) {
	k.SetSource(path)
	return importStory(k, tgt)
}

func (k *Importer) SetSource(path string) {
	//
}

func (k *Importer) Env() *StoryEnv {
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

// AddImplicitAspect declares an assembler specified aspect and its traits
func (k *Importer) AddImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind; k.Once(src) {
		k.Write(&eph.EphKinds{Kinds: kind, From: "kind"})
	}
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		k.Write(&eph.EphAspects{Aspects: aspect, Traits: traits})
		k.Write(&eph.EphKinds{Kinds: kind, Contain: []eph.EphParams{eph.AspectParam(aspect)}})
	}
}

func importStory(k *Importer, tgt jsn.Marshalee) error {
	// presumably, this will have to be fixed eventually
	// so that states can be composited not monolithic.
	ts := chart.MakeEncoder()
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
				} else if scene, ok := flow.GetFlow().(*TestScene); !ok {
					err = errutil.Fmt("trying to post import something other than a flow")
				} else {
					// unpack the name, resolving "CurrentTest" to the name of the current test
					// fix? the most recent test might become the last popped test value
					n := scene.TestName.String()
					if n == TestName_CurrentTest {
						n = k.env.Recent.Test
					} else {
						k.env.Recent.Test = n
					}
					k.Write(&eph.EphBeginDomain{Name: n})
				}
				return
			},
			BlockEnd: func(b jsn.Block, _ interface{}) (err error) {
				k.Write(&eph.EphEndDomain{})
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
