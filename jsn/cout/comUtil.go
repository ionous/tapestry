package cout

import (
	"encoding/json"
	"io"
	"strings"

	"git.sr.ht/~ionous/iffy/jsn"
)

func Marshal(cmd jsn.Marshalee, customFlow CustomFlow) (ret string, err error) {
	var out strings.Builder
	if e := marshal(&out, cmd, customFlow, false); e != nil {
		err = e
	} else {
		ret = out.String()
	}
	return
}

func MarshalIndent(out io.Writer, cmd jsn.Marshalee, customFlow CustomFlow) error {
	return marshal(out, cmd, customFlow, true)
}

func marshal(out io.Writer, cmd jsn.Marshalee, customFlow CustomFlow, indent bool) (err error) {
	if data, e := Encode(cmd, customFlow); e != nil {
		err = e
	} else {
		js := json.NewEncoder(&noNewLine{out: out})
		if indent {
			js.SetIndent("", "  ")
		}
		js.SetEscapeHTML(false)
		if e := js.Encode(data); e != nil {
			err = e
		}
	}
	return
}

// https://github.com/golang/go/issues/37083
// or, they could have just provided an Set<Option> like SetEscapeHTML...
type noNewLine struct {
	out     io.Writer
	pending bool
}

func (n *noNewLine) Write(p []byte) (ret int, err error) {
	if cnt := len(p); cnt > 0 {
		if n.pending {
			n.out.Write([]byte{newline})
			n.pending = false
		}
		pending := p[cnt-1] == newline
		if pending {
			n.pending = pending
			cnt--
		}
		if cnt == 0 {
			ret = 1
		} else {
			c, e := n.out.Write(p[:cnt])
			if pending && c > 0 {
				c++
			}
			ret, err = c, e
		}
	}
	return
}

const newline byte = '\n'
