-- 创建账户表
CREATE TABLE IF NOT EXISTS accounts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    account_name VARCHAR(50) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL
);

-- 插入初始数据
INSERT INTO accounts (account_name, balance) VALUES ('张三', 1000.00);

select * from accounts;

-- 设置事务隔离级别为READ UNCOMMITTED（允许脏读）
SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

-- 事务A：修改余额但不提交
START TRANSACTION;
UPDATE accounts SET balance = 500.00 WHERE account_name = '张三';
-- 注意：这里不执行COMMIT

-- 在另一个会话中执行事务B（或在同一会话中快速切换）
-- 事务B：读取未提交的修改（脏读）
SELECT * FROM accounts WHERE account_name = '张三';
-- 此时会看到余额变为500.00，尽管事务A尚未提交

-- 回到事务A会话
ROLLBACK; -- 回滚事务A
-- 现在实际余额又变回1000.00

-- 在事务B会话中再次查询
SELECT * FROM accounts WHERE account_name = '张三';
-- 现在会看到余额恢复为1000.00，说明之前读到的是脏数据
