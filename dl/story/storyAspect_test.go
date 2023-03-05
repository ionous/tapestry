package story_test

// func TestEventually(t *testing.T) {
// 	asp := story.DefineTraits{
// 		Aspect: assign.T("test"),
// 		Traits: assign.Ts("yes", "no"),
// 	}
// 	var aspect string
// 	var out []string
// 	ts := chart.MakeEncoder()
// 	if e := ts.Marshal(&asp,
// 		chart.Map(&ts, chart.BlockMap{
// 			story.DefineTraits_Type: chart.KeyMap{
// 				story.DefineTraits_Field_Aspect: func(b jsn.Block, v interface{}) (err error) {
// 					aspect = *v.(*string)
// 					return
// 				},
// 				story.DefineTraits_Field_TraitPhrase: func(jsn.Block, interface{}) (err error) {
// 					ts.PushState(story.EveryValueOf(&ts, story.Trait_Type, func(v interface{}) (err error) {
// 						trait := *v.(*string)
// 						out = append(out, aspect+":"+trait)
// 						return
// 					}))
// 					return
// 				},
// 			}})); e != nil {
// 		t.Fatal(e)
// 	} else if diff := pretty.Diff(out, []string{"test:yes", "test:no"}); len(diff) > 0 {
// 		t.Fatal(diff)
// 	} else {
// 		t.Log("okay")
// 	}
// }

// func TestEndBlock(t *testing.T) {
// 	asp := story.Certainties{
// 		PluralKinds: story.PluralKinds{Str: "test"},
// 	}
// 	found := false
// 	ts := chart.MakeEncoder()
// 	if e := ts.Marshal(&asp,
// 		chart.Map(&ts, chart.BlockMap{
// 			story.Certainties_Type: chart.KeyMap{
// 				chart.BlockEnd: func(b jsn.Block, v interface{}) (err error) {
// 					cs := b.(jsn.FlowBlock).GetFlow().(*story.Certainties) // ick
// 					found = cs.PluralKinds.Str == "test"
// 					return
// 				},
// 			}})); e != nil {
// 		t.Fatal(e)
// 	} else if !found {
// 		t.Fatal("end not found")
// 	}
// }

// // visit every value matching the typeName
// func EveryValueOf(m *chart.Machine, typeName string, fn func(interface{}) error) chart.State {
// 	var blocks chart.BlockStack
// 	return &chart.StateMix{
// 		OnBlock: func(b jsn.Block) (err error) {
// 			blocks.Push(b)
// 			return
// 		},
// 		OnKey: func(lede, key string) (err error) {
// 			return
// 		},
// 		OnValue: func(valType string, val interface{}) (err error) {
// 			if valType == typeName {
// 				err = fn(val)
// 			}
// 			return
// 		},
// 		OnEnd: func() {
// 			if _, ok := blocks.Pop(); !ok {
// 				m.FinishState("scope") // pop this.
// 			}
// 		},
// 	}
// }
