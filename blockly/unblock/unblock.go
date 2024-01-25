package unblock

import (
	"encoding/json"
	"errors"
	"fmt"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/markup"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

// where topBlock is the expected topblock type in the file.... ex. story_file
func Decode(dst typeinfo.Inspector, topBlock string, reg TypeCreator, msg json.RawMessage) (err error) {
	var bff File
	if e := json.Unmarshal(msg, &bff); e != nil {
		err = e
	} else if top, ok := bff.FindFirst(topBlock); !ok {
		err = errutil.New("couldnt find story file in block file")
	} else {
		err = DecodeBlock(r.ValueOf(dst), reg, top)
	}
	return
}

// where topBlock is the expected topblock type in the file.... ex. story_file
func DecodeBlock(dst r.Value, reg TypeCreator, top *BlockInfo) (err error) {
	dec := unblock{reg}
	return dec.decodeBlock(walk.Walk(dst), top)
}

type unblock struct {
	factory TypeCreator
}

func (un *unblock) decodeBlock(out walk.Walker, bff *BlockInfo) (err error) {
	if e := readMarkup(out, bff); e != nil {
		err = e
	} else {
		for out.Next() && err == nil {
			f := out.Field()
			termName := upper(f.Name())
			switch t := f.SpecType(); t {
			default:
				err = fmt.Errorf("unhandled type %s", t)

			// simple values live in bff.fields
			case walk.Str, walk.Value:
				if f.Repeats() {
					if fields := bff.SliceFields(termName); len(fields) > 0 {
						err = un.decodeList(out.Walk(), fields)
					}
				} else {
					err = decodeField(out, bff, termName)
				}

			// a member that is a flow; its value lives in an input.
			case walk.Flow:
				if f.Repeats() {
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
			case walk.Slot:
				if f.Repeats() {
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

func readMarkup(flow walk.Walker, bff *BlockInfo) (err error) {
	if c := bff.Icons.Comment; c != nil {
		m := make(map[string]any)
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
		flow.Markup().Set(r.ValueOf(m))
	}
	return
}

// a stack is a repeating slot
func (un *unblock) decodeStack(out walk.Walker, bff *BlockInfo, idx int) (err error) {
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
func (un *unblock) decodeSeries(out walk.Walker, inputs js.MapSlice) (err error) {
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

func (un *unblock) decodeSlot(out walk.Walker, inputType string, next *BlockInfo) (err error) {
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
func (un *unblock) fillSlot(out walk.Walker, typeName string) (err error) {
	if rptr, ok := un.factory.NewType(typeName); !ok {
		err = errutil.Fmt("couldn't create %q", typeName)
	} else {
		out.Value().Set(r.ValueOf(rptr))
	}
	return
}

// repeating flows
func (un *unblock) decodeSlice(out walk.Walker, inputs js.MapSlice) (err error) {
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
func (un *unblock) decodeList(out walk.Walker, fields js.MapSlice) (err error) {
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
func decodeField(out walk.Walker, bff *BlockInfo, fieldName string) (err error) {
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
