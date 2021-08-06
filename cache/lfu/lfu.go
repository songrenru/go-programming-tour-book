package lfu

import (
	"container/heap"
	"github.com/songrenru/cache"
)

// 需求：设计一个lfu的cache

// 1. heap排序
// 2. 定义struct，map + k,v + count

type lfu struct {
	cache map[string]*entry
	hh *entryHeap

	maxBytes int
	usedBytes int

	onEvicted func(key string, value interface{})  // 用户对外部资源进行清理
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	h := make(entryHeap, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		hh:     &h,
		cache:     make(map[string]*entry),
		onEvicted: onEvicted,
	}
}

func (l *lfu) Get(key string) interface{} {
	en, ok := l.cache[key]
	if !ok {
		return nil
	}

	l.update(en, en.value)
	return en.value
}

func (l *lfu) Set(key string, value interface{}) {
	en, ok := l.cache[key]
	if ok {
		l.usedBytes = l.usedBytes - en.Len() + cache.CalculLen(value)
		l.update(en, value)
		return
	}

	for l.usedBytes + cache.CalculLen(value) > l.maxBytes {
		l.DelOldest() // 应该先淘汰，再插入，不然会直接弹出了刚插入的元素
	}
	en = &entry{
		key:   key,
		value: value,
		count: 1,
	}
	heap.Push(l.hh, en)
	l.cache[key] = en
	l.usedBytes += en.Len()
}

func (l *lfu) Del(key string) {
	if en, ok := l.cache[key]; ok {
		heap.Remove(l.hh, en.index)
		l.delElement(en)
	}
}

func (l *lfu) DelOldest() {
	en := heap.Pop(l.hh).(*entry)
	l.delElement(en)
}

func (l *lfu) Len() int {
	return len(l.cache)
}

func (l *lfu) update(en *entry, value interface{}) {
	en.count++
	en.value = value
	heap.Fix(l.hh, en.index)
}

func (l *lfu) delElement(en *entry) {
	delete(l.cache, en.key)
	l.usedBytes -= en.Len()

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

type entry struct {
	key string
	value interface{}
	count int
	index int
}

func (e *entry) Len() int {
	return cache.CalculLen(e.value) + 4 + 4
}

type entryHeap []*entry

func (h *entryHeap) Less(i, j int) bool {
	return (*h)[i].count < (*h)[j].count
}

func (h *entryHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].index, (*h)[j].index = (*h)[j].index, (*h)[i].index
}

func (h *entryHeap) Len() int {
	return len(*h)
}

func (h *entryHeap) Pop() (v interface{}) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *entryHeap) Push(v interface{}) {
	en := v.(*entry)
	en.index = len(*h)
	*h = append(*h, en)
}




