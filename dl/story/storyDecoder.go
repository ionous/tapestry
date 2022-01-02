package story

import (
	"encoding/json"
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"github.com/ionous/errutil"
)

// Story decode assumes a story file format.
// It can read either a normal compact encoded value, or the special story encoded data.
// ( see Encode for details. )
func Decode(dst jsn.Marshalee, msg json.RawMessage, sig cin.Signatures) (err error) {
	for i, b := range msg {
		if !unicode.IsSpace(rune(b)) {
			msg = msg[i:]
			break
		}
	}
	if len(msg) > 0 {
		dec := cin.NewDecoder(sig).
			SetFlowDecoder(CompactFlowDecoder).
			SetSlotDecoder(CompactSlotDecoder)

		switch msg[0] {
		default:
			err = dec.Decode(dst, msg)
		case 'n':
			// null?
			if out, ok := dst.(*Story); !ok {
				err = dec.Decode(dst, msg)
			} else {
				var empty string
				if e := json.Unmarshal(msg, &empty); e != nil {
					err = e
				} else if len(empty) != 0 {
					err = errutil.New("unexpected or invalid json")
				} else {
					out.Paragraph = make([]Paragraph, 0)
				}
			}
		case '[':
			// a quick if hacky way of reading compact paragraphs.
			// ( otherwise we expect to see: {"Story:": [{"Paragraph:": [...cmds...]}, ...more paragraphs...]}
			// instead of the nicer: "[[...cmds...], ...more paragraphs...]"
			if out, ok := dst.(*Story); !ok {
				err = errutil.New("expected story data")
			} else {
				var ps [][]json.RawMessage
				if e := json.Unmarshal(msg, &ps); e != nil {
					err = e
				} else {
				Loop:
					for _, pmsg := range ps {
						var p Paragraph
						for _, msg := range pmsg {
							var s StoryStatement
							if e := dec.Decode(StoryStatement_Slot{&s}, msg); e != nil {
								err = e
								break Loop
							} else {
								p.StoryStatement = append(p.StoryStatement, s)
							}
						}
						out.Paragraph = append(out.Paragraph, p)
					}
				}
			}
		}
	}
	return
}

var CompactSlotDecoder = core.CompactSlotDecoder

// customized reader of compact data
func CompactFlowDecoder(flow jsn.FlowBlock, msg json.RawMessage) (err error) {
	switch typeName := flow.GetType(); typeName {
	default:
		err = chart.Unhandled("CustomFlow")

	case NamedNoun_Type:
		var str string
		if e := json.Unmarshal(msg, &str); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			var out NamedNoun
			str := strings.TrimSpace(str)
			if space := strings.Index(str, " "); space < 0 {
				out.Name.Str = str
			} else {
				jsn.MakeEnum(&out.Determiner, &out.Determiner.Str).SetValue(str[:space])
				out.Name.Str = str[space+1:]
			}
			if !flow.SetFlow(&out) {
				err = errutil.New("could set result to flow", typeName, flow)
			}
		}
	}
	return
}
