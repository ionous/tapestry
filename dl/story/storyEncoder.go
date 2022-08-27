package story

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
)

// Write a story to a story file.
func Encode(src *StoryFile) (interface{}, error) {
	return cout.Encode(src, CompactEncoder)
}

// customized writer of compact data
func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	switch op := flow.GetFlow().(type) {
	case *core.CallPattern:
		// rewrite pattern calls to look like normal operations.
		patName := recase(op.Pattern.Str, true)
		if err = m.MarshalBlock(fakeBlock(patName)); err == nil {
			for _, arg := range op.Arguments {
				argName := recase(arg.Name, false)
				if e := m.MarshalKey(argName, argName); e != nil {
					err = e
					break
				} else if e := rt.Assignment_Marshal(m, &arg.From); e != nil {
					err = e
					break
				}
			}
			m.EndBlock()
		}

	default:
		err = core.CompactEncoder(m, flow)
	}
	return
}

type fakeBlock string

func (fakeBlock) GetType() string          { return "fakeBlock" }
func (fb fakeBlock) GetLede() string       { return string(fb) }
func (fakeBlock) GetFlow() interface{}     { return nil }
func (fakeBlock) SetFlow(interface{}) bool { return false }

// pascal when true, camel when false
func recase(str string, cap bool) string {
	u := lang.Underscore(str)
	rs := []rune(u)
	var i int
	for j, cnt := 0, len(rs); j < cnt; j++ {
		if n := rs[j]; n == '_' {
			cap = true
		} else {
			if !cap {
				n = unicode.ToLower(n)
			} else {
				n = unicode.ToUpper(n)
				cap = false
			}
			rs[i] = n
			i++
		}
	}
	return string(rs[:i])
}

// change from old format composer friendly paragraph blocks into simpler to read and edit lines.
func (op *Story) Reformat() (out []StoryStatement) {
	for i, p := range op.Paragraph {
		// every new paragraph, write a "story break"
		if i > 0 || len(p.Markup) > 0 {
			out = append(out, &StoryBreak{p.Markup})
		}
		// add all the lines of the paragraph to the output.
		for _, s := range p.StoryStatement {
			out = append(out, s)
		}
	}
	return
}

// change from simpler to read story lines into old format composer friendly blocks of paragraphs.
func ReformatStory(lines []StoryStatement) (out Story) {
	var p Paragraph
	for i, el := range lines {
		if br, ok := el.(*StoryBreak); !ok {
			// not a story break, add the statement to the current paragraph.
			p.StoryStatement = append(p.StoryStatement, el)
		} else if i == 0 {
			// if the first statement was a story break,
			// that was just a helper to store the first paragraph's comment.
			p.Markup = br.Markup
		} else {
			// any (other) story breaks generate new paragraphs
			// ( so first, flush our old one )
			out.Paragraph = append(out.Paragraph, p)
			// the comment from the break is the comment of the new paragraph
			p = Paragraph{Markup: br.Markup}
		}
	}
	// flush any pending paragraph
	// ( and technically stories always have at least one paragraph anyway )
	out.Paragraph = append(out.Paragraph, p)
	return
}

// story break is an empty do nothing statement, used as a paragraph marker.
func (op *StoryBreak) ImportPhrase(k *Importer) error { return nil }
