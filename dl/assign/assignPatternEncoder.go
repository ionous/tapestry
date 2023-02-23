package assign

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang"
)

// rewrite pattern calls to look like normal operations.
func EncodePattern(m jsn.Marshaler, op *CallPattern) (err error) {
	patName := recase(op.PatternName, false)
	pb := patternBlock(patName)
	if err = m.MarshalBlock(pb); err == nil {
		for _, arg := range op.Arguments {
			argName := recase(arg.Name, false)
			if e := m.MarshalKey(argName, argName); e != nil {
				err = e
				break
			} else if e := Assignment_Marshal(m, &arg.Value); e != nil {
				err = e
				break
			}
		}
		m.EndBlock()
	}
	return
}

// a fake block that writes the pattern name out as the lede
type patternBlock string

func (pb patternBlock) GetLede() string       { return string(pb) }
func (patternBlock) GetType() string          { return "patternBlock" }
func (patternBlock) GetFlow() interface{}     { return nil }
func (patternBlock) SetFlow(interface{}) bool { return false }

// pascal when true, camel when false
func recase(str string, cap bool) string {
	u := lang.Underscore(str)
	rs := []rune(u)
	var i int
	for j, cnt := 0, len(rs); j < cnt; j++ {
		if n := rs[j]; n == '_' {
			cap = true
		} else {
			if !cap {
				n = unicode.ToLower(n)
			} else {
				n = unicode.ToUpper(n)
				cap = false
			}
			rs[i] = n
			i++
		}
	}
	return string(rs[:i])
}
