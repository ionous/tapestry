package cmdcompact

// define  a custom spec encoder.
// var customSpecEncoder cout.CustomFlow = nil
// var customStoryFlow = story.CompactEncoder
// var customStorySlot cout.CustomSlot = nil

// example removing "trim" for underscore names
// func init() {
// 	customSpecEncoder = func(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
// 		switch op := flow.GetFlow().(type) {
// 		case *spec.FlowSpec:
// 			if op.Trim {
// 				if len(op.Terms) == 0 {
// 					panic("empty terms " + op.Name)
// 				}
// 				if op.Terms[0].Name != "" {
// 					panic("unexpected name " + op.Name + " " + op.Terms[0].Name)
// 				}
// 				if op.Terms[0].Key == "" {
// 					panic("unexpected key " + op.Name)
// 				} else {
// 					op.Terms[0].Name = op.Terms[0].Key
// 					op.Terms[0].Key = "_"
// 				}
// 				op.Trim = false
// 			}
// 		}
// 		// we haven't serialized it -- just poked at its memory
// 		return chart.Unhandled("no custom encoder")
// 	}
// }

//// install a custom encoder to rewrite things
//func init() {
//	customStorySlot = func(m jsn.Marshaler, block jsn.SlotBlock) error {
//		if slot, ok := block.GetSlot(); ok {
//			switch op := slot.(type) {
//			case *story.PatternActions:
//				a := op.Rules
//				for i := len(a)/2 - 1; i >= 0; i-- {
//					opp := len(a) - 1 - i
//					a[i], a[opp] = a[opp], a[i]
//				}
//				block.SetSlot(&story.ExtendPattern{
//					PatternName: assign.T(op.PatternName),
//					Locals:      op.Locals,
//					Rules:       a,
//					Markup:      op.Markup,
//				})
//			}
//		}
//		return chart.Unhandled("always returns unhandled")
//	}
//}
