package compact

import (
	"fmt"
	"path"
)

// the compact format stores source file info inside the metadata.
type Source struct {
	File    string // with extension
	Path    string // enough to locate the file
	Line    int    // a zero-offset printed as one-offset.
	Comment string
}

func (p Source) String() (ret string) {
	if len(p.File) > 0 {
		ret = fmt.Sprintf("%d:%s(%s)", p.Line+1, p.File, p.Path)
	}
	return
}

func MakeSource(m map[string]any) Source {
	var pos Source
	if len(m) > 0 {
		pos.Comment = JoinComment(m)
		if at, ok := m[Position].([]int); ok {
			pos.Line = at[1]
		}
		if at, ok := m[File].(string); ok {
			file := path.Base(at) // extract the file from shared/something.tell
			pos.File = file
			if full, part := len(at), len(file); full > part {
				pos.Path = at[:full-(part+1)] // skip trailing slash before the filename
			}
		}
	}
	return pos
}
