// Package scheduler
// 协调工作的软件有: Zookeeper、Consul、Etcd
// Zookeeper 为 Java 生态群提供了丰富的分布式同步原语(通过Curator库)
// Consul 提供分布式同步原语这事件上不是很积极.
// Etcd 提供了非常好的分布式同步原语
//   - 分布式锁
//   - 分布式读写锁
//   - Leader 选举
package scheduler

/*
	Leader 选举

    Leader 选举常常被应用在主从架构的系统中.
	主从架构中的服务节点分为主(Leader、Master) 和 从(Follower、Slave、Worker) 两种角色.
	1主n从 = 一共n+1个节点.
	主节点常常执行写操作, 从节点常常执行读操作.
	如果读写都在主节点, 从节点只是提供备份功能,主从架构就变成主备架构.
*/