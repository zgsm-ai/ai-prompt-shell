package cache

import (
	"container/list"
	"sync"
)

// LRUCache 实现简单的LRU缓存
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

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		size:   size,
		list:   list.New(),
		values: make(map[string]*list.Element),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.values[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*cacheEntry).value, true
	}
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.values[key]; ok {
		// 更新现有值
		elem.Value.(*cacheEntry).value = value
		c.list.MoveToFront(elem)
		return
	}

	// 添加新元素
	if c.list.Len() >= c.size {
		// 淘汰最久未使用的元素
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
