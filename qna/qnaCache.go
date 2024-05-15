package qna

type qkey struct {
	domain, target, field string
}

func makeKey(domain, target, field string) qkey {
	return qkey{domain, target, field}
}

type cache struct {
	store       map[qkey]any
	cacheErrors bool
}

type cacheMap map[qkey]any

func makeCache(cacheErrors bool) cache {
	return cache{make(cacheMap), cacheErrors}
}

func (c *cache) reset() {
	c.store = make(cacheMap)
}

func (c *cache) ensure(key qkey, build func() (any, error)) (ret any, err error) {
	if v, ok := c.store[key]; ok {
		if e, ok := v.(error); ok {
			err = e
		} else {
			ret = v
		}
	} else {
		if v, e := build(); e == nil {
			c.store[key] = v
			ret = v
		} else {
			err = e
			if c.cacheErrors {
				c.store[key] = e
			}
		}
	}
	return
}
