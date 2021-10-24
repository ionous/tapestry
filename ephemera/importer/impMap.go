package importer

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

// expects that the very next value is a block
type KvMap map[string]func(jsn.Block, interface{}) error

func (kvm KvMap) call(b jsn.Block, key string, val interface{}) (err error) {
	if fn, ok := kvm[key]; ok {
		err = fn(b, val)
	}
	return
}

func Anywhere() *chart.StateMix {
	return &chart.StateMix{
		OnBlock: func(jsn.Block) error { return nil },
		OnKey:   func(_, _ string) error { return nil },
		OnValue: func(_ string, _ interface{}) error { return nil },
	}
}

func EveryValueOf(m *chart.Machine, typeName string, fn func(interface{}) error) (scope *BlockScope) {
	scope = NewBlockScope(m, chart.StateMix{
		OnValue: func(valType string, val interface{}) (err error) {
			if valType == typeName {
				err = fn(val)
			}
			return
		},
	})
	return
}

func Map(m *chart.Machine, typeName string, kvm KvMap) (scope *BlockScope) {
	var atKey string
	scope = NewBlockScope(m, chart.StateMix{
		OnBlock: func(b jsn.Block) (err error) {
			// does the incoming block provide the value of an open key.
			if key := atKey; len(key) > 0 {
				atKey = ""
				if at, ok := scope.At(typeName); ok {
					err = kvm.call(at, key, b)
				}
			}
			return nil
		},
		OnKey: func(lede, key string) error {
			atKey = key
			return nil
		},
		OnValue: func(valType string, val interface{}) (err error) {
			if key := atKey; len(key) > 0 {
				atKey = ""
				if at, ok := scope.At(typeName); ok {
					err = kvm.call(at, key, val)
				}
			}
			return
		},
		OnEnd: func() {
			atKey = ""
		},
	})
	return
}
