package fifo

import "testing"

func TestSetGet(t *testing.T) {
	cache := New(24, nil)

	k, v := "k1", 1
	cache.Set("k1", v)
	got := cache.Get(k)

	if got != v {
		t.Errorf("got %d want %d", got.(int), v)
	}

	cache.Del(k)
	l := cache.Len()
	if l != 0 {
		t.Errorf("length got %d want 0", l)
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := New(16, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{"k1", "k2"}
	if keys[0] != expected[0] || keys[1] != expected[1] {
		t.Errorf("got %v want %v", keys, expected)
	}
	//t.Errorf("got %d want %d", cache.maxBytes, cache.usedBytes)
}
