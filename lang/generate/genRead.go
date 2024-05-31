package generate

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
)

// ex. "Spec:groups:with flow"
func readSpec(groupName string, out *groupContent, msg compact.Message) (err error) {
	if msg.Name != "spec" {
		err = fmt.Errorf("expected a spec, got %s", msg.Name)
	} else if name, e := parseString("spec name", msg.Args[0], ""); e != nil {
		err = e
	} else {
		spec := specData{name, groupName, msg.Markup}
		var slots []string
		for i, cnt := 1, len(msg.Labels); i < cnt && err == nil; i++ {
			str := msg.Labels[i]
			arg := msg.Args[i]
			switch str {
			case "slots":
				// a string or array
				slots, err = parseStrings(arg)
			case "groups":
				// a string or array
				if _, e := parseStrings(arg); e != nil {
					err = e
				}
			default:
				if inner, e := decode.ParseMessage(arg); e != nil {
					err = e
				} else {
					switch str {
					case "with flow":
						if d, e := readFlow(spec, inner, slots); e != nil {
							err = e
						} else {
							out.Flow = append(out.Flow, d)
							out.Reg = out.Reg.addFlow(d)
						}
					case "with group":
						err = readGroupContent(groupName, out, inner)

					case "with slot":
						if d, e := readSlot(spec, inner); e != nil {
							err = e
						} else {
							out.Slot = append(out.Slot, d)
						}
					case "with str":
						if d, e := readStr(spec, inner); e != nil {
							err = e
						} else {
							out.Str = append(out.Str, d)
						}
					case "with num":
						if d, e := readNum(spec, inner); e != nil {
							err = e
						} else {
							out.Num = append(out.Num, d)
						}
					default:
						// fix: the specs should be switched to regular slots
						// might consider ripping off the out "TypeSpec" part
						// especially because so much of tapestry only supports slots for flows
						err = fmt.Errorf("swap no longer supported")
					}
				}
			}
		}
	}
	return
}

// fix? this doesnt handle subgroups in any good way.
func readGroupContent(groupName string, out *groupContent, msg compact.Message) (err error) {
	if msg.Key != "Group contains:" {
		err = fmt.Errorf("expected group definition, have %s", msg.Key)
	} else if msgs, e := parseMessages(msg.Args[0]); e != nil {
		err = e
	} else {
		for _, inner := range msgs {
			if e := readSpec(groupName, out, inner); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func readFlow(spec specData, msg compact.Message, slots []string) (ret flowData, err error) {
	mm := messageMap(msg)
	if lede, e := mm.GetString("", spec.Name); e != nil {
		err = e
	} else if uses, ok := mm["uses"]; !ok {
		err = errors.New("missing flow terms")
	} else if terms, e := parseMessages(uses); e != nil {
		err = e
	} else {
		all := make([]termData, len(terms))
		for i, td := range terms {
			term := messageMap(td)
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
	if err != nil {
		if pos, ok := msg.Markup["pos"].([]int); ok {
			err = fmt.Errorf("%w at line %d col %d", err, pos[1], pos[0])
		}
	}
	return
}

// generate a slot named "n"
func readSlot(spec specData, msg compact.Message) (ret slotData, err error) {
	mm := messageMap(msg)
	if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if cnt := len(uses); cnt > 0 {
		err = errors.New("expected no constraints for slot")
	} else {
		ret = slotData{spec}
	}
	return
}

func readStr(spec specData, msg compact.Message) (ret strData, err error) {
	mm := messageMap(msg)
	if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if ex, e := mm.GetBool("exclusively"); e != nil {
		err = e
	} else if !ex {
		ret, err = readSimpleStr(spec, uses)
	} else {
		ret, err = readEnum(spec, uses)
	}
	return
}

func readSimpleStr(spec specData, uses []compact.Message) (ret strData, err error) {
	if cnt := len(uses); cnt > 0 {
		err = errors.New("expected no constraints for simple string")
	} else {
		ret = strData{
			specData:       spec,
			Options:        nil,
			OptionComments: nil,
		}
	}
	return
}

func readEnum(spec specData, uses []compact.Message) (ret strData, err error) {
	if cnt := len(uses); cnt == 0 {
		err = errors.New("expected string choices")
	} else {
		options := make([]string, cnt)
		var comments []string
		for i, msg := range uses {
			if msg.Key != "Option:" || len(msg.Args) != 1 {
				err = fmt.Errorf("option %d has unexpected option %q", i, msg.Key)
				break
			} else if str, ok := msg.Args[0].(string); !ok {
				err = fmt.Errorf("option %d has unexpected value %q", i, msg.Key)
				break
			} else {
				options[i] = str
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
				Options:        options,
				OptionComments: comments,
			}
		}
	}
	return
}

func readNum(spec specData, msg compact.Message) (ret numData, err error) {
	mm := messageMap(msg)
	if _, ok := mm["exclusively"]; ok {
		err = errors.New("exclusive nums not supported")
	} else if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if len(uses) != 0 {
		err = errors.New("expected no numeric constraints")
	} else {
		ret = numData{spec}
	}
	return
}
