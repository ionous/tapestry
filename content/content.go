// Package content - this go package exists so that tools can embed the standard lib.
package content

import "embed"

//go:embed shared/*.tell
var Shared embed.FS

//go:embed stories/*.tell
var Sample embed.FS

// a template for a default story
var DefaultStory = `
Tapestry:
- # This story was created using 'tap new' by {{ .Author }}.
  #
  Define scene:requires:with:
  - {{ printf "%q" .Story }}
  - - "Tapestry"
  - - Declare: """
      The title of the story is {{ printf "%q," .Title }}
      The author of the story is {{ printf "%q." .Author }}
      The lobby is a room. You are in the lobby.
      The description of the lobby is "An empty space, waiting to filled with life."
      The chest is a container. The chest is in lobby.
      Some coins are in the chest.
      """ # "
`

type DefaultDesc struct {
	Story, Title, Author string
}
