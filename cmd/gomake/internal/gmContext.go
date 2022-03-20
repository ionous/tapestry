package gomake

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
)

// gobal information needed while running templates
type Context struct {
	currentGroup string
	types        rs.TypeSpecs
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
	typeSpec, ok := ctx.types.Types[typeName]
	return typeSpec, ok
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
	if unboxType, ok := unbox[typeName]; ok && !ignoreUnboxing {
		ret = unboxType
	} else {
		ret = ctx.scopeOf(typeName) + pascal(typeName)
	}
	return
}
