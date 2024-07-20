package match_test

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/match"
	"github.com/ionous/tell/charm"
)

func TestTokenTerminal(t *testing.T) {
	if vs, e := match.TokenizeString(`"hello." world`); e != nil {
		t.Fatal(e)
	} else if e := compareTokens(vs, []match.TokenValue{
		{Token: match.Quoted, Value: "hello."},
		{Token: match.Stop, Value: "."},
		{Token: match.String, Value: "world"},
	}); e != nil {
		t.Fatal("mismatch", e)
	}
}

// test (at least) one of each of the possible tokens produced
func TestTokens(t *testing.T) {
	tests := []any{
		// token, string to parse, result:
		/*1*/ match.String, `23`, "23",
		/*2*/ match.String, `hello`, "hello",
		/*3*/ match.Stop, `.`, ".",
		/*4*/ match.Stop, `!`, "!",
		/*5*/ match.Comma, `,`, ",",

		match.Parenthetical,
		`( hello world )`,
		`hello world`,

		// ----------
		match.Quoted,
		`"hello\\world"`,
		`hello\world`,

		// ----------
		match.Quoted,
		"`" + `hello\\world` + "`",
		`hello\\world`,

		// -----
		match.Comment, "# comment", "comment",
		// 		/*9*/ match.Key, "-", "",
		// 		/*10*/ match.Key, "hello:world:", "hello:world:",
		// 		// make sure dash numbers are treated as negative numbers
		// 		/*11*/ match.Number, `-5`, -5,
		// ----------
		match.Quoted,
		`"hello
doc"`,
		`hello doc`,
		// -------------
		match.Quoted,
		strings.Join([]string{
			"'''",
			"hello",
			"line",
			"'''"}, "\n"),
		`hello
line`,
	}

	// test all of the above in both the same and separate buffers
	// at the very least it helps to validate tokens must be separated by whitespace.
	var combined results
	run := match.NewTokenizer(&combined)
	var combinedReader strings.Reader
	combinedParser := charm.MakeParser(&combinedReader)

	for i, whichTest := 0, 1; i < len(tests); i, whichTest = i+3, whichTest+1 {
		wantType := tests[i+0].(match.Token)
		testStr := tests[i+1].(string)
		wantVal := tests[i+2]
		if e := testToken(wantType, testStr, wantVal); e != nil {
			t.Logf("failed single %d: %s", whichTest, e)
			t.Fail()
		} else {
			sep := " "
			if wantType == match.Comment {
				sep = "\n" // comments have to be ended with a newlne
			}
			combinedReader.Reset(testStr + sep) // lets keep the party going
			if next, e := combinedParser.Parse(run); e != io.EOF {
				t.Logf("failed combine parse %d: %s", whichTest, e)
				t.Fail()
			} else {
				last := combined[len(combined)-1]
				if e := compareToken(last, match.TokenValue{Token: wantType, Value: wantVal}); e != nil {
					t.Logf("failed combine compare %d: %s", whichTest, e)
					t.Fail()
				} else {
					run = next
				}
			}
		}
	}
}

func testToken(tokenType match.Token, testStr string, tokenValue any) (err error) {
	var pairs results
	run := match.NewTokenizer(&pairs)
	if e := parseString(testStr+"\n", run); e != io.EOF {
		err = compareValue(e, tokenValue)
	} else if cnt := len(pairs); cnt == 0 {
		err = errors.New("didn't collect any tokens")
	} else {
		last := pairs[cnt-1]
		if e := compareValue(last.Pos, match.Pos{}); e != nil {
			err = e
		} else {
			err = compareToken(last, match.TokenValue{Token: tokenType, Value: tokenValue})
		}
	}
	return
}

// always returns an error; io.EOF means all the input was consumed.
func parseString(str string, state charm.State) (err error) {
	p := charm.MakeParser(strings.NewReader(str))
	_, err = p.Parse(state)
	return
}

type results []match.TokenValue

func (res *results) Decoded(tv match.TokenValue) (_ error) {
	(*res) = append((*res), tv)
	return
}

func compareTokens(have, want []match.TokenValue) (err error) {
	if a, b := len(have), len(want); a != b {
		err = fmt.Errorf("mismatched lengths; have %d, want %d", a, b)
	} else {
		for i, got := range have {
			if e := compareToken(got, want[i]); e != nil {
				err = fmt.Errorf("mismatched at %d %w", i, e)
				break
			}
		}
	}
	return
}

func compareToken(have, want match.TokenValue) (err error) {
	if have.Token != want.Token {
		err = fmt.Errorf("mismatched types want: %s, have: %s", want.Token, have.Token)
	} else {
		err = compareValue(have.Value, want.Value)
	}
	return
}

func compareValue(have any, want any) (err error) {
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
