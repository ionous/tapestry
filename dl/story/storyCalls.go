package story

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/rt/kindsOf"
)

// transforms "story.Determine" into the "core.CallPattern" command.
// while the two commands are equivalent, this provides a hook to verify
// the caller's arguments and the pattern's parameters match up.
func (op *Determine) ImportStub(k *Importer) (interface{}, error) {
	refs, args := op.Arguments.xform(op.Name.String(), kindsOf.Pattern)
	k.Write(&refs)
	return &core.CallPattern{Pattern: op.Name, Arguments: args}, nil
}

func (op *Make) ImportStub(k *Importer) (interface{}, error) {
	refs, args := op.Arguments.xform(op.Name, kindsOf.Record)
	k.Write(&refs)
	return &core.CallMake{Kind: op.Name, Arguments: args}, nil
}

func (op *Send) ImportStub(k *Importer) (interface{}, error) {
	refs, args := op.Arguments.xform(op.Event, kindsOf.Event)
	k.Write(&refs)
	return &core.CallSend{Event: op.Event, Path: op.Path, Arguments: args}, nil
}

func (stubs *Arguments) xform(k string, t kindsOf.Kinds) (refRef eph.EphRefs, retCall core.CallArgs) {
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
	return eph.EphRefs{Kinds: k, From: t.String(), ReferTo: refs}, core.CallArgs{args}
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
