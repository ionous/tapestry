package story

// func Refs(refs ...eph.Ephemera) *eph.EphRefs {
// 	return &eph.EphRefs{Refs: refs}
// }

// func ImportCall(op *assign.CallPattern) *eph.EphRefs {
// 	return refArgs(op.PatternName, kindsOf.Pattern, op.Arguments)
// }

// func refArgs(k string, parentKind kindsOf.Kinds, args []assign.Arg) (ret *eph.EphRefs) {
// 	var refs []eph.EphParams
// 	for _, arg := range args {
// 		refs = append(refs, eph.EphParams{
// 			Name:     arg.Name,
// 			Affinity: affineToAffinity(assign.GetAffinity(arg.Value)),
// 		})
// 	}
// 	ret = Refs(&eph.EphKinds{
// 		Kind: k,
// 		// we dont actually know.
// 		// its probably a pattern, but it could be a record just as well....
// 		// From:    parentKind.String(),
// 		Contain: refs,
// 	})
// 	return
// }

// translate "bool" to "$BOOL", etc.
// note: can return affine.None ( unknown affinity )
// func affineToAffinity(a affine.Affinity) (ret eph.Affinity) {
// 	spec := ret.Compose()
// 	if k, i := spec.IndexOfValue(a.String()); i >= 0 {
// 		ret.Str = k
// 	}
// 	return
// }
