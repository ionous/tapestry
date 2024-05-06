package qna

type qkey struct {
	group, target, field string
}

func makeKey(group, target, field string) qkey {
	return qkey{group, target, field}
}

type cache struct {
	store       map[qkey]cachedValue
	cacheErrors bool
}

type cacheMap map[qkey]cachedValue

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

func (c *cache) cache(build func() (any, error), group, target, field string) (ret any, err error) {
	key := makeKey(group, target, field)
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
