package generate

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
)

// ex. "Spec:groups:with flow"
func readSpec(out *groupContent, msg compact.Message) (err error) {
	if msg.Name != "spec" {
		err = fmt.Errorf("expected a spec, got %s", msg.Name)
	} else if name, e := parseString("spec name", msg.Args[0], ""); e != nil {
		err = e
	} else {
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
				// FIX: add to current groups (push /pop)
			default:
				if inner, e := decode.ParseMessage(arg); e != nil {
					err = e
				} else {
					switch str {
					case "with flow":
						if d, e := readFlow(name, inner, slots); e != nil {
							err = e
						} else {
							d.Markup = msg.Markup
							out.Flow = append(out.Flow, d)
						}
					case "with group":
						err = readGroup(out, name, inner)

					case "with slot":
						if d, e := readSlot(name, inner); e != nil {
							err = e
						} else {
							d.Markup = msg.Markup
							out.Slot = append(out.Slot, d)
						}
					case "with str":
						if d, e := readStr(name, inner); e != nil {
							err = e
						} else {
							d.Markup = msg.Markup
							out.Str = append(out.Str, d)
						}
					case "with num":
						if d, e := readNum(name, inner); e != nil {
							err = e
						} else {
							d.Markup = msg.Markup
							out.Num = append(out.Num, d)
						}
					default:
						// which of course is ironic, because the specs currently use slots
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

func readGroup(out *groupContent, n string, msg compact.Message) (err error) {
	if msg.Key != "Group contains:" {
		err = fmt.Errorf("expected group definition, have %s", msg.Key)
	} else if msgs, e := parseMessages(msg.Args[0]); e != nil {
		err = e
	} else {
		// groupNames.AddName(n)
		for _, inner := range msgs {
			if e := readSpec(out, inner); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func readFlow(n string, msg compact.Message, slots []string) (ret flowData, err error) {
	mm := messageMap(msg)
	if lede, e := mm.GetString("", n); e != nil {
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
				Name:  n,
				Lede:  lede,
				Slots: slots,
				Terms: all,
			}
		}
	}
	return
}

// generate a slot named "n"
func readSlot(n string, msg compact.Message) (ret slotData, err error) {
	mm := messageMap(msg)
	if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if cnt := len(uses); cnt > 0 {
		err = errors.New("expected no constraints for slot")
	} else {
		ret = slotData{Name: n}
	}
	return
}

func readStr(n string, msg compact.Message) (ret strData, err error) {
	mm := messageMap(msg)
	if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if ex, e := mm.GetBool("exclusively"); e != nil {
		err = e
	} else if !ex {
		ret, err = readSimpleStr(n, uses)
	} else {
		ret, err = readEnum(n, uses)
	}
	return
}

func readSimpleStr(n string, uses []compact.Message) (ret strData, err error) {
	if cnt := len(uses); cnt > 0 {
		err = errors.New("expected no constraints for simple string")
	} else {
		ret = strData{Name: n}
	}
	return
}

func readEnum(n string, uses []compact.Message) (ret strData, err error) {
	if cnt := len(uses); cnt == 0 {
		err = errors.New("expected string choices")
	} else {
		options := make([]string, cnt)
		for i, ex := range uses {
			if ex.Key != "Option:" || len(ex.Args) != 1 {
				err = fmt.Errorf("unexpected option %q", ex.Key)
				break
			} else if str, ok := ex.Args[0].(string); !ok {
				err = fmt.Errorf("unexpected value%q", ex.Key)
				break
			} else {
				options[i] = str
			}
		}
		if err == nil {
			ret = strData{Name: n, Options: options}
		}
	}
	return
}

func readNum(n string, msg compact.Message) (ret numData, err error) {
	mm := messageMap(msg)
	if _, ok := mm["exclusively"]; ok {
		err = errors.New("exclusive nums not supported")
	} else if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if len(uses) != 0 {
		err = errors.New("expected no numeric constraints")
	} else {
		ret = numData{Name: n}
	}
	return
}
