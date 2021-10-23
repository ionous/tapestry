package story

import (
	"encoding/json"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/jsn/din"
)

var reg din.Registry

func StoryRegistry() din.Registry {
	if reg == nil {
		for _, slats := range iffy.AllSlats {
			reg.RegisterTypes(slats)
		}
		reg.RegisterTypes(Slats)
	}
	return reg
}

func DecodeDetailedStory(story *Story, msg json.RawMessage) error {
	return din.Decode(story, StoryRegistry(), msg)
}
