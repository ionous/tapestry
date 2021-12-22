// Package affine - names for all of Tapestry's built-in data types: primitives, structured types, and arrays thereof.
package affine

// Affinity - name of one of Tapestry's built-in data type.
type Affinity string

// a helper for input/output of affinities including logging:
// returns the stored string or "unknown affinity" if the string is empty.
func (a Affinity) String() (ret string) {
	if a := string(a); len(a) > 0 {
		ret = a
	} else {
		ret = "unknown affinity"
	}
	return
}

const (
	Bool       Affinity = "bool"
	Number     Affinity = "number"
	NumList    Affinity = "num_list"
	Text       Affinity = "text"
	TextList   Affinity = "text_list"
	Object     Affinity = "object"
	Record     Affinity = "record"
	RecordList Affinity = "record_list"
)

// true if one of three list types
func IsList(a Affinity) bool {
	elAffinity := Element(a)
	return len(elAffinity) > 0
}

// return the affinity of an affine list, or blank if not a list.
func Element(list Affinity) (ret Affinity) {
	switch a := list; a {
	case TextList:
		ret = Text
	case NumList:
		ret = Number
	case RecordList:
		ret = Record
	}
	return
}

// true if a structured type.
func HasFields(a Affinity) bool {
	return a == Record || a == Object
}
