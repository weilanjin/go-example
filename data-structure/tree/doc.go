package tree

// 二叉树常见术语
// - 根节点 root node：二叉树的最顶层节点， 没有父节点
// - 叶节点 leaf node：没有子节点的节点。
// - 边 edge： 节点与节点之间的连接线。 即节点引用 指针
// - 节点所在的层 level： 从顶至底递增，根节点 level 为 1
// - 节点的度 degree：节点的子节点个数
// - 二叉树的高度 height：从根节点到所有叶子节点的最长路径的节点数
// - 节点的深度 depth：从根节点到该节点的最长路径的节点数

// 常见二叉树类型
// 1.完美二叉树 perfect binary tree：所有层的节点都被完全填满。
// 2.完全二叉树 complete binary tree：除了最后一层，其他层的节点都被填满，并且最后一层的节点从左到右依次填满。
// 3.完满二叉树 full binary tree：所有节点都有两个子节点
// 4.平衡二叉树 balanced binary tree：任意节点的左右子树的高度差不超过1

// 二叉树的遍历
// BFS 广度优先搜索 breadth-first search
// - 层序遍历 level order traversal ：从顶部到底部逐层遍历二叉树，并在每一层按照从左到右的顺序访问节点。
// DFS 深度优先搜索 depth-first search
// - 先序遍历 preorder traversal ：先访问根节点，然后递归访问左子树，最后递归访问右子树。
// - 中序遍历 inorder traversal ：先递归访问左子树，然后访问根节点，最后递归访问右子树。
// - 后序遍历 postorder traversal ：先递归访问左子树，然后递归访问右子树，最后访问根节点。
