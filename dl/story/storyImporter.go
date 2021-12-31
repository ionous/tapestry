package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Importer helps read story specific json.
type Importer struct {
	// sometimes the importer needs to define a singleton like type or instance
	oneTime       map[string]bool
	autoCounter   Counters
	entireGame    string
	env           StoryEnv
	activityDepth int
	writer        WriterFun
	Marshal       MarshalFun
	queue         []eph.Ephemera
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

func (k *Importer) Write(op eph.Ephemera) {
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

// a hopefully temporary hack
// fix? could replace with a push/pop of a state processor on execute
func (k *Importer) InProgram() bool {
	return k.activityDepth > 0
}

// AddImplicitAspect declares an assembler specified aspect and its traits
func (k *Importer) AddImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		k.WriteOnce(&eph.EphAspects{Aspects: aspect, Traits: traits})
		k.WriteOnce(&eph.EphKinds{Kinds: kind, Contain: []eph.EphParams{eph.AspectParam(aspect)}})
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
