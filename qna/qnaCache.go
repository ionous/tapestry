package qna

type cache struct {
	store       map[uint64]cachedValue
	cacheErrors bool
}

type cacheMap map[uint64]cachedValue

func makeCache(cacheErrors bool) cache {
	return cache{make(cacheMap), cacheErrors}
}

type cachedValue struct {
	v any
	e error
}

func (c *cache) reset() {
	c.store = make(cacheMap)
}

func (c *cache) cache(build func() (any, error), args ...string) (ret any, err error) {
	if len(args) == 0 {
		panic("key for cache unspecified")
	}
	key := makeKey(args...)
	if n, ok := c.store[key]; ok {
		ret, err = n.v, n.e
	} else {
		if v, e := build(); e == nil {
			c.store[key] = cachedValue{v: v}
			ret = v
		} else {
			err = e
			if c.cacheErrors {
				c.store[key] = cachedValue{e: e}
			}
		}
	}
	return
}
