package gestalt

// matches one pre-existing kind.
// an article can precede every kind
//
// ex. "[the] (container...)"
type Kind struct{}

func (*Kind) Match(q Query, cs []InputState) (ret []InputState) {
	for _, in := range cs {
		ws := in.Words()
		if det := q.SkipArticle(ws); det >= 0 {
			if kind, width := q.FindKind(ws[det:]); width > 0 {
				out := in.Next(det + width)
				out.AddResult(TypeKind, kind)
				ret = append(ret, out)
			}
		}
	}
	return
}
