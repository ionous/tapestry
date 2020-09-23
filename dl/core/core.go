package core

import (
	"github.com/ionous/iffy/dl/composer"
)

var Slots = []composer.Slot{{
	Name: "comparator",
	Type: (*Comparator)(nil),
	Desc: "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc.",
}, {
	Name: "assignment",
	Type: (*Assignment)(nil),
	Desc: "Assignments: Helper used when setting variables.",
}, {
	Name: "object_ref",
	Type: (*ObjectRef)(nil),
	Desc: "Object Reference: Helper used when referring to objects.",
}}

var Slats = []composer.Slat{
	(*Activity)(nil),
	(*Always)(nil),
	(*AllTrue)(nil),
	(*AnyTrue)(nil),

	// Assign turns an Assignment a normal statement.
	(*Assign)(nil),
	(*FromVar)(nil),
	(*FromBool)(nil),
	(*FromNum)(nil),
	(*FromText)(nil),
	(*FromNumList)(nil),
	(*FromTextList)(nil),

	// FIX: Choose scalar/any?
	(*Choose)(nil),
	(*ChooseNum)(nil),
	(*ChooseText)(nil),
	// FIX: compare scalar?
	(*CompareNum)(nil),
	(*CompareText)(nil),

	(*DoNothing)(nil),
	(*ForEachNum)(nil),
	(*ForEachText)(nil),

	(*GetField)(nil),
	(*GetVar)(nil),
	(*HasTrait)(nil),

	(*IsTrue)(nil),    // transparent pass-through of a bool eval
	(*IsNotTrue)(nil), // inverts a bool eval

	(*Bool)(nil),
	(*Number)(nil),
	(*Text)(nil),
	(*Lines)(nil),
	(*Numbers)(nil),
	(*Texts)(nil),

	(*SumOf)(nil),
	(*DiffOf)(nil),
	(*ProductOf)(nil),
	(*QuotientOf)(nil),
	(*RemainderOf)(nil),

	(*SimpleNoun)(nil),
	(*ObjectName)(nil),
	(*KindOf)(nil),
	(*IsKindOf)(nil),
	(*IsExactKindOf)(nil),

	(*PrintNum)(nil),
	(*PrintNumWord)(nil),

	// FIX: take "List" generically
	(*LenOfNumbers)(nil),
	(*LenOfTexts)(nil),
	(*Range)(nil),

	// FIX: should take a speaker, and we should have a default speaker
	(*Say)(nil),
	(*Buffer)(nil),
	(*Span)(nil),
	(*Bracket)(nil),
	(*Slash)(nil),
	(*Commas)(nil),

	(*MakeSingular)(nil),
	(*MakePlural)(nil),

	(*CycleText)(nil),
	(*ShuffleText)(nil),
	(*StoppingText)(nil),

	// FIX: set any; with the type picked during assembly
	(*SetFieldBool)(nil),
	(*SetFieldNum)(nil),
	(*SetFieldText)(nil),
	(*SetFieldNumList)(nil),
	(*SetFieldTextList)(nil),

	(*IsEmpty)(nil),
	(*Includes)(nil),
	(*Join)(nil),

	// comparison
	(*EqualTo)(nil),
	(*NotEqualTo)(nil),
	(*GreaterThan)(nil),
	(*LessThan)(nil),
	(*GreaterOrEqual)(nil),
	(*LessOrEqual)(nil),
}
