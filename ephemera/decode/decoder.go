package decode

import (
	r "reflect"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
	"github.com/ionous/iffy/lang"
)

// ReadRet is similar to reader.ReadMap, except it returns a value.
type ReadRet func(reader.Map) (interface{}, error)

type cmdRec struct {
	name         string
	elem         r.Type
	customReader ReadRet
}

type Override struct {
	Spec     composer.Slat
	Callback ReadRet
}

// Decoder reads programs from json.
type Decoder struct {
	source     string
	cmds       map[string]cmdRec
	issueFn    IssueReport
	IssueCount int
}

func NewDecoder() *Decoder {
	reportNothing := func(reader.Position, error) {}
	return NewDecoderReporter("decoder", reportNothing)
}

func (dec *Decoder) AddCallbacks(overrides []Override) {
	for _, n := range overrides {
		dec.AddCallback(n.Spec, n.Callback)
	}
}

// AddCallback registers a command parser.
func (dec *Decoder) AddCallback(cmd composer.Slat, cb ReadRet) {
	n := specName(cmd)
	if was, exists := dec.cmds[n]; exists && was.customReader != nil {
		panic(errutil.Fmt("conflicting name for spec %q %q!=%T", n, was.elem, cmd))
	} else {
		elem := r.TypeOf(cmd).Elem()
		dec.cmds[n] = cmdRec{n, elem, cb}
	}
}

func specName(cmd composer.Slat) (ret string) {
	spec := cmd.Compose()
	if n := spec.Name; len(n) > 0 {
		ret = n
	} else {
		elem := r.TypeOf(cmd).Elem()
		ret = lang.Underscore(elem.Name())
	}
	return
}

// AddDefaultCallbacks registers default command parsers.
func (dec *Decoder) AddDefaultCallbacks(slats []composer.Slat) {
	for _, cmd := range slats {
		dec.AddCallback(cmd, nil)
	}
}

func (dec *Decoder) ReadSpec(m reader.Map) (ret interface{}, err error) {
	if rptr, e := dec.read(m); e != nil {
		err = e
	} else {
		ret = rptr.Interface()
	}
	return
}

// ReadProg attempts to parse the passed json data as a golang program.
func (dec *Decoder) ReadProg(m reader.Map, outPtr interface{}) (err error) {
	if rptr, e := dec.read(m); e != nil {
		err = e
	} else {
		out := r.ValueOf(outPtr).Elem()
		outType := out.Type()
		if rtype := rptr.Type(); !rtype.AssignableTo(outType) {
			err = errutil.New("incompatible types", rtype.String(), "not assignable to", outType.String())
		} else {
			out.Set(rptr)
		}
	}
	return
}

func (dec *Decoder) read(m reader.Map) (ret r.Value, err error) {
	itemValue, itemType := m, m.StrOf(reader.ItemType)
	if cmd, ok := dec.cmds[itemType]; !ok {
		err = errutil.Fmt("unknown type %q with reading a program at %s", itemType, reader.At(m))
	} else {
		ret, err = dec.readNew(cmd, itemValue)
	}
	return
}

// m is the contents of slotType is a concrete command; returns a pointer to the command
func (dec *Decoder) readNew(cmd cmdRec, m reader.Map) (ret r.Value, err error) {
	if read := cmd.customReader; read != nil {
		if res, e := read(m); e != nil {
			err = e
		} else {
			ret = r.ValueOf(res)
		}
	} else if cmd.elem.Kind() != r.Struct {
		err = errutil.New("expected a struct", cmd.name, "is a", cmd.elem.String())
	} else {
		ptr := r.New(cmd.elem)
		dec.ReadFields(reader.At(m), ptr.Elem(), m.MapOf(reader.ItemValue))
		ret = ptr
	}
	return
}

func (dec *Decoder) ReadFields(at string, out r.Value, in reader.Map) {
	var fields []string
	export.WalkProperties(out.Type(), func(f *r.StructField, path []int) (done bool) {
		token := export.Tokenize(f)
		fields = append(fields, token)
		// we report on missing properties below.
		if inVal, ok := in[token]; !ok {
			// log only if the field is required. not optional.
			if t := tag.ReadTag(f.Tag); !t.Exists("internal") && !t.Exists("optional") {
				// and even then only if its a fixed field
				if f.Type.Kind() != r.Ptr {
					dec.report(at, errutil.Fmt("missing %q", token))
				}
			}
		} else {
			outAt := out.FieldByIndex(path)
			if e := dec.importValue(outAt, inVal); e != nil {
				dec.report(at, errutil.New("error processing field", out.Type().String(), f.Name, e))
			}
		}
		return
	})

	// walk keys of json dictionary:
	for token, _ := range in {
		i, cnt := 0, len(fields)
		for ; i < cnt; i++ {
			if token == fields[i] {
				break
			}
		}
		if i == cnt {
			dec.report(at, errutil.Fmt("unprocessed %q", token))
		}
	}
}

// returns a ptr r.Value
func (dec *Decoder) importSlot(m reader.Map, slotType r.Type) (ret r.Value, err error) {
	itemValue, itemType := m, m.StrOf(reader.ItemType)
	slotName := slotType.Name() // here for debugging; ex. "Comparator"
	if cmd, ok := dec.cmds[itemType]; !ok {
		err = errutil.Fmt("unknown type %q while importing slot %q", itemType, slotName)
	} else if rptr, e := dec.readNew(cmd, itemValue); e != nil {
		err = e
	} else if rtype := rptr.Type(); !rtype.AssignableTo(slotType) {
		err = errutil.New("incompatible types", rtype.String(), "not assignable to", slotName)
	} else {
		ret = rptr
	}
	return
}

func (dec *Decoder) importValue(outAt r.Value, inVal interface{}) (err error) {
	switch outType := outAt.Type(); outType.Kind() {
	case r.Float32, r.Float64:
		err = dec.unpack(inVal, func(v interface{}) (err error) {
			// float64, for JSON numbers
			if n, ok := v.(float64); !ok {
				err = errutil.Fmt("expected a number, have %T", v)
			} else {
				outAt.SetFloat(n)
			}
			return
		})
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		err = dec.unpack(inVal, func(v interface{}) (err error) {
			// float64, for JSON numbers
			if n, ok := v.(float64); !ok {
				err = errutil.New("expected a number")
			} else {
				outAt.SetInt(int64(n))
			}
			return
		})

	case r.Bool:
		// fix? boolean values are stored as enumerations
		err = dec.unpack(inVal, func(v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a string")
			} else {
				outAt.SetBool(str == "$TRUE") // only need to set true: false is the zero value.
			}
			return
		})

	case r.String:
		err = dec.unpack(inVal, func(v interface{}) (err error) {
			// string, for JSON strings
			if str, ok := v.(string); !ok {
				err = errutil.New("expected a string")
			} else {
				outAt.SetString(str)
			}
			return
		})

	case r.Ptr:
		// see if its an optional value.
		ptr := r.New(outAt.Type().Elem())
		if e := dec.importValue(ptr.Elem(), inVal); e != nil {
			err = e
		} else {
			outAt.Set(ptr)
		}

	case r.Struct:
		// b/c of the way optional values are specified,
		// going from r.Struct is easier than from r.Ptr.
		switch spec := outAt.Addr().Interface().(type) {
		case StrType:
			if e := dec.unpack(inVal, func(v interface{}) (err error) {
				if str, ok := v.(string); !ok {
					err = errutil.New("value not a slat")
				} else {
					// validate choice.
					if closed, vs := spec.Choices(); closed {
						var found bool
						for _, choice := range vs {
							if str == "$"+strings.ToUpper(choice) {
								str = choice
								found = true
								break
							}
						}
						if !found {
							err = errutil.New("unknown string", str)
						}
					}
					if err == nil {
						outAt.Field(outAt.NumField() - 1).SetString(str)
					}
				}
				return
			}); e != nil {
				err = e
			}

		case NumType:
			if e := dec.unpack(inVal, func(v interface{}) (err error) {
				if num, ok := v.(float64); !ok {
					err = errutil.New("value not a slat")
				} else {
					// validate choice; fix: tolerance?
					if closed, vs := spec.Choices(); closed {
						var found bool
						for _, choice := range vs {
							if num == choice {
								found = true
								break
							}
						}
						if !found {
							err = errutil.New("unknown value", num)
						}
					}
					if err == nil {
						// handle conversion b/t floats and ints of different widths
						tgt := outAt.Field(outAt.NumField() - 1)
						v := r.ValueOf(num).Convert(tgt.Type())
						tgt.Set(v)
					}
				}
				return
			}); e != nil {
				err = e
			}

		case SwapType:
			if e := dec.unpack(inVal, func(v interface{}) (err error) {
				if data, ok := v.(map[string]interface{}); !ok {
					err = errutil.New("value not a slat")
				} else {
					found := false
					for k, typePtr := range spec.Choices() {
						token := "$" + strings.ToUpper(k)
						if contents, ok := data[token]; ok {
							ptr := r.New(r.TypeOf(typePtr).Elem())
							if e := dec.importValue(ptr.Elem(), contents); e != nil {
								err = e
							} else {
								outAt.Field(outAt.NumField() - 1).Set(ptr)
							}
							found = true
							break
						}
					}
					if !found {
						err = errutil.New("no valid swap data found")
					}
				}
				return
			}); e != nil {
				err = e
			}

		default:
			if e := dec.unpack(inVal, func(v interface{}) (err error) {
				if slot, ok := v.(map[string]interface{}); !ok {
					err = errutil.New("value not a slat")
				} else {
					dec.ReadFields(reader.At(slot), outAt, reader.Map(slot))
				}
				return
			}); e != nil {
				err = e
			}
		}

	case r.Interface:
		if e := dec.unpack(inVal, func(v interface{}) (err error) {
			// map[string]interface{}, for JSON objects
			if slot, ok := v.(map[string]interface{}); !ok {
				err = errutil.New("value not a slot")
			} else if newVal, e := dec.importSlot(slot, outAt.Type()); e != nil {
				dec.report(reader.At(slot), e)
			} else {
				outAt.Set(newVal)
			}
			return
		}); e != nil {
			err = e
		}

	case r.Slice:
		// []interface{}, for JSON arrays
		if items, ok := inVal.([]interface{}); !ok {
			err = errutil.New("expected a slice")
		} else {
			elType := outType.Elem()
			if slice := outAt; len(items) > 0 {
				for _, item := range items {
					if k := elType.Kind(); k != r.Interface {
						el := r.New(elType).Elem()
						if e := dec.importValue(el, item); e != nil {
							err = errutil.Append(err, e)
						} else {
							slice = r.Append(slice, el)
						}
					} else {
						// note: this skips over the slot itself ( ex execute )
						if e := dec.unpack(item, func(v interface{}) (err error) {
							// map[string]interface{}, for JSON objects
							if itemData, ok := v.(map[string]interface{}); !ok {
								// execute has some single nulls sometimes;
								if v != nil {
									err = errutil.Fmt("item data not a slot %T", itemData)
								}

							} else if v, e := dec.importSlot(itemData, elType); e != nil {
								err = e // elType is ex. *story.Paragraph; itemData has a member $STORY_STATEMENT
							} else {
								slice = r.Append(slice, v)
							}
							return
						}); e != nil {
							err = errutil.Append(err, e)
						}
					}
				}
				outAt.Set(slice)
			}
		}
	}
	return
}

// cast inVal to a map, and call setter with contents of "value"
func (dec *Decoder) unpack(inVal interface{}, setter func(interface{}) error) (err error) {
	if item, ok := inVal.(map[string]interface{}); !ok {
		err = errutil.New("expected an item, got:", inVal)
	} else {
		val := item[reader.ItemValue]
		if e := setter(val); e != nil {
			dec.report(reader.At(item), e)
		}
	}
	return
}
