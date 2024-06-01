package cmdweave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
)

// name will be the scene name if no existing scene was found.
func wrapScene(name string, els []story.StoryStatement) []story.StoryStatement {
	// lhs will have the leading comments, body everything else
	if lhs, body := splitStatements(els); len(body) > 0 {
		// no root scene, or the root scene isn't empty?
		// create a file level scene, and put everything inside it.
		if scene, ok := body[0].(*story.DefineScene); !ok || len(scene.Statements) != 0 {
			// create a root scene, the body in it.
			// ( lhs and body share the same backing, so we have to copy one or the other )
			els = append(lhs, &story.DefineScene{
				SceneName:  &literal.TextValue{Value: name},
				Statements: append([]story.StoryStatement{}, body...),
			})
		} else {
			// otherwise, use the empty scene as the file scene
			scene.Statements = body[1:] // everything except the scene
			els = els[:len(lhs)+1]      // everything up to and including the scene
		}
	}
	return els
}

// split so that all leading comments are on the lhs;
// the first statement and everything after are on the rhs
func splitStatements(els []story.StoryStatement) (lhs, rhs []story.StoryStatement) {
	lhs = els
	for i, el := range els {
		if _, ok := el.(*story.StoryNote); !ok {
			lhs, rhs = els[:i], els[i:]
			break
		}
	}
	return
}
