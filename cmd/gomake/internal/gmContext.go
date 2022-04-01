package gomake

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
)

// gobal information needed while running templates
type Context struct {
	currentGroup string
	types        rs.TypeSpecs

	// when generating some kinds of simple types...
	// replace the specified typename with specified primitive.
	// types that map to numbers, etc. are added as unbox automatically.
	unbox map[string]string
}

func (ctx *Context) GroupOf(typeName string) (ret string, okay bool) {
	if types, ok := ctx.types.Types[typeName]; ok {
		if len(types.Groups) > 0 {
			ret = types.Groups[0]
			okay = true
		}
	}
	return
}

// private types dont have a typespec
// ( ex. those declared manually by an implementation )
func (ctx *Context) GetTypeSpec(typeName string) (ret *spec.TypeSpec, okay bool) {
	block, ok := ctx.types.Types[typeName]
	return block, ok
}

func (ctx *Context) TermsOf(block *spec.TypeSpec) []Term {
	flow := block.Spec.Value.(*spec.FlowSpec)
	terms := make([]Term, len(flow.Terms))
	var pubCount int
	for i, t := range flow.Terms {
		pi := -1
		if !t.Private {
			pi = pubCount
			pubCount++
		}
		terms[i] = Term{ctx, flow, t, pi}
	}
	return terms
}

// return the package qualifier for the passed typename
// ( doesnt add the typename to the return value )
func (ctx *Context) scopeOf(typeName string) (ret string) {
	if g, ok := ctx.GroupOf(typeName); ok && g != ctx.currentGroup {
		ret = g + "."
	}
	return
}

// return the fully qualified typename
func (ctx *Context) scopedName(typeName string, ignoreUnboxing bool) (ret string) {
	if unboxType, ok := ctx.unbox[typeName]; ok && !ignoreUnboxing {
		ret = unboxType
	} else {
		ret = ctx.scopeOf(typeName) + pascal(typeName)
	}
	return
}
