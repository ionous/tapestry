package generate

// func TestGenerate(t *testing.T) {
// 	// tell files:
// 	if tests, e := readTests(); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		for _, n := range tests {
// 			if name := n.name; strings.HasPrefix(name, focus) || focus == "" {
// 				for i, msg := range n.msgs {
// 					var str strings.Builder
// 					if e := generateMsg(&str, msg); e != nil {
// 						t.Fatal(name, i, e)
// 					} else {
// 						// tabs and newlines make it hard to test the generated code.
// 						// ignore it for now.
// 						t.Log(name, i, "\n", str.String())
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// var focus = "flow."

// func generateMsg(w io.Writer, msg compact.Message) (err error) {
// 	var pack groupSearch
// 	if tmp, e := genTemplates(&pack); e != nil {
// 		err = e
// 	} else {
// 		var gc groupContent
// 		gen := MakeGenerator( Generator{w: w, tmp: tmp}
// 		if e := readSpec(&gc, msg); e != nil {
// 			err = e
// 		} else {
// 			pack.list = append(pack.list, Group{"", gc})
// 			err = gen.write(gc)
// 		}
// 	}
// 	return
// }

// func readTests() (ret []test, err error) {
// 	err = fs.WalkDir(testdata, ".", func(path string, d fs.DirEntry, e error) (err error) {
// 		if e != nil {
// 			err = e
// 		} else if !d.IsDir() { // the first dir we get is "."
// 			if fp, e := testdata.Open(path); e != nil {
// 				err = e
// 			} else if raw, e := files.ReadRawTell(fp); e != nil {
// 				err = e
// 			} else if msgs, e := parseMessages(raw); e != nil {
// 				err = e
// 			} else {
// 				ret = append(ret, test{
// 					name: d.Name(),
// 					msgs: msgs,
// 					// expect: expect,
// 				})
// 			}
// 		}
// 		return
// 	})
// 	return
// }

// type test struct {
// 	name string
// 	msgs []compact.Message
// 	// expect string
// }

// //go:embed testdata/*.tells
// var testdata embed.FS

// if msgs[1].Key != "Expect:" {
// 				err = fmt.Errorf("missing expectation")
// 			} else if expect, ok := msgs[1].Args[0].(string); !ok {
// 				err = fmt.Errorf("missing expectation")
// 			} else
