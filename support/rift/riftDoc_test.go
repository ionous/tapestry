package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/rift"
)

func TestDoc(t *testing.T) {
	if e := match(t,
		"multiple collections",
		testValue(`
cartParams:
  - name: "Leg type"
    options:
      - "Round"
      - "Square"
  - name: "Leg width"
    options: 
      - "1.00 inches"
      - "2.00 inches"`),
		rift.MapValues{
			{"cartParams:", []any{
				rift.MapValues{{
					"name:", "Leg type",
				}, {
					"options:", []any{
						"Round",
						"Square",
					},
				}}, rift.MapValues{{
					"name:", "Leg width",
				}, {
					"options:", []any{
						"1.00 inches",
						"2.00 inches",
					}},
				},
			},
			},
		},
	); e != nil {
		t.Fatal(e)
	}
}

func TestNestedMap(t *testing.T) {
	if e := match(t,
		"nested map",
		testValue(`
- Field:
  Next: 5`), []any{
			rift.MapValues{{
				Key: "Field:", Value: nil,
			}, {
				Key: "Next:", Value: 5.0,
			},
			}}); e != nil {
		t.Fatal(e)
	}
}

func TestNestedMaps(t *testing.T) {
	if e := match(t,
		"nested maps",
		testValue(`
- Field:
    Next: 5`), []any{
			rift.MapValues{{
				Key: "Field:", Value: rift.MapValues{{
					Key: "Next:", Value: 5.0,
				}},
			}}}); e != nil {
		t.Fatal(e)
	}
}

// in yaml, inline nested maps are invalid
// should they be here too?
// to do, i think Value would need to examine history
// either sniffing prior types or through a flag (ex. require newlines)
// that it can send into NewMapping
func TestMaybeBad(t *testing.T) {
	if e := match(t,
		"inline maps",
		testValue(`- Field: Next: 5`), []any{
			rift.MapValues{{
				Key: "Field:", Value: rift.MapValues{{
					Key: "Next:", Value: 5.0,
				}},
			}}}); e != nil {
		t.Fatal(e)
	}
}
