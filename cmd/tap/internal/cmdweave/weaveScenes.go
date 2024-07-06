package cmdweave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// name will be the scene name if no existing scene was found.
func wrapScene(name, path string, els []story.StoryStatement) []story.StoryStatement {
	note, content := splitHeader(els)
	if hasEmptyScene(content) {
		scene := content[0].(*story.DefineScene)
		// put the file contents into it.
		// copying to avoid aliasing slice memory
		scene.Statements = append([]story.StoryStatement{}, content[1:]...)
		// keep the original statements, ending with that scene
		els = els[:len(note)+1]
	} else if len(els) > 0 {
		// create a new root scene, and put the content in it
		// ( copy to handle overlapping slices )
		copy := append([]story.StoryStatement{}, content...)
		scene := &story.DefineScene{
			SceneName:  &literal.TextValue{Value: name},
			Statements: copy,
			Markup: map[string]any{
				compact.File:    path,
				compact.Comment: note,
			},
		}
		els = []story.StoryStatement{scene}
	}
	return els
}

func hasEmptyScene(content []story.StoryStatement) (okay bool) {
	if len(content) > 0 {
		if scene, ok := content[0].(*story.DefineScene); ok && len(scene.Statements) == 0 {
			okay = true
		}
	}
	return
}

// extract all leading comments, and return any other statements following those
func splitHeader(els []story.StoryStatement) (ret []string, rest []story.StoryStatement) {
	for i, el := range els {
		if note, ok := el.(*story.StoryNote); ok {
			ret = append(ret, note.Text)
		} else {
			rest = els[i:]
			break
		}
	}
	return
}
