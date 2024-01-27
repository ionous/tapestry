package unblock

import (
	"encoding/json"
	"errors"
	"fmt"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/markup"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

// read blockly data into the passed tapestry command (tree).
// where topBlock is the expected topblock type in the file.... ex. story_file
func Decode(dst typeinfo.Instance, topBlock string, reg Creator, msg json.RawMessage) (err error) {
	var bff File
	if e := json.Unmarshal(msg, &bff); e != nil {
		err = e
	} else if top, ok := bff.FindFirst(topBlock); !ok {
		err = errutil.New("couldnt find story file in block file")
	} else {
		err = DecodeBlock(dst, reg, top)
	}
	return
}

// decode a generic block
// public for testing.
func DecodeBlock(dst typeinfo.Instance, reg Creator, top *BlockInfo) (err error) {
	dec := unblock{reg}
	return dec.decodeBlock(inspect.Walk(dst), top)
}

type unblock struct {
	factory Creator
}

func (un *unblock) decodeBlock(out inspect.Iter, bff *BlockInfo) (err error) {
	if e := readMarkup(out, bff); e != nil {
		err = e
	} else {
		for out.Next() && err == nil {
			f := out.Term()
			termName := upper(f.Name)
			switch t := f.Type; t.(type) {
			default:
				if !f.Private { // private fields dont have typeinfo, and thats okay.
					err = fmt.Errorf("unhandled type %s", t.TypeName())
				}

			// simple values live in bff.fields
			case *typeinfo.Str, *typeinfo.Num:
				if f.Repeats {
					if fields := bff.SliceFields(termName); len(fields) > 0 {
						err = un.decodeList(out.Walk(), fields)
					}
				} else {
					err = decodeField(out, bff, termName)
				}

			// a member that is a flow; its value lives in an input.
			case *typeinfo.Flow:
				if f.Repeats {
					if inputs := bff.SliceInputs(termName); len(inputs) > 0 {
						err = un.decodeSlice(out.Walk(), inputs)
					}
				} else {
					if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
						if input, e := bff.ReadInput(idx); e != nil {
							err = e
						} else {
							err = un.decodeBlock(out.Walk(), input.BlockInfo)
						}
					}
				}

			// a member that fills a slot; its value lives in an input.
			case *typeinfo.Slot:
				if f.Repeats {
					if at := bff.Inputs.FindIndex(termName); at >= 0 {
						err = un.decodeStack(out.Walk(), bff, at)
					} else if inputs := bff.SliceInputs(termName); len(inputs) > 0 {
						err = un.decodeSeries(out.Walk(), inputs)
					}
				} else {
					if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
						if input, e := bff.ReadInput(idx); e != nil {
							err = e
						} else {
							err = un.decodeSlot(out, input.Type, input.BlockInfo)
						}
					}
				}
			}
		}
	}
	return
}

func readMarkup(flow inspect.Iter, bff *BlockInfo) (err error) {
	if c := bff.Icons.Comment; c != nil {
		m := flow.Markup(true)
		lines := strings.FieldsFunc(c.Text, func(r rune) bool {
			const newline = '\n'
			return r == newline
		})
		switch len(lines) {
		case 0:
			m[markup.Comment] = ""
		case 1:
			m[markup.Comment] = lines[0]
		default:
			m[markup.Comment] = lines
		}
	}
	return
}

// a stack is a repeating slot
func (un *unblock) decodeStack(out inspect.Iter, bff *BlockInfo, idx int) (err error) {
	if input, e := bff.ReadInput(idx); e != nil {
		err = e
	} else {
		// tbd: sink them into a flat list during count?
		// re: +1: all of the blocks connected to the input, plus the input block itself.
		out.Resize(1 + input.CountNext())
		for next := input; next.BlockInfo != nil; next = next.GetNext() {
			// the typename we want is (munged) in the block file
			if typeName, ok := unstackName(next.Type); !ok {
				err = fmt.Errorf("couldnt unstack %q", next.Type)
				break
			} else {
				out.Next()
				if e := un.decodeSlot(out, typeName, next.BlockInfo); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// repeating non-stacking slots
func (un *unblock) decodeSeries(out inspect.Iter, inputs js.MapSlice) (err error) {
	out.Resize(len(inputs))
	for _, el := range inputs {
		out.Next()
		if input, e := readInput(el); e != nil {
			err = e
			break
		} else if un.decodeSlot(out, input.Type, input.BlockInfo); e != nil {
			err = e
			break
		}
	}
	return
}

func (un *unblock) decodeSlot(out inspect.Iter, inputType string, next *BlockInfo) (err error) {
	if e := un.fillSlot(out, inputType); e != nil {
		err = e
	} else if slot := out.Walk(); !slot.Next() {
		err = errors.New("slot empty after filling it?")
	} else if e := un.decodeBlock(slot.Walk(), next); e != nil {
		err = e
	}
	return
}

// create a blank command to fill the targeted slot.
func (un *unblock) fillSlot(out inspect.Iter, typeName string) (err error) {
	if rptr, ok := un.factory.NewType(typeName); !ok {
		err = errutil.Fmt("couldn't create %q", typeName)
	} else {
		out.RawValue().Set(r.ValueOf(rptr))
	}
	return
}

// repeating flows
func (un *unblock) decodeSlice(out inspect.Iter, inputs js.MapSlice) (err error) {
	out.Resize(len(inputs))
	for _, el := range inputs {
		out.Next()
		if input, e := readInput(el); e != nil {
			err = e
			break
		} else if e := un.decodeBlock(out.Walk(), input.BlockInfo); e != nil {
			err = e
			break
		}
	}
	return
}

// an array of primitives is a list of fields.
func (un *unblock) decodeList(out inspect.Iter, fields js.MapSlice) (err error) {
	out.Resize(len(fields))
	for _, el := range fields {
		out.Next()
		if v, e := readValue(el); e != nil {
			err = e
			break
		} else if !out.SetValue(v) {
			err = errutil.Fmt("couldn't assign from %T", v)
		}
	}
	return
}

// simple values live in bff.fields
func decodeField(out inspect.Iter, bff *BlockInfo, fieldName string) (err error) {
	if field, ok := bff.Fields.Find(fieldName); ok {
		var value any
		if e := json.Unmarshal(field.Msg, &value); e != nil {
			err = e
		} else if !out.SetValue(value) {
			err = fmt.Errorf("couldnt store %T", value)
		}
	}
	return
}
