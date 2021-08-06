package lru

import (
	"container/list"
	"github.com/songrenru/cache"
	"time"
)

// 需求：设计一个lru的cache

// 1. db link
// 2. 定义struct，map + k,v + count

type entry struct {
	key string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalculLen(e.value)
}

type lru struct {
	maxBytes int
	usedBytes int

	queue *list.List
	cache map[string]*list.Element

	onEvicted func(key string, value interface{})  // 用户对外部资源进行清理
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	time.Sleep(time.Millisecond)
	return &lru{
		maxBytes:  maxBytes,
		queue:     list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (f *lru) Get(key string) interface{} {
	ele, ok := f.cache[key]
	if !ok {
		return nil
	}

	f.queue.MoveToBack(ele)
	return ele.Value.(*entry).value
}

func (f *lru) Set(key string, value interface{}) {
	ele, ok := f.cache[key]
	if ok {
		f.queue.MoveToBack(ele)
		en := ele.Value.(*entry)
		f.usedBytes = f.usedBytes - en.Len() + cache.CalculLen(value)
		en.value = value
		return
	}

	en := &entry{key: key, value: value}
	ele = f.queue.PushBack(en)
	f.cache[key] = ele
	f.usedBytes += en.Len()
	// 内存管理
	for f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *lru) Del(key string) {
	if ele, ok := f.cache[key]; ok {
		f.delElement(ele)
	}
}

func (f *lru) DelOldest() {
	f.delElement(f.queue.Front())
}

func (f *lru) Len() int {
	return f.queue.Len()
}

func (f *lru) delElement(e *list.Element) {
	f.queue.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)

	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}


