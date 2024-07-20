package match

import "github.com/ionous/tell/charm"

// expose for testing
func DecodeTestDoc(n AfterDocument) charm.State {
	includeComments, lineOfs := false, 0
	return decodeDoc(includeComments, lineOfs, n)
}
