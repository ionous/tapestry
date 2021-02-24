package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

var Slots = []composer.Slot{{
	Type: (*Comparator)(nil),
	Desc: "Comparison Types: Helper for comparing values.",
}, {
	Type: (*Brancher)(nil),
	Desc: "Helper for choose action.",
}, {
	Type: (*FromSourceFields)(nil),
	Desc: "Helper for getting fields.",
}, {
	Type: (*IntoTargetFields)(nil),
	Desc: "Helper for setting fields.",
}}

var Slats = []composer.Composer{
	(*Activity)(nil),
	(*Argument)(nil),
	(*Arguments)(nil),
	// some boolean tests:
	(*Always)(nil),
	(*AllTrue)(nil),
	(*AnyTrue)(nil),
	(*CompareNum)(nil),
	(*CompareText)(nil),
	(*HasDominion)(nil),
	(*IsNotTrue)(nil), // inverts a bool eval

	// Assign turns an rt.Assignment a normal statement.
	(*Assign)(nil),
	(*Variable)(nil),
	(*FromBool)(nil),
	(*FromNum)(nil),
	(*FromText)(nil),
	(*FromRecord)(nil),
	(*FromNumbers)(nil),
	(*FromTexts)(nil),
	(*FromRecords)(nil),

	// literals
	(*Bool)(nil),
	(*Make)(nil),
	(*Number)(nil),
	(*Numbers)(nil),
	(*Lines)(nil),
	(*Text)(nil),
	(*Texts)(nil),
	// return a value
	(*Var)(nil),
	(*ChooseNum)(nil), // FIX: Choose scalar/any?
	(*ChooseText)(nil),

	// math
	(*SumOf)(nil),
	(*DiffOf)(nil),
	(*ProductOf)(nil),
	(*QuotientOf)(nil),
	(*RemainderOf)(nil),

	//object
	(*ObjectExists)(nil),
	(*NameOf)(nil),
	(*KindOf)(nil),
	(*IsKindOf)(nil),
	(*IsExactKindOf)(nil),
	(*KindsOf)(nil),
	(*HasTrait)(nil),

	(*PrintNum)(nil),
	(*PrintNumWord)(nil),

	// FIX: should take a speaker, and we should have a default speaker
	(*Say)(nil),
	(*Buffer)(nil),
	(*Span)(nil),
	(*Bracket)(nil),
	(*Slash)(nil),
	(*Commas)(nil),
	// text transforms
	(*MakeSingular)(nil),
	(*MakePlural)(nil),
	(*MakeUppercase)(nil),
	(*MakeLowercase)(nil),
	(*MakeTitleCase)(nil),
	(*MakeSentenceCase)(nil),
	(*MakeReversed)(nil),
	//
	// loops and sequences
	(*While)(nil),
	(*Next)(nil),
	(*Break)(nil),
	(*CycleText)(nil),
	(*ShuffleText)(nil),
	(*StoppingText)(nil),

	// text
	(*IsEmpty)(nil),
	(*Includes)(nil),
	(*Join)(nil),
	(*Matches)(nil),

	// comparison
	(*EqualTo)(nil),
	(*NotEqualTo)(nil),
	(*GreaterThan)(nil),
	(*LessThan)(nil),
	(*GreaterOrEqual)(nil),
	(*LessOrEqual)(nil),
	// get at field
	(*GetAtField)(nil),
	(*FromObj)(nil),
	(*FromRec)(nil),
	(*FromVar)(nil),
	// put at field
	(*PutAtField)(nil),
	(*IntoObj)(nil),
	(*IntoVar)(nil),
	// choose action (if)
	(*ChooseAction)(nil),
	(*ChooseMore)(nil),
	(*ChooseNothingElse)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return cmdErrorCtx(op, "", err)
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	// avoid triggering errutil panics for break statements
	if _, ok := err.(DoInterrupt); !ok {
		e := &composer.CommandError{Cmd: op, Ctx: ctx}
		err = errutil.Append(err, e)
	}
	return err
}
