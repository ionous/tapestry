package web

import (
	"mime"
)

// see: https://github.com/golang/go/issues/32350
// https://go-review.googlesource.com/c/go/+/406894/ will be fixed in 1.19
func init() {
	_ = mime.AddExtensionType(".js", "text/javascript")
}
