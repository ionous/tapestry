package internal

import (
	"hash/fnv"
	"io"
)

type SlotMessage struct {
	Name string // name of the message in pascal like text
	Desc string
	Sigs []Sig // list of fields
}

type SlatMessage struct {
	*Cmd
}

// return a list of complete signatures
// command_name#_first_label#_second_label#
type Sig struct {
	Raw      string
	Numbered string // command_name#_first_label#_second_label#
	Camel    string // commandName_firstLabel_secondLabel
	crc      uint32
	Type     string
}

func (f *SlotMessage) Format() (ret string) {
	return "%-20s %-20s = 0x%09xd"
}

func (f *SlatMessage) Format() (ret string) {
	return "%-20s %-20s = %3d"
}

func (f *SlatMessage) Camel() (ret string) {
	return Camel(f.Lede)
}

func (x *Sig) Crc() uint32 {
	if x.crc == 0 {
		h := fnv.New32a()
		io.WriteString(h, x.Raw)
		x.crc = h.Sum32()
	}
	return x.crc
}
