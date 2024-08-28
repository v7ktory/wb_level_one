package cache

import (
	"context"
	"fmt"

	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/internal/repo/pgdb"
)

type Cache[KeyT comparable, ValueT any] interface {
	Get(key KeyT) (ValueT, bool)
	Put(key KeyT, value ValueT)
}
type LRUCache[KeyT comparable, ValueT any] struct {
	capacity int
	cache    map[KeyT]*node[KeyT, ValueT]
	list     *list[KeyT, ValueT]
}

func NewLRUCache[KeyT comparable, ValueT any](capacity int) Cache[KeyT, ValueT] {
	return &LRUCache[KeyT, ValueT]{
		capacity: capacity,
		cache:    make(map[KeyT]*node[KeyT, ValueT]),
		list:     newList[KeyT, ValueT](),
	}
}
func Warmup(ctx context.Context, pgRepo *pgdb.PgRepo, lru Cache[string, *entity.Order]) error {
	orders, err := pgRepo.GetLRUOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to warm up cache: %w", err)
	}
	for _, order := range orders {
		if order != nil {
			lru.Put(order.UID, order)
		}
	}

	return nil
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
