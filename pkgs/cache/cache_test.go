package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/f4tal-err0r/discord_faas/pkgs/cache"
)

func TestCache_SetAndGet(t *testing.T) {
	c := cache.New()

	// Test setting and getting a value
	c.Set("foo", "bar", 5*time.Second)
	if val, exists := c.Get("foo"); !exists || val != "bar" {
		t.Errorf("expected 'bar', got '%v'", val)
	}

	// Test getting a non-existent value
	if _, exists := c.Get("baz"); exists {
		t.Errorf("expected 'baz' to not exist")
	}

	// Test value expiration
	c.Set("expire", "value", 1*time.Second)
	time.Sleep(2 * time.Second)
	if _, exists := c.Get("expire"); exists {
		t.Errorf("expected 'expire' to be expired")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := cache.New()
	key := "concurrent"
	value := "access"

	// Start multiple goroutines to set a value concurrently
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Set(key, value, 5*time.Second)
		}()
	}

	wg.Wait()

	// Verify that the value was set correctly
	if val, exists := c.Get(key); !exists || val != value {
		t.Errorf("expected '%s', got '%v'", value, val)
	}
}
