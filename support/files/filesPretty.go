package files

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

// take existing json data and prettify it.
// on error, logs and returns the empty string
func Indent(str string) (ret string) {
	var b bytes.Buffer
	if e := json.Indent(&b, []byte(str), "", "  "); e != nil {
		log.Println("indention error", e)
	} else {
		ret = b.String()
	}
	return
}

// take existing json data and minimize it.
// on error, logs and returns the empty string
func Compact(str string) (ret string) {
	var b bytes.Buffer
	if e := json.Compact(&b, []byte(str)); e != nil {
		log.Println("indention error", e)
	} else {
		ret = b.String()
	}
	return
}

// remove go's trailing newline.
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
