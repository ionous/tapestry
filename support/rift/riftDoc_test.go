package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/rift"
)

func TestDoc(t *testing.T) {
	if e := match(t,
		"multiple collections",
		testValue(`cartParams:
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
