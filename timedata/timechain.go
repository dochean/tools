package timedata

import "time"

type TimeChain struct {
	head *timeNode
	length int
}

type timeNode struct {
	Value interface{}
	t time.Time
	next *timeNode
}

func NewTimeChain() *TimeChain{
	return &TimeChain{new(timeNode), 0}
}

func (t *TimeChain) Add(d time.Duration, v interface{}) {
	n := &timeNode{v, time.Now().Add(d), nil}
	add(t.head, n)
	t.length++
}

func (t *TimeChain) PeekTime() time.Time {
	if t.length == 0 {
		return time.Now()
	}
	return t.head.next.t
}

func (t *TimeChain) Len() int {
	return t.length
}

func (t *TimeChain) DeleteHead() interface{} {
	if t.length == 0 {
		return nil
	}
	n := deleteHead(t.head)
	t.length--
	return n.Value
}

func add(head, node *timeNode) {
	h, n := head, head.next
	for ; n != nil; h, n = n, n.next {
		if node.t.Before(n.t) {
			h.next, node.next = node, n
			break
		}
	}
	if node.next == nil {
		h.next = node
	}
}

func deleteHead(head *timeNode) *timeNode {
	if head.next == nil {
		return nil
	}
	n := head.next
	head.next = n.next
	return n
}