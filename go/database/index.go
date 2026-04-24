package database

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

/*
	查询分为三种类型的磁盘操作
		1. Scan the whole(全部的) data set.(No idex is used)
		2. Point query: Query the index by a specific (指定) Key。
		3. Range query: Query the index by a range.(The index is sorted)

	B-Tree
		Balanced binary trees can be queried and updated in O(log n) time.
		and can be range-queried.
		B-Tree is a balanced n tree.

	LSM Tree
		Log-structured Merge Tree.
		How to query:
			1. An LSM-Tree contains(包含) multiple levels of data.
			2. Each level is sorted and split into multiple files.
			3. A point query starts at the top level, if the key is not found, the search continues(继续) to the next level.
			4. A range query merges the results from all levels. higher levels have more priority when merging.(合并时高级别优先级更高)
		How to update:
			1. When updating a key, the key is inserted into a file form the top level first.
			2. If the file size exceeds a threshold（超过阈值）, merge it with the next level.
			3. The file size threshold increases exponentially(指数) with each level.

	B-Tree: The Ideas
		Each node of a B-Tree contains multiple keys and multiple links to its children.
		所有键都用于决定下一个子节点。
			   [1,  4,  9]
			   /    |    \
		      v     v     v
		[1, 2, 3] [4, 6] [9, 11, 12]
		如果一个节点太大而无法容纳在一页上，则将其拆分为两个节点。
		如果一个节点太小，尝试将其与兄弟节点合并。

	B-Tree: Operations
		B+ Tree 仅将值存储在叶节点中，内部节点仅包含键。

		将叶子节点拆分为 2 个节点后，指向父节点，如果大小增加，这可能触发进一步的拆分。
		 parent            parent
		 / ｜ \      =>   / ｜  ｜ \
		L1 L2 L6         L1 L3 L4 L6
		After the root node is split, a new root node is added, This is how a B-tree grows.
						  new_root
							/  \
		  root				N1 N2
		 / ｜ \      =>	  / ｜  ｜ \
		L1 L2 L6		 L1 L3 L4 L6

	Immutable (不可变的) Data Structures
		append-only (只追加)
		copy-on-write (写时复制)
		persistent data structures (持久化数据结构)

		在向叶节点插入key时，没有就地修改节点，而是创建一个新节点（包含待更新节点的所有keys和new key），父节点也必须指向新的节点。
		优点：
			1. 避免数据损坏，即使更新中断，旧版本的数据仍然保持不变。
			2. Easy concurrency, Readers can operate concurrently with writes.
*/

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	if node1max > BTREE_PAGE_SIZE {
		panic("node1max > BTREE_PAGE_SIZE")
	}
}

const HEADER = 4

// B-Tree Node
// A node consists of:
// 1. A fixed-sized header containing the type of the node (leaf node or internal node) and the number of keys.
// 2. A list of pointers to the child nodes(used by internal nodes).
// 3. A list of offsets pointing to each key-value pair.
// 4. Packed(包装) KV pairs.
//
// leaf nodes and internal nodes use the same format.
/*
	第二行 字段大小 bytes 为单位
	| type | nkeys |  pointers  |  offsets   | key-values | unused |
	|  2B  |   2B  | nkeys × 8B | nkeys × 2B |     ...    |        |

	header: 4 Byte
		[type] is the node type : leaf or internal
		[nkeys] is the number of keys in the node(也是 child pointers number)

	KV pair: 键值对都以其大小为前缀，对于内部节点，值的大小为0
	| key_size | val_size | key | val |
	|    2B    |    2B    | ... | ... |

	a leaf node {"k1": "hi", "k3", "hello"} is encoded as:
	| type | nkeys | pointers | offsets |            key-values           | unused |
	|   2  |   2   | nil nil  |  8 19   | 2 2 "k1" "hi"  2 5 "k3" "hello" |        |
	|  2B  |  2B   |   2×8B   |  2×4B   | 4B + 2B + 2B + 4B + 2B + 5B     |        |

	offsets: 第一个键值对的偏移始终为 0
	  8 是第二个键值对的偏移量
	  19 是第二个键值对末尾之后的偏移量
*/
type BNode []byte // can be dumped to the disk

// header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(btype, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// 读取和写入指针数组（针对内部节点）

func (node BNode) getPtr(idx uint16) uint64 {
	if idx > node.nkeys() {
		panic("invalid index")
	}
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	if idx > node.nkeys() {
		panic("invalid index")
	}
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// offsets list
// - 偏移量相对于第一对 KV 的位置
// - 第一对KV的偏移量始终为 0，因此不会存储在列表中
// - 将偏移量存储到偏移量列表中最后一个KV对的末尾，该列表用于确定节点的大小
func offsetPos(node BNode, idx uint16) uint16 {
	if idx < 1 || node.nkeys() < idx {
		panic("invalid index")
	}
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := offsetPos(node, idx)
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) setOffset(idx uint16, val uint16) {
	pos := offsetPos(node, idx)
	binary.LittleEndian.PutUint16(node[pos:], val)
}

// key-values
// 偏移量列表用于快速定位第n个KV对

func (node BNode) kvPos(idx uint16) uint16 {
	if idx <= node.nkeys() {
		panic("invalid index")
	}
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	if idx <= node.nkeys() {
		panic("invalid index")
	}
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	if idx <= node.nkeys() {
		panic("invalid index")
	}
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos+0:])
	vlen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+klen:][:vlen]
}

// node size in bytes
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys()) // uses the offset value of the last key
}

// nodeLookupLE
// find the last position of that is less than or equal to the key
// returns the first kid(小孩) node whose range intersects the key. (kid[i] <= key)
func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	i := uint16(0)
	for ; i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp == 0 {
			return i
		}
		if cmp > 0 {
			return i - 1
		}
	}
	return i - 1
}

// add a new key to a leaf node
func leafInsert(new, old BNode, idx uint16, key, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)                   // copy the keys before 'idx'
	nodeAppendKV(new, idx, 0, key, val)                    // the new key
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx) // keys form 'idx'
}

// if the key already exists, update it
func leafUpdate(new, old BNode, idx uint16, key, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys())
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-idx-1)
}

// copy multiple KVs into the position
func nodeAppendRange(new, old BNode, dstNew, srcOld, n uint16) {
	for i := uint16(0); i < n; i++ {
		dst, src := dstNew+i, srcOld+i
		nodeAppendKV(new, dst, old.getPtr(src), old.getKey(src), old.getVal(src))
	}
}

// copy a KV into the position
//
// example:
// node := make(BNode, BTREE_PAGE_SIZE)
// node.setHeader(BNODE_LEAF, 2)
// nodeAppendKV(node, 0, 0, []byte("k1"), []byte("hi"))
// nodeAppendKV(node, 1, 0, []byte("k3"), []byte("hello"))
//
// params:
// [idx] element index
// [ptr] 是第 n 个子指针，叶子节点未使用该指针
func nodeAppendKV(new BNode, idx uint16, ptr uint64, key, val []byte) {
	// ptrs
	new.setPtr(idx, ptr)
	// 4-bytes KVs size
	pos := new.kvPos(idx) // uses the offset of the previous key
	binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))
	// KV data
	copy(new[pos+4:], key)
	copy(new[pos+4+uint16(len(key)):], val)
	// update the offset of the next key
	new.setOffset(idx+1, new.getOffset(idx)+uint16(len(key)+len(val)+4))
}

// KV size limit

const (
	BTREE_PAGE_SIZE    = 4096 // 4KB 字节（典型的操作系统page size） 更大的 page 8k 16k
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

// 将键插入节点会增加其大小，导致其超过页面大小， 节点被分成多个较小的节点。

// nodeSplit2
// 大于允许的节点被分割成两个节点。
func nodeSplit2(left, right, old BNode) {
	if old.nkeys() < 2 {
		return
	}
	nleft := old.nkeys() / 2 // the initial guess 最初的估计
	// 尝试去容纳左半部分
	left_bytes := func() uint16 {
		return HEADER + 8*nleft + 2*nleft + old.getOffset(nleft)
	}
	for left_bytes() > BTREE_PAGE_SIZE {
		nleft--
	}
	if nleft < 1 {
		return
	}
	// 尝试去容纳右半部分
	right_bytes := func() uint16 {
		return old.nbytes() - left_bytes() + HEADER
	}
	for right_bytes() > BTREE_PAGE_SIZE {
		nleft++
	}
	if nleft > old.nkeys() {
		return
	}
	nright := old.nkeys() - nleft

	// new nodes
	left.setHeader(old.btype(), nleft)
	nodeAppendRange(left, old, 0, 0, nleft)

	right.setHeader(old.btype(), nright)
	nodeAppendRange(right, old, 0, nleft, nright)
}

// 胖节点被分成3个节点（中间有一个大KV对）
// split a node if it's too big. the results are 1-3 nodes.
func nodeSplit3(old BNode) (uint16, [3]BNode) {
	if old.nbytes() < BTREE_PAGE_SIZE {
		old = old[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old} // not split
	}
	left := make(BNode, 2*BTREE_PAGE_SIZE) // 稍后可能被拆分
	right := make(BNode, BTREE_PAGE_SIZE)
	nodeSplit2(left, right, old)
	if left.nbytes() <= BTREE_PAGE_SIZE {
		left = left[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}
	// the left node is still(仍然) too large
	leftleft := make(BNode, BTREE_PAGE_SIZE)
	middle := make(BNode, BTREE_PAGE_SIZE)
	nodeSplit2(leftleft, middle, left)
	return 3, [3]BNode{leftleft, middle, right}
}

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

// 磁盘 B+ 树 依赖磁盘 I/O 来解引用指针（page number）

type BTree struct {
	// 不能使用内存指针，该指针是引用磁盘页面而不是内存节点
	root uint64 // pointer (a nonzero page number)
	// callbacks for managing on-disk pages
	get func(uint64) BNode // dereference(间接引用) a pointer
	new func(BNode) uint64 // allocate a new page
	del func(uint64)       // deallocate(释放) a page
}

// insert a KV into a node, the result might(可能) be split into 2 nodes
// the caller（调用者） is responsible（负责） for deallocating(释放) the input node
// and splitting and allocating result nodes
func treeInsert(tree *BTree, node BNode, key, val []byte) BNode {
	// the result node.
	// it's allowed to be bigger than (大于) 1 page and will be split if so
	new := make(BNode, 2*BTREE_PAGE_SIZE)

	// where to insert the key ?
	idx := nodeLookupLE(node, key)

	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			// found the key, update it
			leafUpdate(new, node, idx, key, val)
		} else {
			// insert it after the position.
			leafInsert(new, node, idx+1, key, val)
		}
	case BNODE_NODE:
		// recursive(递归) insertion to the kid node.
		kptr := node.getPtr(idx)
		knode := treeInsert(tree, tree.get(kptr), key, val)
		// after insertion, split the result
		nsplit, split := nodeSplit3(knode)
		// deallocate the old kid node
		tree.del(kptr)
		// updata the kid links
		nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
	default:
		panic("bad node!")
	}
	return new
}

func nodeReplaceKidN(tree *BTree, new, old BNode, idx uint16, kids ...BNode) {
	// replace the kid pointers
	for i, kid := range kids {
		ptr := tree.new(kid)
		old.setPtr(idx+uint16(i), ptr)
	}
	// replace the offsets
	for i := uint16(1); i <= uint16(len(kids)); i++ {
	}
}

// Insert insert a new key or update an existing key
// - Create the root node if the tree is empty
// - Add a new root if the root node is split -- 如果root节点分裂，则创建新的root节点
func (tree *BTree) Insert(key, val []byte) error {
	// if err := checkLimit(key, val); err != nil {
	// 	return err
	// }

	if tree.root == 0 { // create the first node
		root := make(BNode, BTREE_PAGE_SIZE)
		root.setHeader(BNODE_LEAF, 2)

		// 一个虚拟键，这使得树覆盖了整个键空间。因此查找总是可以找到一个包含节点。
		nodeAppendKV(root, 0, 0, nil, nil)
		nodeAppendKV(root, 1, 0, key, val)
		tree.root = tree.new(root)
		return nil
	}
	node := treeInsert(tree, tree.get(tree.root), key, val)
	nsplit, split := nodeSplit3(node) // 如果 root 节点分裂，则扩容它
	tree.del(tree.root)
	if nsplit > 1 {
		root := make(BNode, BTREE_PAGE_SIZE)
		root.setHeader(BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.new(knode), knode.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(split[0])
	}
	return nil
}

// TODO
func (tree *BTree) Delete(key []byte) (bool, error) { return false, nil }

// 更新后的子节点应该和兄弟姐妹节点合并吗？
func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode) {
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}
	if idx > 0 {
		sibling := BNode(tree.get(node.getPtr(idx - 1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return 1, sibling // left
		}
	}

	if idx+1 < node.nbytes() {
		sibling := BNode(tree.get(node.getPtr(idx + 1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return +1, sibling // right
		}
	}
	return 0, BNode{}
}

type C struct {
	tree  BTree
	ref   map[string]string // the reference data
	pages map[uint64]BNode  // in-memory pages
}

func newC() *C {
	pages := map[uint64]BNode{}
	return &C{
		tree: BTree{
			get: func(ptr uint64) BNode {
				node, _ := pages[ptr]
				return node
			},
			new: func(b BNode) uint64 {
				ptr := uint64(uintptr(unsafe.Pointer(&b[0])))
				pages[ptr] = b
				return ptr
			},
			del: func(ptr uint64) {
				delete(pages, ptr)
			},
		},
		ref:   map[string]string{},
		pages: pages,
	}
}

func (c *C) add(key string, val string) {
	c.tree.Insert([]byte(key), []byte(val))
	c.ref[key] = val // reference data
}