package web

import (
	"strconv"

	"github.com/ionous/errutil"
)

// implements the flag.Value interface for reading ports from the commandline
type Port int

func (p Port) Int() (ret int) {
	return int(p)
}

func (p Port) String() string {
	return strconv.Itoa(p.Int())
}

func (p *Port) Set(s string) (err error) {
	if len(s) == 0 || s == "false" {
		*p = -1 // returns the default port and false
	} else if s == "true" {
		*p = 1 // returns the default port and true
	} else if i, e := strconv.Atoi(s); e != nil {
		err = e // not a number?
	} else if !portIsValid(i) {
		err = errutil.New("expected a port in the range 1024-49151. got", i)
	} else {
		*p = Port(i)
	}
	return
}

func portIsValid(i int) bool {
	return i >= 1024 && i < 49152
}

func (p Port) GetPort(defaultPort int) (ret int, wasValid bool) {
	if n := p.Int(); n == 1 {
		ret, wasValid = defaultPort, true
	} else if !portIsValid(n) {
		ret, wasValid = defaultPort, false
	} else {
		ret, wasValid = n, true
	}
	return
}
