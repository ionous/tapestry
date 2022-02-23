package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func importCall(op *core.CallPattern) *eph.EphRefs {
	return refArgs(op.Pattern.String(), kindsOf.Pattern, op.Arguments)
}

func refArgs(k string, parentKind kindsOf.Kinds, args []rt.Arg) (ret *eph.EphRefs) {
	var refs []eph.EphParams
	for _, arg := range args {
		args = append(args, rt.Arg{
			Name: arg.Name, // string
			From: arg.From, // assignment
		})
		//
		refs = append(refs, eph.EphParams{
			Name:     arg.Name,
			Affinity: infinityToAffinity(arg.From),
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

func infinityToAffinity(a interface{ Affinity() affine.Affinity }) (ret eph.Affinity) {
	if a != nil {
		ret = affineToAffinity(a.Affinity())
	}
	return
}

// note: can return "" ( unknown affinity )
func affineToAffinity(a affine.Affinity) (ret eph.Affinity) {
	spec := ret.Compose()
	if k, i := spec.IndexOfValue(a.String()); i >= 0 {
		ret.Str = k
	}
	return
}
