package story

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
)

func EveryValueOf(m *chart.Machine, typeName string, fn func(interface{}) error) chart.State {
	var blocks BlockStack
	return &chart.StateMix{
		OnBlock: func(b jsn.Block) (err error) {
			blocks.Push(b)
			return
		},
		OnKey: func(lede, key string) (err error) {
			return
		},
		OnValue: func(valType string, val interface{}) (err error) {
			if valType == typeName {
				err = fn(val)
			}
			return
		},
		OnEnd: func() {
			if _, ok := blocks.Pop(); !ok {
				m.FinishState("scope") // pop this.
			}
		},
	}
}
