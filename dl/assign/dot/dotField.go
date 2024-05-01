package dot

import (
	"strings"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// field name to pick a value from a record.
type Field string

// raw string
func (dot Field) Field() string {
	return string(dot)
}

// print friendly string
func (dot Field) writeTo(b *strings.Builder) {
	b.WriteRune('.')
	b.WriteString(dot.Field())
}

func (dot Field) Peek(c Cursor) (Cursor, error) {
	return c.GetAtField(string(dot))
}

func (dot Field) Poke(c Cursor, newValue g.Value) (err error) {
	return c.SetAtField(string(dot), newValue)
}
