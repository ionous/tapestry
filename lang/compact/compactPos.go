package compact

import (
	"fmt"
	"path"
)

// the compact format stores source file info inside the metadata.
type Source struct {
	File    string // enough to locate the data
	Line    int    // a zero-offset printed as one-offset.
	Comment string
}

// return the source position as line:base(path)
func (p Source) String() (ret string) {
	if len(p.File) > 0 {
		var rel string
		base := path.Base(p.File) // extract the file from shared/something.tell
		if p.Line < 0 {
			rel = "internal"
		} else if full, part := len(p.File), len(base); full > part {
			rel = p.File[:full-(part+1)] // skip trailing slash before the filename
		}
		// padding the number sorts better
		ret = fmt.Sprintf("%s, %s, %3d", base, rel, p.Line+1)
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
		if file, ok := m[File].(string); ok {
			pos.File = file
		}
	}
	return pos
}
