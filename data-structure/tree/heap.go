package tree

// 堆 heap 是一种满足特定条件的完全二叉树
// - 小顶推 min heap : 任意节点的值 <= 其子节点的值
// - 大顶堆 max heap : 任意节点的值 >= 其子节点的值
// 完全二叉树
// - 最底层节点靠右填充，其他层的节点都被填满。
// - 二叉树的根节点被称为 “推顶”，将底层最靠右的节点称为 “堆底”。
// - 对于大顶堆（小顶推），堆顶元素（根节点）的值是最大（最小）的
//
// 「推通常用于实现优先队列，大顶堆相当于元素按从大到小的顺序出队的优先队列」

// Heap 通过实现 heap.Interface 来构建整数大顶堆
type Heap struct {
	data []int
}

func (h *Heap) Push(v any) {
	iv, _ := v.(int)
	h.data = append(h.data, iv)
}

func (h *Heap) Pop() any {
	if len(h.data) == 0 {
		return nil
	}
	last := h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	return last
}

func (h *Heap) Len() int {
	return len(h.data)
}

func (h *Heap) Less(i, j int) bool {
	return h.data[i] > h.data[j]
}

func (h *Heap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Top 获取堆顶元素
func (h *Heap) Top() any {
	return h.data[0]
}

// 获取左节点的索引
func (h *Heap) left(i int) int {
	return 2*i + 1
}

// 获取右节点的索引
func (h *Heap) right(i int) int {
	return 2*i + 2
}
