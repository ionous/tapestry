package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

// Write a story to a story file.
func Encode(src *StoryFile) (interface{}, error) {
	return cout.Encode(src, CompactEncoder)
}

// customized writer of compact data
var CompactEncoder = core.CompactEncoder

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
func (op *StoryBreak) PostImport(k *imp.Importer) error { return nil }
