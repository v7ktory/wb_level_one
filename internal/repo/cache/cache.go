package cache

type CacheRepo[KeyT comparable, ValueT any] interface {
	Get(key KeyT) (ValueT, bool)
	Put(key KeyT, value ValueT)
}
type LRUCache[KeyT comparable, ValueT any] struct {
	capacity int
	cache    map[KeyT]*node[KeyT, ValueT]
	list     *list[KeyT, ValueT]
}

func NewLRUCache[KeyT comparable, ValueT any](capacity int) CacheRepo[KeyT, ValueT] {
	return &LRUCache[KeyT, ValueT]{
		capacity: capacity,
		cache:    make(map[KeyT]*node[KeyT, ValueT]),
		list:     newList[KeyT, ValueT](),
	}
}

func (lru *LRUCache[KeyT, ValueT]) Get(key KeyT) (ValueT, bool) {
	if node, found := lru.cache[key]; found {
		lru.list.moveToFront(node)
		return node.value, true
	}
	var value ValueT
	return value, false
}

func (lru *LRUCache[KeyT, ValueT]) Put(key KeyT, value ValueT) {
	if node, found := lru.cache[key]; found {
		lru.list.moveToFront(node)
		node.value = value
		return
	}
	if len(lru.cache) == lru.capacity {
		back := lru.list.back()
		if back != nil {
			lru.list.remove(back)
			delete(lru.cache, back.key)
		}
	}
	newNode := &node[KeyT, ValueT]{key, value, nil, nil}
	lru.list.pushToFront(newNode)
	lru.cache[key] = newNode
}
