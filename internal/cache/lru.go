package cache

import (
	"container/list"
	"sync"
)

// LRUCache implements simple LRU cache
type LRUCache struct {
	size   int
	list   *list.List
	values map[string]*list.Element
	mu     sync.Mutex
}

type cacheEntry struct {
	key   string
	value interface{}
}

/**
 * Create new LRU cache instance
 * @param size Maximum number of entries in cache
 * @return New LRUCache instance
 */
func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		size:   size,
		list:   list.New(),
		values: make(map[string]*list.Element),
	}
}

/**
 * Get value from cache by key
 * @param c LRUCache instance
 * @param key Lookup key
 * @return Value and existence flag
 */
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.values[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*cacheEntry).value, true
	}
	return nil, false
}

/**
 * Add or update value in cache
 * @param c LRUCache instance
 * @param key Entry key
 * @param value Entry value
 * Will evict least recently used item if cache is full
 */
func (c *LRUCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.values[key]; ok {
		// Update existing value
		elem.Value.(*cacheEntry).value = value
		c.list.MoveToFront(elem)
		return
	}

	// Add new element
	if c.list.Len() >= c.size {
		// Evict least recently used element
		lastElem := c.list.Back()
		if lastElem != nil {
			delete(c.values, lastElem.Value.(*cacheEntry).key)
			c.list.Remove(lastElem)
		}
	}

	newEntry := &cacheEntry{key: key, value: value}
	elem := c.list.PushFront(newEntry)
	c.values[key] = elem
}
