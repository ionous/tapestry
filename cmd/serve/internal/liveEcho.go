package serve

// read from stdin
// go func() {
// r := bufio.NewReader(os.Stdin)
// 	for {
// 		if in, e := r.ReadString('\n'); e != nil {
// 			break
// 		} else if cnt := len(in); cnt > 0 {
// 			msgs <- in[:cnt-1] // trim the newline.
// 		}
// 	}
// }()
