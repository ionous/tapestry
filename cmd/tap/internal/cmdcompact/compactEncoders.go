package cmdcompact

import (
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

// define  a custom spec encoder.
var customSpecEncoder cout.CustomFlow = nil
var customStoryEncoder = story.CompactEncoder

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

// install a custom encoder to rewrite things
// func init() {
// 	customStoryEncoder = func(m jsn.Marshaler, flow jsn.FlowBlock) error {
// 		switch op := flow.GetFlow().(type) {
// 		case *story.AspectProperty:
// 			swap(&op.UserComment, &op.Comment)
// 		case *story.BoolProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.NumberProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.NumListProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.RecordProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.RecordListProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.TextListProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.TextProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		}
// 		return core.CompactEncoder(m, flow)
// 	}
// }
