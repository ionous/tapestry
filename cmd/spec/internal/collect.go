package internal

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	r "reflect"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/export"
)

type Collect struct {
	all    []export.Dict
	slots  []r.Type
	groups Groups
}

func (c *Collect) AddGroup(out export.Dict, group string) {
	if c.groups == nil {
		c.groups = make(Groups)
	}
	c.groups.addGroup(out, group)
}

func (c *Collect) AddSlot(slot composer.Slot) {
	spec := getSpec(slot.Type)
	if spec.Group != "internal" {
		i := r.TypeOf(slot.Type).Elem()
		//
		if len(spec.Name) == 0 {
			spec.Name = slot.Name
		}
		if len(spec.Desc) == 0 {
			spec.Desc = export.Prettify(slot.Name)
		}
		if len(spec.Group) == 0 {
			spec.Group = slot.Group
		}
		out := export.Dict{
			"name": spec.Name,
			"desc": spec.Desc,
			"uses": "slot",
		}
		addDesc(out, slot.Desc)
		c.AddGroup(out, spec.Group)
		c.all = append(c.all, out)
		c.slots = append(c.slots, i)
	}
}

func (c *Collect) AddSlat(cmd composer.Slat) {
	if spec := cmd.Compose(); spec.Group != "internal" {
		rtype := r.TypeOf(cmd).Elem()
		if len(spec.Name) == 0 {
			panic(fmt.Sprintln("missing name for type", rtype.Name()))
		}
		//
		with := make(export.Dict)
		if slotNames := slotsOf(rtype, c.slots); len(slotNames) > 0 {
			with["slots"] = slotNames
		}
		out := export.Dict{
			"name": spec.Name,
			"uses": "run",
			"with": with,
		}
		// missing spec, missing slots.
		if len(spec.Spec) != 0 {
			out["spec"] = spec.Spec
		} else {
			tokens, params := parse(rtype)
			with["params"] = params
			with["tokens"] = updateTokens(spec.Spec, tokens)
		}
		addDesc(out, spec.Desc)
		c.AddGroup(out, spec.Group)
		c.all = append(c.all, out)
	}
}

func (c *Collect) FlushGroups() {
	c.all = c.groups.appendGroups(c.all)
}

func (c *Collect) Sort() {
	sort.Slice(c.all, func(idx, jdx int) (ret bool) {
		i, j := c.all[idx], c.all[jdx]
		uses := strings.Compare(i["uses"].(string), j["uses"].(string))
		switch uses {
		case 0:
			ret = i["name"].(string) < j["name"].(string)
		case -1:
			ret = false
		case 1:
			ret = true
		}
		return

	})
}

func (c *Collect) Marshal() (ret []byte, err error) {
	return json.MarshalIndent(c.all, "", "  ")
}