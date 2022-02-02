package bgen_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
)

// story lines should be a block with no output, and one stacking input
// the stacks should all use the "stacked_kinds_of_kind" type
func TestStoryLines(t *testing.T) {
	if e := testBlocks(&story.StoryLines{
		Lines: []story.StoryStatement{
			&story.KindsOfKind{
				PluralKinds:  story.PluralKinds{"cats"},
				SingularKind: story.SingularKind{"animal"},
			},
			&story.KindsOfKind{
				PluralKinds:  story.PluralKinds{"cats"},
				SingularKind: story.SingularKind{"animal"},
			},
		},
	}, `{
  "id": "test-1",
  "type": "story_lines",
  "extraState": {
    "LINES": 2
  },
  "inputs": {
    "LINES": {
      "block": {
        "id": "test-2",
        "type": "_kinds_of_kind_stack",
        "extraState": {
          "PLURAL_KINDS": 1,
          "SINGULAR_KIND": 1
        },
        "fields": {
          "PLURAL_KINDS": "cats",
          "SINGULAR_KIND": "animal"
        },
        "next": {
          "block": {
            "id": "test-3",
            "type": "_kinds_of_kind_stack",
            "extraState": {
              "PLURAL_KINDS": 1,
              "SINGULAR_KIND": 1
            },
            "fields": {
              "PLURAL_KINDS": "cats",
              "SINGULAR_KIND": "animal"
            }
          }
        }
      }
    }
  }
}`); e != nil {
		t.Fatal(e)
	}
}
