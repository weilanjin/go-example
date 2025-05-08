-- 索引选择追踪

-- 开启 trace 记录（仅当前 session 有效）
set session optimizer_trace="enabled=on",end_markers_in_json=on;
-- 1. 使用 or 导致无法走索引（可以使用索引（索引合并））
SELECT * FROM customer WHERE username = 'john_doe' OR age = 25;

-- 查看优化器的执行逻辑
select * from information_schema.OPTIMIZER_TRACE;