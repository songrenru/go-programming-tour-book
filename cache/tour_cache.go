package cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type TourCache struct {
	mainCache *safeCache
	getter Getter
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.get(key)
	if val != nil {
		return val
	}

	if t.getter != nil {
		val = t.getter.Get(key)
		if val != nil {
			t.mainCache.set(key, val)
			return val
		}
	}

	return nil
}

func (t *TourCache) Stat() *Stat {
	return t.mainCache.stat()
}


