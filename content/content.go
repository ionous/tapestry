// Package content - this go package exists so that tools can embed the standard lib.
package content

import "embed"

//go:embed shared/*.tell
var Shared embed.FS

//go:embed stories/*.tell
var Sample embed.FS

// a template for a default story
var DefaultStory = //
`# --------------------------------------------------------
# This story was created by {{ .Author }} using 'tap new'.
# --------------------------------------------------------
The title of the story is {{ printf "%q." .Title }}
The author of the story is {{ printf "%q." .Author }}
The Empty Space is a room. You are in the space.
The description of the space is "An empty space, waiting to be filled with life."
`

type DefaultDesc struct {
	Story, Title, Author string
}
