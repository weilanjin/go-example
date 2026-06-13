package pagination

// 传统的 LIMIT offset, size 进行分页查询，
// 会引发非常严重的深度分页（Deep Pagination）性能问题。
// 因为数据库需要扫描过 offset 条记录后才开始返回结果。
/*
	- 当 offset 到达百万级别时，此查询会极度缓慢甚至超时
	SELECT * FROM operate_log ORDER BY id ASC LIMIT 1000000, 1000;
*/
