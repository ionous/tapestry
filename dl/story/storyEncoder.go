package story

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

// Write a story to a story file.
func Encode(src *Story) (interface{}, error) {
	weave := src.reformat()
	return cout.Encode(&weave, CompactEncoder)
}

// Read a story from a story file.
func Decode(dst jsn.Marshalee, msg json.RawMessage, sig cin.Signatures) error {
	return cin.NewDecoder(sig).
		SetFlowDecoder(CompactFlowDecoder).
		SetSlotDecoder(CompactSlotDecoder).
		Decode(dst, msg)
}

// customized writer of compact data
func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	switch ptr := flow.GetFlow().(type) {
	case *Story:
		lines := ptr.reformat()
		err = lines.Marshal(m)

	default:
		err = core.CompactEncoder(m, flow)
	}
	return
}

var CompactSlotDecoder = core.CompactSlotDecoder

// customized reader of compact data
func CompactFlowDecoder(m jsn.Marshaler, flow jsn.FlowBlock, msg json.RawMessage) (err error) {
	switch typeName := flow.GetType(); typeName {
	default:
		err = chart.Unhandled("CustomFlow")

	case Story_Type:
		var lines StoryLines
		if e := lines.Marshal(m); e != nil {
			err = e
		} else {
			story := lines.reformat()
			flow.SetFlow(&story)
		}
	}
	return
}

// change from old format composer friendly paragraph blocks into simpler to read and edit lines.
func (op *Story) reformat() (out StoryLines) {
	for i, p := range op.Paragraph {
		// every new paragraph, write a "story break"
		if i > 0 || len(p.UserComment) > 0 {
			out.Lines = append(out.Lines, &StoryBreak{p.UserComment})
		}
		// add all the lines of the paragraph to the output.
		for _, s := range p.StoryStatement {
			out.Lines = append(out.Lines, s)
		}
	}
	return
}

// change from simpler to read story lines into old format composer friendly blocks of paragraphs.
func (op *StoryLines) reformat() (out Story) {
	var p Paragraph
	for i, el := range op.Lines {
		if br, ok := el.(*StoryBreak); !ok {
			// not a story break, add the statement to the current paragraph.
			p.StoryStatement = append(p.StoryStatement, el)
		} else if i == 0 {
			// if the first statement was a story break,
			// that was just a helper to store the first paragraph's comment.
			p.UserComment = br.UserComment
		} else {
			// any (other) story breaks generate new paragraphs
			// ( so first, flush our old one )
			out.Paragraph = append(out.Paragraph, p)
			// the comment from the break is the comment of the new paragraph
			p = Paragraph{UserComment: br.UserComment}
		}
	}
	// flush any pending paragraph
	// ( and technically stories always have at least one paragraph anyway )
	out.Paragraph = append(out.Paragraph, p)
	return
}

// story break is an empty do nothing statement, used as a paragraph marker.
func (op *StoryBreak) ImportPhrase(k *Importer) error { return nil }
