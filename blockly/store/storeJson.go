package store

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/web/js"
)

type BlockFile struct {
	TopBlocks `json:"blocks"`
}

type TopBlocks struct {
	LanguageVersion float64     `json:"languageVersion"`
	Blocks          []BlockInfo `json:"blocks"`
}

type BlockInfo struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	ExtraState map[string]int         `json:"extraState"`
	Inputs     map[string]*BlockInput `json:"inputs"`
	Fields     js.MapSlice            `json:"fields"`
	Next       *BlockInput            `json:"next"`
}

type BlockInput struct {
	*BlockInfo `json:"block"`
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
func (bi *BlockInput) CountNext() (ret int) {
	for ab := bi; ab.Next != nil; ab = ab.Next {
		ret++
	}
	return
}

// return the number of term# formatted fields
func (bi *BlockInfo) CountField(term string) (retStart, retCnt int) {
	next := term + strconv.Itoa(retCnt)
	if at := bi.Fields.FindIndex(next); at >= 0 {
		retStart, retCnt = at, 1
		for _, f := range bi.Fields[at+1:] {
			if is := f.Key == term+strconv.Itoa(retCnt); !is {
				break
			}
			retCnt++
		}
	}
	return
}
