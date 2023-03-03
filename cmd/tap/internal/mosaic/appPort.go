package mosaic

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
	if s == "false" {
		*p = 0 // same as if nothing was specified
	} else if s == "true" {
		*p = 1
	} else if i, e := strconv.Atoi(s); e != nil {
		err = e
	} else if !portIsValid(i) {
		err = errutil.New("expected a port in the range 1024-49151. got", i)
	} else {
		*p = Port(i)
	}
	return
}

func portIsValid(i int) bool {
	return i == 0 || i == 1 || (i >= 1024 && i < 49152)
}

func (p Port) GetPort(defaultPort int) (ret int, wasSet bool) {
	if n := p.Int(); n == 0 || !portIsValid(n) {
		ret, wasSet = defaultPort, false
	} else if n == 1 {
		ret, wasSet = defaultPort, true
	} else {
		ret, wasSet = n, true
	}
	return
}
