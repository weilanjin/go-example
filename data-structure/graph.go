package main

// 图 graph 是一种非线性数据结构，它由顶点（vertex）和边（edge）组成。
// 图 G = {V, E}
// 一组顶点 V = {v1, v2, v3, ..., vn}
// 一组边 E = {(v1, v2), (v1, v3), (v2, v3), ..., en}
//
// 无向图 (undirected graph)：边没有方向，即 (v1, v2) 和 (v2, v1) 是相同的边。【微信的好友】
// 有向图 (directed graph)：边有方向，即 (v1, v2) 表示 v1 到 v2 的边，(v2, v1) 表示 v2 到 v1 的边。【微博关注与被关注】
//
// 根据所有顶点是否连通
// - 连通图 (connected graph)：所有顶点都是连通的，即任意两个顶点之间都有路径。
// - 非连通图 (disconnected graph)：存在顶点之间没有路径。
//
// 有权图 (weighted graph)：边有权值，即 (v1, v2) 的权值是 x，(v2, v1) 的权值是 y，x != y 【《王者荣耀》系统会根据共同游戏时间来计算玩家之间的“亲密度”】
//
// 常用术语
// 邻接 adjacency: 顶点 v1 和 v2 之间有边，则称 v1 和 v2 为邻接的顶点。
// 路径 path: 从顶点 v1 到顶点 v2 的一条边序列。
// 度 degree:  顶点 v 的度，即顶点 v 的邻接顶点的个数。
//
// 图的常见应用
// 社交网络： 用户（顶点）， 好友关系（边）， 潜在的好友关系（图计算问题）
// 地铁路线： 站点（顶点）， 站点间的连通性（边）， 最短路线推荐（图计算问题）
// 太阳系： 星体（顶点）， 星体间的万有引力作用（边）， 行星轨道计算（图计算问题）

type Graph struct {
	// 邻接表 k 顶点，v 该顶点的所有邻居顶点
	adjList map[int][]int
}

func NewGraph(edges [][]int) *Graph {
	g := &Graph{
		adjList: make(map[int][]int),
	}
	// 添加所有顶点和边
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		g.adjList[u] = append(g.adjList[u], v)
		g.adjList[v] = append(g.adjList[v], u)
	}
	return g
}

// 获取顶点数量
func (g *Graph) Size() int {
	return len(g.adjList)
}

func (g *Graph) addEdge(u, v int) {
	g.adjList[u] = append(g.adjList[u], v)
	g.adjList[v] = append(g.adjList[v], u)
}