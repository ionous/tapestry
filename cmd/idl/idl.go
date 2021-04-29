package main

import (
	"hash/fnv"
	"io"
	"os"
	"path"
	r "reflect"
	"sort"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/cmd/idl/internal"
	"git.sr.ht/~ionous/iffy/dl/composer"
)

var temp = internal.Cap

type Pack struct {
	Name  string
	Hash  string
	Deps  internal.Deps
	Slots []internal.SlotMessage
	Slats []internal.SlatMessage
}

func (p *Pack) Sort() {
	sort.Slice(p.Slots, func(i, j int) bool {
		return p.Slots[i].Name < p.Slots[j].Name
	})
	sort.Slice(p.Slats, func(i, j int) bool {
		return p.Slats[i].Name < p.Slats[j].Name
	})
}

type All struct {
	Packages map[string]Pack
}

// capnp compile -I$GOPATH/src/zombiezen.com/go/capnproto2/std -ogo ./allCmds.capnp
// capnp compile -I$GOPATH/src/zombiezen.com/go/capnproto2/std -ogo ./allCmds.capnp ./core/core.capnp ./debug/debug.capnp  ./grammar/grammar.capnp  ./list/list.capnp  ./rel/rel.capnp  ./render/render.capnp
func main() {
	all := All{make(map[string]Pack)}
	all.makeSlots()
	all.makeSlats()
	base := os.ExpandEnv("$GOPATH/src/git.sr.ht/~ionous/iffy/idl")
	for k, p := range all.Packages {
		// ex. src/git.sr.ht/~ionous/iffy/idl/core
		dir := path.Join(base, k)
		if e := os.MkdirAll(dir, os.ModePerm); e != nil {
			panic(e)
		} else {
			fn := path.Join(dir, k+temp.Ext)
			if fp, e := os.Create(fn); e != nil {
				panic(e)
			} else {
				p.Sort()
				p.Name = k
				hash := fnv.New64a()
				io.WriteString(hash, fn)
				p.Hash = strconv.FormatUint((1<<63)|hash.Sum64(), 16)
				temp.Pack.Must(fp, p)
			}
		}
	}
	{
		fn := path.Join(base, "allCmds"+temp.Ext)
		if fp, e := os.Create(fn); e != nil {
			panic(e)
		} else {
			var slots []internal.SlotMessage
			var deps internal.Deps
			for k, p := range all.Packages {
				deps = deps.AddDep(k)
				slots = append(slots, p.Slots...)
			}
			sort.Strings(deps)
			sort.Slice(slots, func(i, j int) bool {
				return slots[i].Name < slots[j].Name
			})
			temp.All.Must(fp, map[string]interface{}{
				"Slots": slots,
				"Deps":  deps,
			})
		}
	}
}

var allCmds internal.Cmds = internal.MakeCommands()

func (all *All) makeSlots() {
	// write slots
	for _, slots := range iffy.AllSlots {
		for _, slot := range slots {
			var sigs []internal.Sig
			for _, slat := range internal.ImplementorsOf(slot.Type) {
				cmd := allCmds.Add(internal.MakeCommand(slat))
				if cmd != nil {
					sigs = append(sigs, cmd.Signatures()...)
				}
			}
			sort.Slice(sigs, func(i, j int) bool {
				return sigs[i].Raw < sigs[j].Raw
			})
			name := internal.Pascal(composer.SlotName(slot))
			msg := internal.SlotMessage{
				Name: name,
				Desc: internal.ClipDesc(slot.Desc),
				Sigs: sigs,
			}
			//
			pack := shortPack(r.TypeOf(slot.Type).Elem())
			p := all.Packages[pack]
			p.Slots = append(p.Slots, msg)
			all.Packages[pack] = p
		}
	}
}

func (all *All) makeSlats() {
	for _, cmd := range allCmds {
		//
		pack := shortPack(cmd.Type)
		p := all.Packages[pack]
		msg := internal.SlatMessage{cmd}
		p.Slats = append(p.Slats, msg)
		// accumulate dependencies
		for _, place := range cmd.Places {
			inner := shortPack(place.Type)
			if inner != pack {
				p.Deps = p.Deps.AddDep(inner)
			}
		}
		all.Packages[pack] = p
	}
}

func shortPack(rtype r.Type) string {
	ns := strings.Split(rtype.PkgPath(), "/")
	return ns[len(ns)-1]
}
