package story

//         {
//           "--": "Give a kind a trait",
//           "Spec:slots:groups:with flow:": [
//             "certainties",
//             "story_statement",
//             "properties",
//             {
//               "Flow:uses:": [
//                 "certainties",
//                 [
//                   {
//                     "Term:name:": [
//                       "_",
//                       "plural_kinds"
//                     ]
//                   },
//                   {
//                     "Term:": "are_being"
//                   },
//                   {
//                     "Term:": "certainty"
//                   },
//                   {
//                     "Term:type:": [
//                       "trait",
//                       "text"
//                     ]
//                   }
//                 ]
//               ]
//             }
//           ]
//         },
// { --> from
//   "Certainties:areBeing:certainty:trait:": [
//     "rooms",
//     "are",
//     "usually",
//     "lit"
//   ]
// },

// // horses are usually fast.
// func (op *Certainties) Weave(cat *weave.Catalog) (err error) {
// 	certaintiesNotImplemented.PrintOnce()
// 	// if certainty, e := op.Certainty.ImportString(k); e != nil {
// 	// 	err = e
// 	// } else if trait, e := NewTrait(w, op.Trait); e != nil {
// 	// 	err = e
// 	// } else if kind, e := NewPluralKinds(w, op.PluralKinds); e != nil {
// 	// 	err = e
// 	// } else {
// 	// 	k.NewCertainty(certainty, trait, kind)
// 	// }
// 	return
// }

// var certaintiesNotImplemented eph.PrintOnce = "certainties not implemented"

//func (op *Certainty) ImportString(cat *weave.Catalog) (ret string, err error) {
// 	if str, ok := composer.FindChoice(op, op.Str); !ok {
// 		err = ImportError(op, errutil.Fmt("%w %q", InvalidValue, op.Str))
// 	} else {
// 		ret = str
// 	}
// 	return
// }

//  {
//           "--": "Whether an trait applies to a kind of noun.",
//           "Spec:groups:with str:": [
//             "certainty",
//             "properties",
//             {
//               "Str exclusively:uses:": [
//                 true,
//                 [
//                   {
//                     "Option:": "usually"
//                   },
//                   {
//                     "Option:": "always"
//                   },
//                   {
//                     "Option:": "seldom"
//                   },
//                   {
//                     "Option:": "never"
//                   }
//                 ]
//               ]
//             }
//           ]
//         },
