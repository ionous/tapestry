package files

import (
	"encoding/json"
	"io"

	"github.com/ionous/tell"
	"github.com/ionous/tell/encode"
)

// matches https://pkg.go.dev/encoding/json#Encoder.Encode
type PlainEncoder interface {
	Encode(from any) error
}

// matches https://pkg.go.dev/encoding/json#Decoder.Decode
type PlainDecoder interface {
	Decode(into any) error
}

type JsonFlags int

const (
	RawJson JsonFlags = 0
	Pretty  JsonFlags = 1 << iota
	EscapeHtml
)

func MakeJsonFlags(pretty, escape bool) (ret JsonFlags) {
	if pretty {
		ret |= Pretty
	}
	if escape {
		ret |= EscapeHtml
	}
	return
}

func JsonEncoder(w io.Writer, f JsonFlags) PlainEncoder {
	pretty := f&Pretty != 0
	if !pretty {
		w = &noNewLine{out: w}
	}
	js := json.NewEncoder(w)
	if pretty {
		js.SetIndent("", "  ")
	}
	js.SetEscapeHTML(f&EscapeHtml != 0)
	return js
}

func TellEncoder(w io.Writer) PlainEncoder {
	enc := tell.
		NewEncoder(w).
		SetSequencer(encode.MakeSequence, encode.NoComments).
		SetMapper(makeMapping, makeComments)
	return enc
}
