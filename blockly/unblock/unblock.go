package unblock

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/markup"
	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

type unblock struct {
	factory TypeCreator
}

func (un *unblock) decodeBlock(out walk.Walker, bff *BlockInfo) (err error) {

	for it := out.Walk(); it.Next() && err == nil; {
		f := it.Field()
		termName := upper(f.Name())
		// a member that repeats:
		// could be a slice of a specific flow of a series of slots ( of numbered inputs ),
		// a stack ( starting in input, continuing in next ), a list ( of numbered fields )
		// we could return true/false depending on whether the block file has data
		// but first we have to find out where that data lives
		if f.Repeats() { // SliceBlock
			if fields := bff.SliceFields(termName); len(fields) > 0 {
				un.decodeList(it.Walk(), fields)

			} else if inputs := bff.SliceInputs(termName); len(inputs) > 0 {
				err = un.decodeSeriesSlice(out.Walk(), f, inputs)

			} else if at := bff.Inputs.FindIndex(termName); at >= 0 {

				if input, e := bff.ReadInput(at); e != nil {
					log.Println(e)
				} else {
					// might be nicer if count could grow, rather than counting in advance
					// could also sink them into a flat list as we count.
					cnt := 1 + input.CountNext() // all of the next blocks connected to the input, plus the input block itself.
					out.Resize(cnt)
					err = un.decodeStack(out, input.BlockInfo)
				}
			} else {
				switch t := f.SpecType(); t {
				default:
					err = fmt.Errorf("unhandled type %s", t)

					// a member that is a flow; its value lives in an input.
				case walk.Flow:
					if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
						if input, e := bff.ReadInput(idx); e != nil {
							log.Println(e)
						} else {
							err = un.decodeBlock(out.Walk(), input.BlockInfo)
						}
					}

					// a member that fills a slot; its value lives in an input.
				case walk.Slot:
					if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
						if input, e := bff.ReadInput(idx); e != nil {
							err = e
						} else if e := un.fillSlot(out, input.Type); e != nil {
							err = e
						} else {
							err = un.decodeBlock(out.Walk(), input.BlockInfo)
						}
					}
				// case :
				// a member that's a swap uses both a field and an input.
				// for simple values ( strs in swaps ) there will be a faux block type for that input.

				case walk.Swap:
					// the field holds a combo box with swap's choice
					if idx := bff.Fields.FindIndex(termName); idx >= 0 {
						var choice string
						if e := json.Unmarshal(bff.Fields[idx].Msg, &choice); e != nil {
							err = e
						} else if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
							// the input hold the value of the swap's block
							if input, e := bff.ReadInput(idx); e != nil {
								err = e
							} else {
								err = un.decodeSwap(out, choice, input.BlockInfo)
							}
						}
					}

					// simple values live in bff.fields
				case walk.Str, walk.Value:
					var value any
					if field, ok := bff.Fields.Find(termName); !ok {
						err = jsn.Missing
					} else if e := json.Unmarshal(field.Msg, &value); e != nil {
						err = e
					} else {
						err = storeValue(out, f, value)
					}
				}
			}
		}
	}

	panic("readMarkup")

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
func (un *unblock) decodeStack(out walk.Walker, next *BlockInfo) (err error) {
	// the typename we want is (munged) in the block file
	if typeName, ok := unstackName(next.Type); !ok {
		err = fmt.Errorf("couldnt unstack %q", next.Type)
	} else if e := un.fillSlot(out, typeName); e != nil {
		err = e
	} else {
		err = un.decodeBlock(out.Walk(), next)
	}
	return
}

// read the insides of a swap:
// it could be a flow filling the input....
// or a fake block wrapping a primitive value ( a "standalone" )
func (un *unblock) decodeSwap(out walk.Walker, choice string, bff *BlockInfo) (err error) {
	// //if !swap.SetSwap(choice) {

	// // see: block.newSwap & shape.writeStandalone:
	// // the field name is the name of the spec type )
	// field := strings.ToUpper(typeName)
	// if idx := bff.Fields.FindIndex(field); idx < 0 {
	// 	err = jsn.Missing // the block might be missing, and that's okay.
	// } else {
	// 	field := bff.Fields[idx]
	// 	err = storeValue(pv, field.Msg) // pv is the destination
	// }
	// return next
	panic("huh")
}

// fix: refactor
func (un *unblock) decodeSeriesSlice(out walk.Walker, f walk.Field, inputs js.MapSlice) (err error) {
	switch f.SpecType() {
	case walk.Slot:
		err = un.decodeSeries(out, inputs)
	case walk.Flow:
		err = un.decodeSlice(out, inputs)
	default:
		err = errors.New("unexpected slice")
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
		} else if e := un.fillSlot(out, input.Type); e != nil {
			err = e
			break
		} else if e := un.decodeBlock(out.Walk(), input.BlockInfo); e != nil {
			err = e
			break
		}
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

func (un *unblock) fillSlot(out walk.Walker, typeName string) (err error) {
	if rptr, ok := un.factory.NewType(typeName); !ok {
		err = errutil.Fmt("couldn't create %q", typeName)
	} else {
		out.Value().Set(r.ValueOf(rptr))
	}
	return
}

// an array of primitives is a list of fields .
func (un *unblock) decodeList(out walk.Walker, fields js.MapSlice) (err error) {
	out.Resize(len(fields))
	for _, el := range fields {
		out.Next()
		if v, e := readValue(el); e != nil {
			err = e
			break
		} else {
			// fix: might need to convert via the regular decoder.
			// especially because this can panic
			out.Value().Set(r.ValueOf(v))
		}
	}
	return
}

func storeValue(out walk.Walker, f walk.Field, v any) (err error) {
	panic("")
}
