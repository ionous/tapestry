package unblock

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/web/js"
)

// describes a blockly file
// https://developers.google.com/blockly/guides/configure/web/serialization?hl=en
type File struct {
	TopBlocks `json:"blocks"`
}

// file formatting info inside every blockly json file
type TopBlocks struct {
	LanguageVersion float64     `json:"languageVersion"`
	Blocks          []BlockInfo `json:"blocks"`
}

// a stored blockly block
// recursive in the sense that inputs can contain other blocks
// note while in shape definitions fields are stored inside of inputs;
// in block data fields get stored separate from inputs, keyed by field name.
type BlockInfo struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	ExtraState map[string]int `json:"extraState"`
	Icons      Icons          `json:"icons"`
	Inputs     js.MapSlice    `json:"inputs"`
	Fields     js.MapSlice    `json:"fields"`
	Next       *Input         `json:"next"`
}

// a blockly input contains a blockly block.
type Input struct {
	*BlockInfo `json:"block"`
}

// blockly stores comments inside of the icon representing the comment
type Icons struct {
	Comment *Comment `json:"comment"`
}

// this is how blockly handles comments
// fix? move to make part of the mutation ( as extra data )
type Comment struct {
	Text   string `json:"text"`
	Pinned bool   `json:"pinned"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func (tb *TopBlocks) FindFirst(typeName string) (ret *BlockInfo, okay bool) {
	for _, b := range tb.Blocks {
		if b.Type == typeName {
			ret, okay = &b, true
			break
		}
	}
	return
}

// return the depth of the linked list starting in next.
func (bi *Input) CountNext() (ret int) {
	for ab := bi; ab.Next != nil; ab = ab.Next {
		ret++
	}
	return
}

// return the number of term# formatted fields
func (bi *BlockInfo) ReadInput(idx int) (ret Input, err error) {
	return readInput(bi.Inputs[idx])
}

// return the number of term# formatted fields
func (bi *BlockInfo) SliceFields(term string) js.MapSlice {
	return sliceTerm(term, bi.Fields)
}

// return the number of term# formatted inputs
func (bi *BlockInfo) SliceInputs(term string) js.MapSlice {
	return sliceTerm(term, bi.Inputs)
}

// "inputs": { "CONTAINS0": {"block":{...}}, "CONTAINS1": {"block":{...}}, ... }
// fix: nothing tests this.
func sliceTerm(term string, msgs js.MapSlice) (ret js.MapSlice) {
	var retStart, retCnt int
	next := term + strconv.Itoa(retCnt)
	if at := msgs.FindIndex(next); at >= 0 {
		retStart, retCnt = at, 1
		for _, f := range msgs[at+1:] {
			if is := f.Key == term+strconv.Itoa(retCnt); !is {
				break
			}
			retCnt++
		}
	}
	return msgs[retStart:retCnt]
}
