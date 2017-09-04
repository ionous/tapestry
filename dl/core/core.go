package core

type Classes struct {
	*NumberCounter
	*TextCounter
	*ObjCounter
}

type Counters struct {
	*CycleCounter
	*ShuffleCounter
	*StoppingCounter
}

type Commands struct {
	*Add
	*Sub
	*Mul
	*Div
	*Mod
	//
	*AllTrue
	*AnyTrue
	*Bool
	*Buffer
	*SetState
	*Choose
	*ChooseNum
	*ChooseObj
	*ChooseText
	*ClassName
	*CompareNum
	*CompareText
	*CycleText
	*DoNothing
	*EqualTo
	// *Error
	*ForEachNum
	*ForEachObj
	*ForEachText
	*Get
	// *GoCall
	*GreaterThan
	// *Inc
	*Includes
	// *Is/State -> use Get
	*IsEmpty
	*IsSameClass
	*IsSimilarClass
	*IsNot
	// *IsNum
	// *IsObj
	// *IsText
	// *IsValid
	*Len
	*LesserThan
	*NotEqualTo
	*Num
	*Numbers
	*Object
	// *Objects
	// *ObjListContains
	// *ObjListIsEmpty
	// *Object
	*PrintBracket
	*PrintList
	*PrintNum
	*PrintNumWord
	*PrintSpan
	*PrintText
	*Range
	*RelatedList
	*RelationEmpty
	*SetBool
	*SetNum
	*SetText
	*SetObj
	*ShuffleText
	*StoppingText
	// *State
	// *StopNow
	*Text
	*Texts
	// *Using
}
