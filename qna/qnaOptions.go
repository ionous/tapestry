package qna

import (
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

// uses its own cache to preserve values across domain changes
type Options struct {
	options map[string]rt.Value
}

// indexed by meta.Options
var optionTypes = [...]affine.Affinity{affine.Bool, affine.Bool, affine.Bool, affine.Text}

func NewOptions() Options {
	// if this fails, then optionsTypes has to be adjusted to match the meta options
	var staticAssert [1]struct{}
	_ = staticAssert[int(meta.NumOptions)-len(optionTypes)]
	//
	out := make(map[string]rt.Value)
	for opt := meta.Options(0); opt < meta.NumOptions; opt++ {
		if v, e := rt.ZeroValue(optionTypes[int(opt)]); e == nil {
			n := opt.String()
			out[n] = v
		}
	}
	return Options{out}
}

func (m *Options) SetOption(opt meta.Options, v rt.Value) (err error) {
	return m.SetOptionByName(opt.String(), v)
}

// change an existing option.
func (m *Options) SetOptionByName(name string, v rt.Value) (err error) {
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
func (m Options) OptionByName(name string) (ret rt.Value, err error) {
	if was, ok := m.options[name]; !ok {
		err = errutil.New("no such option", name)
	} else {
		ret = was
	}
	return
}

func (m Options) Option(opt meta.Options) (ret rt.Value, err error) {
	return m.OptionByName(opt.String())
}

func (m Options) IsOption(opt meta.Options) (okay bool) {
	if flag, e := m.Option(opt); e != nil {
		log.Printf("couldn't determine option %s because %s.\n", opt, e)
	} else if flag.Affinity() == affine.Bool && flag.Bool() {
		okay = true
	}
	return
}

func (m Options) cacheErrors() (okay bool) {
	if m, ok := m.options[meta.CacheErrors.String()]; ok {
		okay = m.Affinity() == affine.Bool && m.Bool()
	}
	return
}
