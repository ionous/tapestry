package weave

import (
	"strconv"
)

// Counters is a helper to generate semi-unique names for a group.
type Counters map[string]uint64

func (m *Counters) Next(name string) string {
	c := (*m)[name] + 1
	(*m)[name] = c // COUNTER:#
	return name + "_" + strconv.FormatUint(c, 36)
}

// NewCounter generates a unique string, and uses local markup to try to create a stable one.
// instead consider  "PreImport" could be used to write a key into the markup if one doesnt already exist.
// and a free function could also extract what it needs from any op's markup.
// ( then Schedule wouldn't need Catalog for counters )
func (k *Catalog) NewCounter(name string, markup map[string]any) (ret string) {
	// fix: use a special "id" marker instead?
	if at, ok := markup["comment"].(string); ok && len(at) > 0 {
		ret = at
	} else {
		ret = k.autoCounter.Next(name)
	}
	return
}
