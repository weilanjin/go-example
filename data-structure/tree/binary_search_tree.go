package tree

// 二叉搜索树
type BinarySearchTree struct {
	root *TreeNode[int]
}

func NewBinarySearchTree() *BinarySearchTree {
	return &BinarySearchTree{}
}

// 查找节点
func (t *BinarySearchTree) Search(num int) *TreeNode[int] {
	node := t.root
	for node != nil {
		if num < node.Val {
			node = node.Left
		} else if num > node.Val {
			node = node.Right
		} else {
			return node
		}
	}
	return nil
}

func (t *BinarySearchTree) Insert(num int) {
	node := NewTreeNode(num)
	cur := t.root
	// 若树为空，则初始化根节点
	if t.root == nil {
		t.root = node
		return
	}
	var pre *TreeNode[int]
	for cur != nil {
		if cur.Val == num {
			return
		}
		pre = cur
		if cur.Val < num {
			cur = cur.Right
		} else {
			cur = cur.Left
		}
	}
	if pre.Val < num {
		pre.Right = node
	} else {
		pre.Left = node
	}
}

func (t *BinarySearchTree) Remove(num int) {
	cur := t.root
	if cur == nil {
		return
	}
	var pre *TreeNode[int]
	for cur != nil {
		if cur.Val == num {
			break
		}
		pre = cur
		if cur.Val < num {
			cur = cur.Right // 待删除节点在右子树中
		} else {
			cur = cur.Left // 待删除节点在左子树中
		}
	}
	if cur == nil {
		return
	}

	// 子节点数为 0 或 1
	if cur.Left == nil && cur.Right == nil {
		var child *TreeNode[int]
		if cur.Left != nil { // 取出待删除节点的子节点
			child = cur.Left
		} else {
			child = cur.Right
		}
		if cur != t.root { // 删除节点 cur
			if pre.Left == cur {
				pre.Left = child
			} else {
				pre.Right = child
			}
		} else { // 若删除节点为根节点，则重新指定根节点
			t.root = child
		}
	} else { // 子节点数 > 1
		tmp := cur.Right
		for tmp.Left != nil {
			tmp = tmp.Left
		}
		t.Remove(tmp.Val) // 递归删除节点 tmp
		cur.Val = tmp.Val
	}
}
