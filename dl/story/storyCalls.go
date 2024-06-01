package story

// func Refs(refs ...eph.Ephemera) *eph.Refs {
// 	return &eph.Refs{Refs: refs}
// }

// func ImportCall(op *call.CallPattern) *eph.Refs {
// 	return refArgs(op.PatternName, kindsOf.Pattern, op.Arguments)
// }

// func refArgs(k string, parentKind kindsOf.Kinds, args []assign.Arg) (ret *eph.Refs) {
// 	var refs []eph.Params
// 	for _, arg := range args {
// 		refs = append(refs, eph.Params{
// 			Name:     arg.Name,
// 			Affinity: affineToAffinity(assign.GetAffinity(arg.Value)),
// 		})
// 	}
// 	ret = Refs(&eph.Kinds{
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
