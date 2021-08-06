package lfu

import "testing"

func TestSet(t *testing.T) {
	cache := New(24, nil)

	k, v := "k1", 1
	cache.Set(k, v)
	got := cache.Get(k).(int)
	if got != v {
		t.Errorf("length got %d want %d", got, v)
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := New(32, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Get("k1")
	cache.Get("k2")
	cache.Set("k3", 3)
	cache.Set("k4", 4)

	expected := []string{"k2", "k3"}
	if keys[0] != expected[0] || keys[1] != expected[1] {
		t.Errorf("got %v want %v", keys, expected)
	}
}
