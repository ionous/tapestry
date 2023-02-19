package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func Refs(refs ...eph.Ephemera) *eph.EphRefs {
	return &eph.EphRefs{Refs: refs}
}

func ImportCall(op *core.CallPattern) *eph.EphRefs {
	return refArgs(op.Pattern.String(), kindsOf.Pattern, op.Arguments)
}

func refArgs(k string, parentKind kindsOf.Kinds, args []core.Arg) (ret *eph.EphRefs) {
	var refs []eph.EphParams
	for _, arg := range args {
		refs = append(refs, eph.EphParams{
			Name:     arg.Name,
			Affinity: affineToAffinity(arg.Value.Affinity()),
		})
	}
	ret = Refs(&eph.EphKinds{
		Kinds: k,
		// we dont actually know.
		// its probably a pattern, but it could be a record just as well....
		// From:    parentKind.String(),
		Contain: refs,
	})
	return
}

// translate "bool" to "$BOOL", etc.
// note: can return affine.None ( unknown affinity )
func affineToAffinity(a affine.Affinity) (ret eph.Affinity) {
	spec := ret.Compose()
	if k, i := spec.IndexOfValue(a.String()); i >= 0 {
		ret.Str = k
	}
	return
}
