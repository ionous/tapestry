package web

import (
	"net/http"
	"strconv"
)

// turn a status code into an error
type Status int

func (s Status) Error() (ret string) {
	c := int(s)
	if res := http.StatusText(c); len(res) > 0 {
		ret = "unknown http status code " + strconv.Itoa(c)
	} else {
		ret = res
	}
	return
}
