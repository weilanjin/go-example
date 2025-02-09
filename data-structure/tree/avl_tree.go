package tree

type AVLTree[T any] struct {
	val         T   // 节点值
	height      int // 节点高度(指该节点到它的最远叶节点的距离（边的数量）)
	Left, Right *AVLTree[T]
}

func NewAVLTree[T any](val T) *AVLTree[T] {
	return &AVLTree[T]{
		val:    val,
		height: 1,
	}
}

func (t *AVLTree[T]) Height(node *AVLTree[T]) int {
	if node == nil {
		return 0
	}
	return node.height
}

func (t *AVLTree[T]) updateHeight(node *AVLTree[T]) {
	if node == nil {
		return
	}
	// 节点高度等于最高子树高度
	node.height = max(t.Height(node.Left), t.Height(node.Right)) + 1
}

// 平衡因子
func (t *AVLTree[T]) balanceFactor(node *AVLTree[T]) int {
	if node == nil {
		return 0
	}
	// 左子树高度减去右子树高度
	return t.Height(node.Left) - t.Height(node.Right)
}

// 右旋转
func (t *AVLTree[T]) rightRotate(node *AVLTree[T]) *AVLTree[T] {
	left := node.Left
	node.Left = left.Right
	left.Right = node
	// 更新节点高度
	t.updateHeight(node)
	t.updateHeight(left)
	return left
}

// 左旋转
func (t *AVLTree[T]) leftRotate(node *AVLTree[T]) *AVLTree[T] {
	right := node.Right
	node.Right = right.Left
	right.Left = node
	// 更新节点高度
	t.updateHeight(node)
	t.updateHeight(right)
	return right
}
