package compact

import (
	r "reflect"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/export"
	"git.sr.ht/~ionous/iffy/export/tag"
	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

// Compacter reads programs from json.
type Compacter struct {
	source     string
	cmds       map[string]composer.Composer
	issueFn    IssueReport
	IssueCount int
}
type IssueReport func(reader.Position, error)

func NewCompacter() *Compacter {
	reportNothing := func(reader.Position, error) {}
	return NewCompacterReporter(reportNothing)
}

func NewCompacterReporter(reporter IssueReport) *Compacter {
	return &Compacter{cmds: make(map[string]composer.Composer), issueFn: reporter}
}

func (m *Compacter) SetSource(source string) *Compacter {
	m.source = source
	return m
}

func (m *Compacter) report(ofs string, err error) {
	m.issueFn(reader.Position{Source: m.source, Offset: ofs}, err)
	m.IssueCount++
}

func (dec *Compacter) AddTypes(slats []composer.Composer) {
	for _, cmd := range slats {
		dec.addType(cmd)
	}
}

// AddCallback registers a command parser.
func (dec *Compacter) addType(cmd composer.Composer) {
	n := composer.SpecName(cmd)
	if was, exists := dec.cmds[n]; exists {
		panic(errutil.Fmt("conflicting name for spec %q %q!=%T", n, was, cmd))
	} else {
		dec.cmds[n] = cmd
	}
}

func (dec *Compacter) Compact(m reader.Map) (ret interface{}, err error) {
	return dec.readItem(m)
}

// note: this cant read slots, and it cant read repeated items
func (dec *Compacter) readItem(m reader.Map) (ret interface{}, err error) {
	itemType, itemAt := m.StrOf(reader.ItemType), reader.At(m)
	if cmd, ok := dec.cmds[itemType]; !ok {
		err = errutil.Fmt("unknown type %q at %q", itemType, itemAt)
	} else {
		spec := cmd.Compose()
		switch use := spec.Uses; use {
		case composer.Type_Flow:
			ret = dec.readFields(cmd, itemAt, m.MapOf(reader.ItemValue))
		case composer.Type_Swap:
			ret, err = dec.readSwap(cmd, m)
		case composer.Type_Str:
			ret, err = dec.readString(spec, m)
		case composer.Type_Num:
			ret, err = dec.readFloat(spec, m)
		default:
			err = errutil.New("unhandled", use, "at", itemAt)
		}
	}
	return
}

func (dec *Compacter) readFields(cmd composer.Composer, at string, srcData reader.Map) (ret reader.Map) {
	var sig sigBuilder
	spec := cmd.Compose()
	sig.WriteLede(spec)
	//
	var tokens []string
	var fields []*r.StructField
	// build a flat list of the fields we expect
	refType := r.TypeOf(cmd).Elem()
	export.WalkProperties(refType, func(field *r.StructField, path []int) (done bool) {
		if t := tag.ReadTag(field.Tag); !t.Exists("internal") {
			token := export.Tokenize(field.Name)
			if !t.Exists("optional") || srcData.Contains(token) {
				tokens = append(tokens, token)
				fields = append(fields, field)
				//
				label, _ := t.Find("label")
				if len(label) == 0 {
					dec.report(at, errutil.Fmt("missing label for %s.%s at %s", refType, field.Name, at))
					label = lang.Underscore(field.Name)
				}
				sig.WriteLabel(label)
			}
		}
		return
	})
	// report on keys that are in the src data but not our list of tokens.
	for key, _ := range srcData {
		i, cnt := 0, len(tokens)
		for ; i < len(tokens); i++ {
			if key == tokens[i] {
				break
			}
		}
		if i == cnt {
			dec.report(at, errutil.Fmt("unprocessed %q", key))
		}
	}
	// finally: process the tokens.
	var args []interface{}
	for i, token := range tokens {
		var val interface{}
		// see if we have data for this token.
		if fieldData, exists := srcData[token]; !exists {
			dec.report(at, errutil.Fmt("missing %s at %s", token, at))
		} else if v, e := dec.readField(fields[i].Type, fieldData); e != nil {
			dec.report(at, errutil.Fmt("%s compacting %s.%s", e, refType.String(), token))
		} else {
			val = v
		}
		args = append(args, val)
	}
	if n, cnt := sig.String(), len(args); cnt == 0 {
		ret = reader.Map{n: true}
	} else if cnt == 1 {
		ret = reader.Map{n: args[0]}
	} else {
		ret = reader.Map{n: args}
	}
	return
}

// to be here means the field exists and is non-optional
func (dec *Compacter) readField(refType r.Type, el interface{}) (ret interface{}, err error) {
	// expand pointers
	k := refType.Kind()
	if k == r.Ptr {
		refType = refType.Elem()
		k = refType.Kind()
	}
	//
	switch k {
	case r.Interface:
		if item, ok := el.(map[string]interface{}); !ok {
			err = errutil.Fmt("expected a map, got %T", el)
		} else {
			ret, err = dec.readSlot(item)
		}

	case r.Slice:
		if items, ok := el.([]interface{}); !ok {
			err = errutil.New("expected a slice")
		} else if k := refType.Elem().Kind(); k != r.Interface {
			ret, err = dec.readItems(items)
		} else {
			ret, err = dec.readSlots(items)
		}

	default:
		if item, ok := el.(map[string]interface{}); !ok {
			err = errutil.Fmt("expected a map, got %T", el)
		} else {
			ret, err = dec.readItem(item)
		}
	}
	return
}

func (dec *Compacter) readString(spec composer.Spec, m reader.Map) (ret interface{}, err error) {
	if el, ok := m[reader.ItemValue]; !ok {
		err = errutil.New("missing value while reading string")
	} else if v, ok := el.(string); !ok {
		err = errutil.New("expected a string")
	} else {
		if v[0] != '$' {
			if spec.OpenStrings {
				ret = v
			} else {
				err = errutil.New("invalid string", v)
			}
		} else {
			if s, i := spec.IndexOfChoice(v); i < 0 {
				err = errutil.New("invalid string", v)
			} else if v == "$TRUE" {
				ret = true
			} else if v == "$FALSE" {
				ret = false
			} else {
				ret = []string{s}
			}
		}
	}
	return
}

func (dec *Compacter) readFloat(cmd composer.Spec, m reader.Map) (ret float64, err error) {
	if el, ok := m[reader.ItemValue]; !ok {
		err = errutil.New("missing value while reading string")
	} else if v, ok := el.(float64); !ok {
		err = errutil.New("expected a string")
	} else {
		ret = v
	}
	return
}

// for now, doesnt worry if the slot implementation is valid
// it just converts whatever is there
func (dec *Compacter) readSlot(m reader.Map) (ret interface{}, err error) {
	if el, ok := m[reader.ItemValue]; !ok {
		err = errutil.New("missing value while reading slot")
	} else if data, ok := el.(map[string]interface{}); !ok {
		err = errutil.New("expected a map")
	} else {
		ret, err = dec.readItem(data)
	}
	return
}

//
func (dec *Compacter) readSwap(cmd composer.Composer, m reader.Map) (ret reader.Map, err error) {
	if el, ok := m[reader.ItemValue]; !ok {
		err = errutil.New("missing value while reading swap")
	} else if data, ok := el.(map[string]interface{}); !ok {
		err = errutil.New("expected a map")
	} else {
		// try the first key
		for _, swap := range data {
			if swap, ok := swap.(map[string]interface{}); ok {
				swap := reader.Box(swap)
				if i, e := dec.readItem(swap); e != nil {
					err = e
				} else if m, ok := i.(reader.Map); ok {
					ret = m
				} else {
					itemType, itemAt := swap.StrOf(reader.ItemType), reader.At(swap)
					if cmd, ok := dec.cmds[itemType]; !ok {
						err = errutil.Fmt("unknown type %q at %q", itemType, itemAt)
					} else {
						var sig sigBuilder
						spec := cmd.Compose()
						sig.WriteLede(spec)
						n := sig.String()
						ret = reader.Map{n: i}
					}
				}
			}
			break
		}
		if ret == nil && err == nil {
			err = errutil.New("no valid swap data found")
		}
	}
	return
}

// primitives
func (dec *Compacter) readItems(items []interface{}) (ret []interface{}, err error) {
	for _, el := range items {
		if item, ok := el.(map[string]interface{}); !ok {
			err = errutil.New("expected an item, got:", el)
		} else if v, e := dec.readItem(item); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = append(ret, v)
		}
	}
	return
}

func (dec *Compacter) readSlots(slots []interface{}) (ret []interface{}, err error) {
	for _, el := range slots {
		if slot, ok := el.(map[string]interface{}); !ok {
			err = errutil.New("expected a slot, got:", el)
		} else if slat, ok := slot[reader.ItemValue].(map[string]interface{}); !ok {
			err = errutil.New("expected a alot, got:", slot)
		} else if v, e := dec.readItem(slat); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = append(ret, v)
		}
	}
	return
}

// cast inVal to a map, and call setter with contents of "value"
func (dec *Compacter) read(inVal interface{}, setter func(reader.Map, interface{}) error) (err error) {
	if item, ok := inVal.(map[string]interface{}); !ok {
		err = errutil.New("expected an item, got:", inVal)
	} else {
		val := item[reader.ItemValue]
		if e := setter(item, val); e != nil {
			dec.report(reader.At(item), e)
		}
	}
	return
}
