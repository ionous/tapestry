// Package content - this go package exists so that tools can embed the standard lib.
package content

import "embed"

//go:embed shared/*.tell
var Shared embed.FS // we want the _index file; so this uses *.tell

//go:embed stories
var Sample embed.FS // not using *.tell excludes files starting with _

// a template for a default story
var DefaultStory = //
`# {{ .Title }}
# This story was created by {{ .Author }} using 'tap new'.

The title of the story is {{ printf "%q." .Title }}
The author of the story is {{ printf "%q." .Author }}
The Empty Space is a room. You are in the space.
The description of the space is "An empty space, waiting to be filled with life."
`

type DefaultDesc struct {
	Story, Title, Author string
}
