package unblock

import (
	"encoding/json"
	"strconv"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/web/js"
)

type File struct {
	TopBlocks `json:"blocks"`
}

type TopBlocks struct {
	LanguageVersion float64 `json:"languageVersion"`
	Blocks          []Info  `json:"blocks"`
}

type Info struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	ExtraState map[string]int `json:"extraState"`
	Icons      Icons          `json:"icons"`
	Inputs     js.MapSlice    `json:"inputs"`
	Fields     js.MapSlice    `json:"fields"`
	Next       *Input         `json:"next"`
}

type Input struct {
	*Info `json:"block"`
}

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

func (tb *TopBlocks) FindFirst(typeName string) (ret *Info, okay bool) {
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
func (bi *Info) CountFields(term string) (retStart, retCnt int) {
	return count(term, bi.Fields)
}

// return the number of term# formatted inputs
func (bi *Info) CountInputs(term string) (retStart, retCnt int) {
	return count(term, bi.Inputs)
}

func (bi *Info) ReadInput(i int) (ret Input, err error) {
	if e := json.Unmarshal(bi.Inputs[i].Msg, &ret); e != nil {
		err = e
	} else if ret.Info == nil {
		err = jsn.Missing
	}
	return
}

func count(term string, msgs js.MapSlice) (retStart, retCnt int) {
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
	return
}
