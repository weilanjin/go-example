package main

// 链表节点结构
type ListNode[T any] struct {
	Val  T            // 节点值
	Next *ListNode[T] // 指向下一个节点的指针
}

// 创建一个新链表
func NewListNode[T any](val T) *ListNode[T] {
	return &ListNode[T]{
		Val: val,
	}
}

// 在链表的节点 n 之后插入节点 p
func Insert[T any](n, p *ListNode[T]) {
	p.Next = n.Next
	n.Next = p
}

// 删除链表节点 n 之后的首个节点
func Remove[T any](n *ListNode[T]) {
	if n.Next == nil {
		return
	}
	n.Next = n.Next.Next
}

func (l *ListNode[T]) Access(index int) *ListNode[T] {
	if index < 0 {
		return nil
	}
	cur := l
	for range index {
		if cur.Next == nil { // 链表长度不足 index
			return nil
		}
		cur = cur.Next
	}
	return cur
}

/*
func (l *ListNode[T]) Index(target T) int {
	idx := 0
	for l != nil {
		if l.Val == target {
			return idx
		}
		l = l.Next
		idx++

	}
	return -1
}
*/
