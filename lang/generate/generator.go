package generate

import (
	"errors"
	"fmt"
	"io"
	"text/template"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

type generator struct {
	w      io.Writer
	tmp    *template.Template
	groups Groups
}

// // in: some spec data in "tell" format.
// func generate(out io.Writer, in io.Reader) (err error) {
// 	var m map[string]any
// 	if e := files.ReadTell(in, &m); e != nil {
// 		err = e
// 	} else if msg, e := decode.DecodeMessage(m); e != nil {
// 		err = e
// 	} else {
// 		err = generateMsg(out, msg)
// 	}
// 	return
// }

func generateMsg(w io.Writer, msg compact.Message) (err error) {
	var pack packageHelper
	if tmp, e := genTemplates(&pack); e != nil {
		err = e
	} else {
		var gc groupContent
		gen := generator{w: w, tmp: tmp}
		if e := gen.readSpec(&gc, msg); e != nil {
			err = e
		} else {
			pack.groups = append(pack.groups, groupData{"", gc})
			err = gen.write(gc)
		}
	}
	return
}

func (gen *generator) write(gc groupContent) (err error) {
	for i, cnt := 0, len(gc.Slot); i < cnt && err == nil; i++ {
		n := gc.Slot[i]
		err = gen.runTemplate("slot", n)
	}
	for i, cnt := 0, len(gc.Flow); i < cnt && err == nil; i++ {
		n := gc.Flow[i]
		err = gen.runTemplate("flow", n)
	}
	for i, cnt := 0, len(gc.Str); i < cnt && err == nil; i++ {
		n := gc.Str[i].(strData)
		if len(n.Options) > 0 {
			err = gen.runTemplate("enum", n)
		} else {
			err = gen.runTemplate("string", n)
		}
	}
	for i, cnt := 0, len(gc.Num); i < cnt && err == nil; i++ {
		n := gc.Num[i]
		err = gen.runTemplate("num", n)
	}
	return
}

// ex. "Spec:groups:with flow"
func (gen *generator) readSpec(out *groupContent, msg compact.Message) (err error) {
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
				if inner, e := parseMessage(arg); e != nil {
					err = e
				} else {
					switch str {
					case "with flow":
						if d, e := gen.readFlow(name, inner, slots); e != nil {
							err = e
						} else {
							out.Flow = append(out.Flow, d)
						}
					case "with group":
						err = gen.readGroup(out, name, inner)

					case "with slot":
						if d, e := gen.readSlot(name, inner); e != nil {
							err = e
						} else {
							out.Slot = append(out.Slot, d)
						}
					case "with str":
						if d, e := gen.readStr(name, inner); e != nil {
							err = e
						} else {
							out.Str = append(out.Str, d)
						}
					case "with num":
						if d, e := gen.readNum(name, inner); e != nil {
							err = e
						} else {
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

func (gen *generator) readGroup(out *groupContent, n string, msg compact.Message) (err error) {
	if msg.Key != "Group contains:" {
		err = fmt.Errorf("expected group definition, have %s", msg.Key)
	} else if msgs, e := parseMessages(msg.Args[0]); e != nil {
		err = e
	} else {
		groups := gen.groups
		gen.groups.AddGroup(n)
		for _, inner := range msgs {
			if e := gen.readSpec(out, inner); e != nil {
				err = e
				break
			}
		}
		gen.groups = groups
	}
	return
}

func (gen *generator) readFlow(n string, msg compact.Message, slots []string) (ret flowData, err error) {
	mm := messageMap(msg)
	if lede, e := mm.GetString("", ""); e != nil {
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
func (gen *generator) readSlot(n string, msg compact.Message) (ret slotData, err error) {
	mm := messageMap(msg)
	if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if cnt := len(uses); cnt > 0 {
		err = errors.New("expected no constraints for slot")
	} else {
		ret = slotData{n}
	}
	return
}

func (gen *generator) readStr(n string, msg compact.Message) (ret strData, err error) {
	mm := messageMap(msg)
	if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if ex, e := mm.GetBool("exclusively"); e != nil {
		err = e
	} else if !ex {
		ret, err = gen.readSimpleStr(n, uses)
	} else {
		ret, err = gen.readEnum(n, uses)
	}
	return
}

func (gen *generator) readSimpleStr(n string, uses []compact.Message) (ret strData, err error) {
	if cnt := len(uses); cnt > 0 {
		err = errors.New("expected no constraints for simple string")
	} else {
		ret = strData{n, nil}
	}
	return
}

func (gen *generator) readEnum(n string, uses []compact.Message) (ret strData, err error) {
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
			ret = strData{n, options}
		}
	}
	return
}

func (gen *generator) readNum(n string, msg compact.Message) (ret numData, err error) {
	mm := messageMap(msg)
	if _, ok := mm["exclusively"]; ok {
		err = errors.New("exclusive nums not supported")
	} else if uses, e := parseMessages(mm["uses"]); e != nil {
		err = e
	} else if len(uses) != 0 {
		err = errors.New("expected no numeric constraints")
	} else {
		ret = numData{n}
	}
	return
}
func (gen *generator) runTemplate(name string, data any) error {
	return gen.tmp.ExecuteTemplate(gen.w, name+".tmpl", data)
}
