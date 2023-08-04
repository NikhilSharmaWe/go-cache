package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	// sync.RWMutex is more flexible than sync.Mutex and provides a read-write lock.
	// It allows multiple goroutines to acquire the lock
	// for reading (i.e., read operations) simultaneously, as long as no
	// goroutine is holding the lock for writing
	// (i.e., a write operation). This means that read operations can be performed concurrently,
	// but only one write operation can occur at a time.
	// Read this: https://blog.devgenius.io/how-to-speed-up-concurrent-go-routines-with-mutex-by-upto-50-51863bfbea8d

	data map[string][]byte
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	c.lock.Unlock()

	delete(c.data, string(key))

	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keyStr := string(key)

	val, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)
	}

	// fmt.Println("Value", val)

	log.Printf("GET %s = %s\n", string(key), string(val))

	return val, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(key)] = value

	if ttl > 0 { // this is the case when ttl is not set
		go func() {
			<-time.After(ttl)
			delete(c.data, string(key))
		}()
	}

	return nil
}
