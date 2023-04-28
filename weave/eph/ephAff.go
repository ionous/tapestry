package eph

import "git.sr.ht/~ionous/tapestry/dl/composer"

// Affinity requires a predefined string.
type Affinity struct {
	Str string
}

func (op *Affinity) String() string {
	return op.Str
}

const Affinity_Bool = "$BOOL"
const Affinity_Number = "$NUMBER"
const Affinity_Text = "$TEXT"
const Affinity_Record = "$RECORD"
const Affinity_NumList = "$NUM_LIST"
const Affinity_TextList = "$TEXT_LIST"
const Affinity_RecordList = "$RECORD_LIST"

func (*Affinity) Compose() composer.Spec {
	return composer.Spec{
		Uses: composer.Type_Str,
		Choices: []string{
			Affinity_Bool, Affinity_Number, Affinity_Text, Affinity_Record, Affinity_NumList, Affinity_TextList, Affinity_RecordList,
		},
		Strings: []string{
			"bool", "number", "text", "record", "num_list", "text_list", "record_list",
		},
	}
}
