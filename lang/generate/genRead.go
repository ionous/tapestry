package generate

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// ex. Spec:requires:contains:
func readSpec(group string, mm MessageMap) (ret groupContent, err error) {
	if req, e := mm.GetStrings("requires"); e != nil {
		err = e
	} else if contents, e := mm.GetMessages("contains"); e != nil {
		err = e
	} else {
		ret.Requires = req
		for _, msg := range contents {
			if e := readEntry(group, &ret, msg); e != nil {
				if pos, ok := msg.Markup[compact.Position].([]int); ok {
					err = fmt.Errorf("%w at line %d col %d", e, pos[1], pos[0])
					break
				}
			}
		}
	}
	return
}

func readEntry(group string, out *groupContent, msg compact.Message) (err error) {
	mm := MakeMessageMap(msg)
	if name, e := mm.GetString("", ""); e != nil {
		err = e
	} else {
		spec := specData{name, group, msg.Markup}
		switch msg.Lede {
		case "flow":
			if d, e := readFlow(spec, mm); e != nil {
				err = e
			} else {
				out.Flow = append(out.Flow, d)
				out.Reg = out.Reg.addFlow(d)
			}
		case "slot":
			if slot, e := readSlot(spec, mm); e != nil {
				err = e
			} else {
				out.Slot = append(out.Slot, slot)
			}
		case "str":
			if d, e := readStr(spec, mm); e != nil {
				err = e
			} else {
				out.Str = append(out.Str, d)
			}
		case "num":
			if d, e := readNum(spec, mm); e != nil {
				err = e
			} else {
				out.Num = append(out.Num, d)
			}
		default:
			err = fmt.Errorf("unknown message %q", name)
		}
	}
	return
}

// ex. Flow:slots:lede:terms:
// only the first term is required: the name.
func readFlow(spec specData, mm MessageMap) (ret flowData, err error) {
	if slots, e := mm.GetStrings("slots"); e != nil {
		err = e
	} else if lede, e := mm.GetString("lede", spec.Name); e != nil {
		err = e
	} else if terms, e := mm.GetMessages("terms"); e != nil {
		err = e
	} else {
		all := make([]termData, len(terms))
		for i, td := range terms {
			term := MakeMessageMap(td)
			if label, e := term.GetString("", ""); e != nil {
				err = e
				break
			} else if name, e := term.GetString("name", label); e != nil {
				err = e
				break
			} else if termType, e := term.GetString("type", name); e != nil {
				err = e
				break
			} else if private, e := term.GetBool("private"); e != nil {
				err = e
				break
			} else if optional, e := term.GetBool("optional"); e != nil {
				err = e
				break
			} else if repeats, e := term.GetBool("repeats"); e != nil {
				err = e
				break
			} else {
				all[i] = termData{
					name,
					label,
					termType,
					private,
					optional,
					repeats,
					td.Markup,
				}
			}
		}
		if err == nil {
			ret = flowData{
				specData: spec,
				Lede:     lede,
				Slots:    slots,
				Terms:    all,
			}
		}
	}
	return
}

// ex. Slot: "slotname"
func readSlot(spec specData, mm MessageMap) (ret slotData, _ error) {
	ret = slotData{specData: spec}
	return
}

// ex. Num:
func readNum(spec specData, mm MessageMap) (ret numData, _ error) {
	ret = numData{spec}
	return
}

// ex. Str:options:
func readStr(spec specData, mm MessageMap) (ret strData, err error) {
	if options, e := mm.GetMessages("options"); e != nil {
		err = e
	} else if cnt := len(options); cnt == 0 {
		ret = strData{specData: spec}
	} else {
		names := make([]string, cnt)
		var comments []string
		for i, msg := range options {
			if msg.Key != "Option:" || len(msg.Args) != 1 {
				err = fmt.Errorf("option %d has unexpected option %q", i, msg.Key)
				break
			} else if str, ok := msg.Args[0].(string); !ok {
				err = fmt.Errorf("option %d has unexpected value %q", i, msg.Key)
				break
			} else {
				names[i] = str
				if lines, e := compact.ExtractComment(msg.Markup); e != nil {
					err = fmt.Errorf("option %d has %w", i, e)
					break
				} else if len(lines) > 0 {
					if comments == nil {
						comments = make([]string, cnt)
					}
					// fix: newlines in comments; i'd like if comments were normalized to a single line
					// with newlines as literal \n(s) and appropriate escaping
					// maybe lift comments out of markup for types, since we handle them explicitly.
					comments[i] = strings.Join(lines, " ")
				}
			}
		}
		if err == nil {
			ret = strData{
				specData:       spec,
				Options:        names,
				OptionComments: comments,
			}
		}
	}
	return
}
