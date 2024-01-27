package compact

// constants for the markup maps in the dl.
const (
	Comment = "comment"
	Markup  = "--"
)

// read a user comment from markup, normalizing it as an array of strings
func UserComment(markup map[string]any) (ret []string) {
	switch cmt := markup[Comment].(type) {
	case string:
		ret = []string{cmt}
	case []string:
		ret = cmt
	case []interface{}:
		lines := make([]string, len(cmt))
		for i, el := range cmt {
			if str, ok := el.(string); !ok {
				lines = nil
				break
			} else {
				lines[i] = str
			}
		}
		ret = lines
	}
	return
}
