package cache

type node[KeyT comparable, ValueT any] struct {
	key   KeyT
	value ValueT
	prev  *node[KeyT, ValueT]
	next  *node[KeyT, ValueT]
}

type list[KeyT comparable, ValueT any] struct {
	head *node[KeyT, ValueT]
	tail *node[KeyT, ValueT]
}

func newList[KeyT comparable, ValueT any]() *list[KeyT, ValueT] {
	var key KeyT
	var value ValueT
	list := &list[KeyT, ValueT]{
		head: &node[KeyT, ValueT]{key, value, nil, nil},
		tail: &node[KeyT, ValueT]{key, value, nil, nil},
	}
	list.head.next = list.tail
	list.tail.prev = list.head
	return list
}

func (l *list[KeyT, ValueT]) pushToFront(node *node[KeyT, ValueT]) {
	node.prev = l.head
	node.next = l.head.next
	l.head.next.prev = node
	l.head.next = node
}

func (l *list[KeyT, ValueT]) remove(node *node[KeyT, ValueT]) {
	if node == nil {
		return
	}
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
}

func (l *list[KeyT, ValueT]) moveToFront(node *node[KeyT, ValueT]) {
	if node == nil {
		return
	}
	l.remove(node)
	l.pushToFront(node)
}

func (l *list[KeyT, ValueT]) back() *node[KeyT, ValueT] {
	if l.tail == l.head {
		return nil
	}
	return l.tail.prev
}
