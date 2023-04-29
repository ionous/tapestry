package weave

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// request future processing from the catalog's importer.
type Schedule interface {
	Schedule(*Catalog) error
}

// PreImport happens at the opening of a json block and it can transform the value into something completely new.
type PreImport interface {
	PreImport(*Catalog) (interface{}, error)
}

// StoryStatement - import a single story statement.
// used during weave, and expects that the runtime is the importer's own runtime.
// ( as opposed to the story's playtime. )
func StoryStatement(run rt.Runtime, op Schedule) (err error) {
	if k, ok := run.(*Weaver); !ok {
		err = errutil.Fmt("runtime %T doesn't support story statements", run)
	} else {
		err = op.Schedule(k.d.catalog)
	}
	return
}

func (k *Catalog) SetSource(x string) {
	k.cursor = x
}

// Env - used for comments to determine if they should turn into log statements.
// todo: remove?
func (k *Catalog) Env() *Environ {
	return &k.env
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
