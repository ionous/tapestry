package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

var Slots = []composer.Slot{{
	Type: (*Comparator)(nil),
	Desc: "Comparison types: Helper for comparing values.",
}, {
	Type: (*Brancher)(nil),
	Desc: "Helper for choose action.",
}, {
	Type: (*FromSourceFields)(nil),
	Desc: "Helper for getting fields.",
}, {
	Type: (*IntoTargetFields)(nil),
	Desc: "Helper for setting fields.",
}, {
	Type: (*Trigger)(nil),
	Desc: "Trigger types: Helper for counting values.",
}}

var Slats = []composer.Composer{
	(*Activity)(nil),
	(*AllTrue)(nil),
	(*Always)(nil),
	(*AnyTrue)(nil),
	(*Argument)(nil),
	(*Arguments)(nil),
	(*Assign)(nil),
	(*Bool)(nil),
	(*Bracket)(nil),
	(*Break)(nil),
	(*Buffer)(nil),
	(*Capitalize)(nil),
	(*ChooseAction)(nil),
	(*ChooseMore)(nil),
	(*ChooseNothingElse)(nil),
	(*ChooseNum)(nil), // FIX: Choose scalar/any?
	(*ChooseText)(nil),
	(*ChooseValue)(nil),
	(*Commas)(nil),
	(*CompareNum)(nil),
	(*CompareText)(nil),
	(*CountOf)(nil),
	(*CycleText)(nil),
	(*Determine)(nil),
	(*DiffOf)(nil),
	(*During)(nil),
	(*EqualTo)(nil),
	(*FromBool)(nil),
	(*FromNum)(nil),
	(*FromNumbers)(nil),
	(*FromObj)(nil),
	(*FromRec)(nil),
	(*FromRecord)(nil),
	(*FromRecords)(nil),
	(*FromText)(nil),
	(*FromTexts)(nil),
	(*FromVar)(nil),
	(*GetAtField)(nil),
	(*GreaterOrEqual)(nil),
	(*GreaterThan)(nil),
	(*HasDominion)(nil),
	(*HasTrait)(nil),
	(*IdOf)(nil),
	(*Includes)(nil),
	(*IntoObj)(nil),
	(*IntoVar)(nil),
	(*IsEmpty)(nil),
	(*IsExactKindOf)(nil),
	(*IsKindOf)(nil),
	(*IsNotTrue)(nil), // inverts a bool eval
	(*Join)(nil),
	(*KindOf)(nil),
	(*KindsOf)(nil),
	(*LessOrEqual)(nil),
	(*LessThan)(nil),
	(*Lines)(nil),
	(*Make)(nil),
	(*MakeLowercase)(nil),
	(*MakePlural)(nil),
	(*MakeReversed)(nil),
	(*MakeSentenceCase)(nil),
	(*MakeSingular)(nil),
	(*MakeTitleCase)(nil),
	(*MakeUppercase)(nil),
	(*Matches)(nil),
	(*NameOf)(nil),
	(*Newline)(nil),
	(*Next)(nil),
	(*NotEqualTo)(nil),
	(*Number)(nil),
	(*Numbers)(nil),
	(*ObjectExists)(nil),
	(*Paragraph)(nil),
	(*PrintNum)(nil),
	(*PrintNumWord)(nil),
	(*ProductOf)(nil),
	(*PutAtField)(nil),
	(*QuotientOf)(nil),
	(*RemainderOf)(nil),
	(*Response)(nil),
	(*Row)(nil),
	(*Rows)(nil),
	(*Say)(nil),
	(*Send)(nil),
	(*SetTrait)(nil),
	(*ShuffleText)(nil),
	(*Slash)(nil),
	(*Softline)(nil),
	(*Span)(nil),
	(*StoppingText)(nil),
	(*SumOf)(nil),
	(*Text)(nil),
	(*Texts)(nil),
	(*TriggerCycle)(nil),
	(*TriggerOnce)(nil),
	(*TriggerSwitch)(nil),
	(*Var)(nil),
	(*Variable)(nil),
	(*While)(nil),
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
