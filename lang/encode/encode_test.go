package encode_test

import (
	_ "embed"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/shortcut"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/test/debug"
)

func TestEncodeStory(t *testing.T) {
	var enc encode.Encoder // story doesnt have its own custom encoder.
	if n, e := enc.Customize(shortcut.Encoder).Encode(&debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else {
		var b strings.Builder
		if e := files.JsonEncoder(&b, files.RawJson).Encode(n); e != nil {
			t.Fatal(e)
		} else if str := b.String(); str != debug.FactorialJs {
			t.Fatal(str)
		}
	}
}
