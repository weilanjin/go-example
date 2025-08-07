-- 索引下推 测试

-- 查看ICP相关的状态变量
SHOW STATUS LIKE 'Handler_read%';

-- 不使用ICP的情况(MySQL 5.6之前)
SET optimizer_switch='index_condition_pushdown=off';
EXPLAIN SELECT * FROM employees 
WHERE department = 'Engineering' AND salary > 80000 AND first_name LIKE 'First1%';

-- 使用ICP的情况
SET optimizer_switch='index_condition_pushdown=on';
EXPLAIN SELECT * FROM employees 
WHERE department = 'Engineering' AND salary > 80000 AND first_name LIKE 'First1%';

-- 启用ICP并测试执行时间
SET optimizer_switch='index_condition_pushdown=on';
SELECT SQL_NO_CACHE * FROM employees 
WHERE department = 'Marketing' AND salary > 60000 AND first_name LIKE 'First2%';

-- 禁用ICP并测试执行时间
SET optimizer_switch='index_condition_pushdown=off';
SELECT SQL_NO_CACHE * FROM employees 
WHERE department = 'Marketing' AND salary > 60000 AND first_name LIKE 'First2%';