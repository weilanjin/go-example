-- 创建无主键表
CREATE TABLE test_no_pk (name VARCHAR(20));

-- 插入数据
INSERT INTO test_no_pk VALUES ('Alice'), ('Bob');

-- 使用 ROW_NUMBER() 作为行标识符
SELECT ROW_NUMBER() OVER() AS row_id, name FROM test_no_pk;