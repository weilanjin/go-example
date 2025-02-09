package tree

import "container/list"

// TreeNode 二叉树节点 binary tree
// 是一种非线性数据结构， 代表 祖先和后代之间的派生关系。
// 每个节点都有两个引用
// - 分别指向 左子节点 left-child node
// - 分别指向 右子节点 right-child node
// - 该节点被称为 父节点 parent node
type TreeNode[T any] struct {
	Val   T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

func NewTreeNode[T any](val T) *TreeNode[T] {
	return &TreeNode[T]{
		Val: val,
	}
}

// BFS 广度优先搜索
// 层序遍历
func (t *TreeNode[T]) LevelOrder() []T {
	var res []T
	if t == nil {
		return res
	}

	// 初始化队列，加入根节点
	q := list.New()
	q.PushBack(t)
	for q.Len() > 0 {
		node := q.Remove(q.Front()).(*TreeNode[T])
		res = append(res, node.Val)

		if node.Left != nil {
			q.PushBack(node.Left) // 左子节点入队
		}
		if node.Right != nil {
			q.PushBack(node.Right) // 右子节点入队
		}
	}
	return res
}

// DFS 深度先搜索

// 前序遍历
func (t *TreeNode[T]) PreOrder() []T {
	var res []T
	if t == nil {
		return res
	}

	// 根、左、右
	res = append(res, t.Val)
	res = append(res, t.Left.PreOrder()...)
	res = append(res, t.Right.PreOrder()...)
	return res
}

// 中序遍历
func (t *TreeNode[T]) InOrder() []T {
	var res []T
	if t == nil {
		return res
	}

	// 左、根、右
	res = append(res, t.Left.InOrder()...)
	res = append(res, t.Val)
	res = append(res, t.Right.InOrder()...)
	return res
}

// 后序遍历
func (t *TreeNode[T]) PostOrder() []T {
	var res []T
	if t == nil {
		return res
	}

	// 左、右、根
	res = append(res, t.Left.PostOrder()...)
	res = append(res, t.Right.PostOrder()...)
	res = append(res, t.Val)
	return res
}
