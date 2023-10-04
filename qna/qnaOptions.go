package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

// uses its own cache to preserve values across domain changes
type Options struct {
	options map[string]g.Value
}

var optionTypes = [meta.NumOptions]affine.Affinity{affine.Bool}

func NewOptions() Options {
	var staticAssert [1]struct{}
	_ = staticAssert[int(meta.NumOptions)-len(optionTypes)]
	//
	out := make(map[string]g.Value)
	for opt := meta.Options(0); opt < meta.NumOptions; opt++ {
		if v, e := g.NewDefaultValue(optionTypes[int(opt)], ""); e == nil {
			n := opt.String()
			out[n] = v
		}
	}
	return Options{out}
}

func (m *Options) SetOption(opt meta.Options, v g.Value) (err error) {
	return m.SetOptionByName(opt.String(), v)
}

// change an existing option.
func (m *Options) SetOptionByName(name string, v g.Value) (err error) {
	// see meta.Options: new options cannot be added dynamically at runtime
	if was, ok := m.options[name]; !ok {
		err = errutil.New("no such option", name)
	} else if a, va := was.Affinity(), v.Affinity(); a != va {
		err = errutil.New("option", name, "requires", a, "had", va, v)
	} else {
		m.options[name] = v
	}
	return
}

// return an existing option.
func (m Options) OptionByName(name string) (ret g.Value, err error) {
	if was, ok := m.options[name]; !ok {
		err = errutil.New("no such option", name)
	} else {
		ret = was
	}
	return
}

func (m Options) Option(opt meta.Options) (ret g.Value, err error) {
	return m.OptionByName(opt.String())
}
