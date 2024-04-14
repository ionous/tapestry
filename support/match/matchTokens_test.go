package flex_test

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/flex"
	"github.com/ionous/tell/charm"
)

// test (at least) one of each of the possible tokens produced
func TestTokens(t *testing.T) {
	tests := []any{
		// token, string to parse, result:
		/*1*/ flex.String, `23`, "23",
		/*2*/ flex.String, `hello`, "hello",
		/*3*/ flex.Stop, `.`, '.',
		/*4*/ flex.Stop, `!`, '!',
		/*5*/ flex.Comma, `,`, ',',

		flex.Parenthetical,
		`( hello world )`,
		`hello world`,

		// ----------
		flex.Quoted,
		`"hello\\world"`,
		`hello\world`,

		// ----------
		flex.Quoted,
		"`" + `hello\\world` + "`",
		`hello\\world`,

		// -----
		flex.Comment, "# comment", "comment",
		// 		/*9*/ flex.Key, "-", "",
		// 		/*10*/ flex.Key, "hello:world:", "hello:world:",
		// 		// make sure dash numbers are treated as negative numbers
		// 		/*11*/ flex.Number, `-5`, -5,
		// ----------
		flex.Quoted,
		`"""
hello
doc
"""`,
		`hello doc`,
		// -------------
		flex.Quoted,
		strings.Join([]string{
			"```",
			"hello",
			"line",
			"```"}, "\n"),
		`hello
line`,
	}

	// test all of the above in both the same and separate buffers
	// at the very least it helps to validate tokens must be separated by whitespace.
	var combined results
	run := flex.NewTokenizer(&combined)

	for i := 0; i < len(tests); i += 3 {
		wantType := tests[i+0].(flex.Type)
		testStr := tests[i+1].(string)
		wantVal := tests[i+2]
		whichTest := 1 + i/3
		if e := testToken(wantType, testStr, wantVal); e != nil {
			t.Logf("failed single %d: %s", whichTest, e)
			t.Fail()
		} else {
			sep := " "
			if wantType == flex.Comment {
				sep = "\n" // comments have to be ended with a newlne
			}
			if next, e := charm.Parse(testStr+sep, run); e != nil {
				t.Logf("failed combine parse %d: %s", whichTest, e)
				t.Fail()
			} else {
				last := combined[len(combined)-1]
				if e := last.compare(wantType, wantVal); e != nil {
					t.Logf("failed combine compare %d: %s", whichTest, e)
					t.Fail()
				} else {
					run = next
				}
			}
		}
	}
}

func testToken(tokenType flex.Type, testStr string, tokenValue any) (err error) {
	var pairs results
	run := flex.NewTokenizer(&pairs)
	if _, e := charm.Parse(testStr+"\n", run); e != nil {
		err = compare(e, tokenValue)
	} else if cnt := len(pairs); cnt == 0 {
		err = errors.New("didn't collect any tokens")
	} else {
		last := pairs[cnt-1]
		if e := compare(last.pos, flex.Pos{}); e != nil {
			err = e
		} else {
			err = last.compare(tokenType, tokenValue)
		}
	}
	return
}

type results []result

type result struct {
	pos        flex.Pos
	tokenType  flex.Type
	tokenValue any
}

func (res *results) Decoded(pos flex.Pos, tokenType flex.Type, tokenValue any) (_ error) {
	(*res) = append((*res), result{pos, tokenType, tokenValue})
	return
}

// compare everything except pos
func (res results) compare(expects results) (err error) {
	if have, want := len(res), len(expects); have != want {
		log.Println(res)
		err = fmt.Errorf("failed test have %d != want %d", have, want)
	} else {
		for k, el := range res {
			want := expects[k]
			if e := el.compare(want.tokenType, want.tokenValue); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func (p result) compare(wantType flex.Type, wantValue any) (err error) {
	if tt := p.tokenType; tt != wantType {
		err = fmt.Errorf("mismatched types want: %s, have: %s", wantType, tt)
	} else {
		err = compare(p.tokenValue, wantValue)
	}
	return
}

func compare(have any, want any) (err error) {
	if haveErr, ok := have.(error); !ok {
		if !reflect.DeepEqual(have, want) {
			err = fmt.Errorf("mismatched want: %v(%T) have: %v(%T)", want, want, have, have)
		}
	} else {
		if expectErr, ok := want.(error); !ok {
			err = fmt.Errorf("failed %v", haveErr)
		} else if !strings.HasPrefix(haveErr.Error(), expectErr.Error()) {
			err = fmt.Errorf("failed %v, expected %v", haveErr, expectErr)
		}
	}
	return
}
