package web

import (
	"mime"
	"net/url"
	"strconv"
	"strings"
)

// Given a port and a path return a url string `http://path...:port[/]`
func Endpoint(port int, parts ...string) (ret string) {
	ret = ":" + strconv.Itoa(port)
	if len(parts) > 0 {
		u := url.URL{Scheme: "http", Host: parts[0] + ret, Path: strings.Join(parts[1:], "/")}
		if len(parts) > 0 {
			ret = u.String() + "/"
		}
	}
	return
}

// see: https://github.com/golang/go/issues/32350
// https://go-review.googlesource.com/c/go/+/406894/ will be fixed in 1.19
func init() {
	_ = mime.AddExtensionType(".js", "text/javascript")
}
