package qna

type cache map[uint64]cachedValue

type cachedValue struct {
	v interface{}
	e error
}

func (c *cache) cache(build func() (interface{}, error), args ...string) (ret interface{}, err error) {
	if len(args) == 0 {
		panic("missing key for cache")
	}
	key := makeKey(args...)
	if n, ok := (*c)[key]; ok {
		ret, err = n.v, n.e
	} else {
		var n cachedValue
		if v, e := build(); e != nil {
			err, n.e = e, e
		} else {
			ret, n.v = v, v
		}
		(*c)[key] = n
	}
	return
}
