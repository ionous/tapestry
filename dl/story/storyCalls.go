package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func importPattern(op *core.CallPattern) *eph.EphRefs {
	return refArgs(op.Pattern.String(), kindsOf.Pattern, op.Arguments.Args)
}

func (op *Make) ImportStub(k *Importer) (interface{}, error) {
	refs, args := op.Arguments.xform(op.Name, kindsOf.Record)
	k.WriteEphemera(refs)
	return &core.CallMake{Kind: op.Name, Arguments: args}, nil
}

func (op *Send) ImportStub(k *Importer) (interface{}, error) {
	// note: this used to pass "kindOf.Event" but we dont need to be so strict.
	refs, args := op.Arguments.xform(op.Event, kindsOf.Pattern)
	k.WriteEphemera(refs)
	return &core.CallSend{Event: op.Event, Path: op.Path, Arguments: args}, nil
}

func (stubs *Arguments) xform(k string, t kindsOf.Kinds) (retRefs *eph.EphRefs, retCall core.CallArgs) {
	var args []core.CallArg
	var refs []eph.EphParams
	if stubs != nil {
		for _, arg := range stubs.Args {
			args = append(args, core.CallArg{
				Name: arg.Name, // string
				From: arg.From, // assignment
			})
			//
			refs = append(refs, eph.EphParams{
				Name:     arg.Name,
				Affinity: infinityToAffinity(arg.From),
			})
		}
	}
	retCall = core.CallArgs{Args: args}
	retRefs = Refs(&eph.EphKinds{
		Kinds:   k,
		From:    t.String(),
		Contain: refs,
	})
	return
}

func refArgs(k string, t kindsOf.Kinds, args []core.CallArg) (ret *eph.EphRefs) {
	var refs []eph.EphParams
	for _, arg := range args {
		args = append(args, core.CallArg{
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
		Kinds:   k,
		From:    t.String(),
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
