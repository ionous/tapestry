package story

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"github.com/ionous/errutil"
)

var AllSlats = append(
	tapestry.AllSlats,
	Slats,
)

var AllSignatures = append(
	tapestry.AllSignatures,
	Signatures,
)

var reg composer.TypeRegistry

func Registry() composer.TypeRegistry {
	if reg == nil {
		for _, slats := range AllSlats {
			reg.RegisterTypes(slats)
		}
	}
	return reg
}

func CompactDecode(msg json.RawMessage) (ret StoryFile, err error) {
	errutil.Panic = true
	err = Decode(&ret, msg, AllSignatures)
	errutil.Panic = false
	return
}

func DetailedDecode(msg json.RawMessage) (ret StoryFile, err error) {
	err = din.Decode(&ret, Registry(), msg)
	return
}
