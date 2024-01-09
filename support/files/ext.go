package files

import (
	"path/filepath"
)

type Extension int

//go:generate stringer -type=Extension -linecomment
const (
	BlockExt   Extension = iota + 1 // .block
	CompactExt                      // .if
	SpecExt                         // .ifspecs
	TellSpec                        // .tells
	TellStory                       // .tell
)

func Ext(path string) (ret Extension) {
	ext := filepath.Ext(path)
	for i := 1; i < len(_Extension_index); i++ {
		if n := Extension(i); n.String() == ext {
			ret = n
			break
		}
	}
	return
}

func (ext Extension) Spec() bool {
	return ext == TellSpec || ext == SpecExt
}

func (ext Extension) Story() bool {
	return ext == CompactExt || ext == TellStory
}

// blockly's format.
func (ext Extension) Blockly() bool {
	return ext == BlockExt
}

// uses the format described by github.com/ionous/tell.
func (ext Extension) Tell() bool {
	return ext == TellSpec || ext == TellStory
}

// json message format ( excludes the deprecated detailed format )
func (ext Extension) Json() bool {
	return ext == CompactExt || ext == SpecExt
}
