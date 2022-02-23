package text

import (
	"bytes"
	"io"
	"testing"
)

func TestText(t *testing.T) {
	var tests = []string{
		// 1. nothing special
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit.`,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit.`,
		// 2. hard lines
		"\nLorem ipsum dolor\nsit amet,\n\nconsectetur\n\nadipiscing elit.",
		"\nLorem ipsum dolor\nsit amet,\n\nconsectetur\n\nadipiscing elit.",
		// 3. soft new lines
		"\rLorem ipsum dolor\rsit amet,\r\rconsectetur\n\radipiscing elit.",
		"Lorem ipsum dolor\nsit amet,\nconsectetur\nadipiscing elit.",
		// 4. soft paragraphs
		"\vLorem ipsum dolor\vsit amet,\v\vconsectetur\n\vadipiscing elit.",
		"Lorem ipsum dolor\n\nsit amet,\n\nconsectetur\n\nadipiscing elit.",
		// 5. bold
		"nadipiscing <b>elit</b>",
		"nadipiscing **elit**",
		// 6. italic
		"nadipiscing <i>elit</i>",
		"nadipiscing *elit*",
		// 7. strike
		"nadipiscing <s>elit</strike>",
		"nadipiscing ~~elit~~",
		// 8. underline
		"nadipiscing <u>elit</u>",
		"nadipiscing __elit__",
		// 9. divider
		"nadipiscing<hr>elit",
		"nadipiscing\n-----------------------------\nelit",
		// 10. hard lines
		"<br>Lorem ipsum dolor<br>sit amet,<br><br>consectetur<br><br>adipiscing elit.",
		"\nLorem ipsum dolor\nsit amet,\n\nconsectetur\n\nadipiscing elit.",
		// 11. soft new lines
		"<wbr>Lorem ipsum dolor<wbr>sit amet,<wbr><wbr>consectetur<br><wbr>adipiscing elit.",
		"Lorem ipsum dolor\nsit amet,\nconsectetur\nadipiscing elit.",
		// 12. soft paragraphs
		"<p>Lorem ipsum dolor<p>sit amet,<p><p>consectetur<br><p>adipiscing elit.",
		"Lorem ipsum dolor\n\nsit amet,\n\nconsectetur\n\nadipiscing elit.",
		// 13. lists
		`Lorem` +
			`<ul>` +
			`<li>ipsum</li>` +
			`<li>dolor sit</li>` +
			`<ol>` +
			/**/ `<li>amet</li>` +
			`</ol>` +
			`<li>consectetur adipiscing</li>` +
			`</ul>` +
			`elit.`,
		//
		`Lorem
  - ipsum
  - dolor sit
    - amet
  - consectetur adipiscing
elit.`,
		// 14. malformed ( no change )
		"if x < 5; x= x<<5; < br >; < br/>; <123><></>",
		"if x < 5; x= x<<5; < br >; < br/>; <123><></>",
		// 15. unknown tag
		"<beep><bop>",
		"<beep><bop>",
		// 16. explicitly closing self-closing tags should be eaten silently with no change to the output
		"</p></br></wbr>",
		"",
	}
	//
	for i, cnt := 0, len(tests); i < cnt; i += 2 {
		var buf bytes.Buffer
		test, want := tests[i], tests[i+1]
		which := i/2 + 1
		if n, e := io.WriteString(Html2Text(&buf), test); e != nil {
			t.Fatal(e)
		} else if wantLen := len(test); n != wantLen {
			t.Fatal(which, "mismatched count", n, "!=", wantLen)
		} else if res := buf.String(); res != want {
			t.Fatal(res)
		} else {
			t.Log("okay", which)
		}
	}
}
