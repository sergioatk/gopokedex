package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 1 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "http://example.com",
			val: []byte("test_data"),
		},
		{
			key: "http://example.com/path",
			val: []byte("test_more_data"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key %s", c.key)
				return
			}

			if string(val) != string(c.val) {
				t.Errorf("expected to find value %s", c.key)
				return
			}
		})
	}
}

func TestReadLoop(t *testing.T) {
	const interval = 1 * time.Second
	const waitTime = interval + time.Millisecond*10

	cache := NewCache(interval)

	cacheKey := "http://example.com/read-loop"

	cache.Add(cacheKey, []byte("some_value"))

	_, ok := cache.Get(cacheKey)
	if !ok {
		t.Errorf("expected to find key %s", cacheKey)
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(cacheKey)
	if ok {
		t.Errorf("expected not to find key %s", cacheKey)
	}
}
