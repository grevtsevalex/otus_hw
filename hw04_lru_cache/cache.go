package hw04lrucache

type Key string

type CacheElement struct {
	MapKey Key         // ключ, по которому он лежит в словаре
	Value  interface{} // значение
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, hasItem := c.items[key]
	if hasItem {
		cacheEl := CacheElement{
			MapKey: key,
			Value:  value,
		}
		item.Value = cacheEl
		c.queue.MoveToFront(item)
		return true
	}

	item = c.queue.PushFront(CacheElement{
		MapKey: key,
		Value:  value,
	})
	c.items[key] = item

	if c.queue.Len() > c.capacity {
		lruItem := c.queue.Back()
		c.queue.Remove(lruItem)
		switch v := lruItem.Value.(type) {
		case CacheElement:
			delete(c.items, v.MapKey)
		default:
			return false
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, hasItem := c.items[key]
	if !hasItem {
		return nil, false
	}

	c.queue.MoveToFront(item)
	switch v := item.Value.(type) {
	case CacheElement:
		return v.Value, true
	default:
		return nil, false
	}
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem)
	c.queue.Reset()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
