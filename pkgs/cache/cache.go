package cache

import (
	"sync"
	"time"
)

type Cache struct {
	Data map[string]Data
	sync.RWMutex
}

type Data struct {
	value interface{}
	ttl   time.Time
}

func New() *Cache {
	return &Cache{
		Data: make(map[string]Data),
	}
}

func (c *Cache) Set(k string, v interface{}, exp time.Duration) {
	c.Lock()
	defer c.Unlock()

	ttlExpire := time.Now().Add(exp)
	c.Data[k] = Data{
		value: v,
		ttl:   ttlExpire,
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.Data[k]; !ok || time.Now().After(v.ttl) {
		delete(c.Data, k)
		return nil, false
	} else {
		return v.value, true
	}
}

func (c *Cache) Delete(k string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Data, k)
}

func (c *Cache) Cleanup() {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.Data {
		if time.Now().After(v.ttl) {
			delete(c.Data, k)
		}
	}
}
