package cmdweave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// ensure there's a scene to contain the passed statements
// name will be the scene name if no existing scene was found.
func wrapScene(scene *story.DefineScene, name, path string, src []story.StoryStatement) []story.StoryStatement {
	note, content := splitHeader(src)
	// note: we're implicitly dependent on the scene index ( if any. )
	if fileScene, ok := hasScene(content); ok {
		// put the rest of the file content into it.
		// ( copy to avoid aliasing slice memory. )
		fileScene.Statements = append([]story.StoryStatement{}, content[1:]...)
		// keep the starting comments and fileScene
		src = src[:len(note)+1]
	} else if scene == nil {
		// without a scene index, we need a scene for the file itself.
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
		src = []story.StoryStatement{scene}
	}
	return src
}

func hasScene(content []story.StoryStatement) (ret *story.DefineScene, okay bool) {
	if len(content) > 0 {
		ret, okay = content[0].(*story.DefineScene)
	}
	return
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
