package query

type Key struct {
	Domain, Target, Field string
}

func MakeKey(domain, target, field string) Key {
	return Key{domain, target, field}
}

// for logs and errors
func (k Key) String() string {
	return k.Domain + "::" + k.Target + "." + k.Field
}

type Cache struct {
	CacheMap
	cacheErrors bool
}

func (c *Cache) Reset() {
	c.CacheMap = make(CacheMap)
}

// things that implement TextMarshaler will save.
type CacheMap map[Key]any

func MakeCache(cacheErrors bool) Cache {
	return Cache{make(CacheMap), cacheErrors}
}
func (c *Cache) Get(k Key) (ret any, okay bool) {
	ret, okay = c.CacheMap[k]
	return
}
func (c *Cache) Store(k Key, d any) {
	c.CacheMap[k] = d
}

func (c *Cache) Ensure(key Key, build func() (any, error)) (ret any, err error) {
	if v, ok := c.CacheMap[key]; ok {
		if e, ok := v.(error); ok {
			err = e
		} else {
			ret = v
		}
	} else {
		if v, e := build(); e == nil {
			c.CacheMap[key] = v
			ret = v
		} else {
			err = e
			if c.cacheErrors {
				c.CacheMap[key] = e
			}
		}
	}
	return
}
