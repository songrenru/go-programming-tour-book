package cache

import (
	"fmt"
	"runtime"
	"sync"
)

type Cache interface {
	Get(key string) interface{};
	Set(key string, value interface{});
	Del(key string)
	DelOldest()
	Len() int
}

type Value interface {
	Len() int
}

func CalculLen(value interface{}) int {
	var n int
	switch v := value.(type) {
	case Value:
		n = v.Len()
	case string:
		if runtime.GOARCH == "amd64" {
			n = 16 + len(v)
		} else {
			n = 8 + len(v)
		}
	case uint8, int8, bool:
		n = 1
	case uint16, int16:
		n = 2
	case int32, uint32, float32:
		n = 4
	case int64, uint64, float64:
		n = 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			n = 8
		} else {
			n = 4
		}
	case complex64:
		n = 8
	case complex128:
		n = 16
	default:
		panic(fmt.Sprintf("%T is not implement cache.Value", value))
	}

	return n
}

const DefaultMaxBytes = 1 << 29

type safeCache struct {
	m sync.RWMutex
	cache Cache

	nhit, nget int
}

func NewSafeCache(cache Cache) *safeCache {
	return &safeCache{cache: cache}
}

func (s *safeCache) set(key string, value interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	s.cache.Set(key, value)
}

func (s *safeCache) get(key string) interface{} {
	s.m.Lock() // todo 这里感觉也该用写锁，如LRU的cache，get也会操作数据
	defer s.m.Unlock()

	s.nget++
	value := s.cache.Get(key)
	if value == nil {
		return nil
	}

	s.nhit++
	return value
}

func (s *safeCache) stat() *Stat {
	s.m.RLock()
	defer s.m.RUnlock()
	return &Stat{
		NHit: s.nhit,
		NGet: s.nget,
	}
}

type Stat struct {
	NHit, NGet int
}
