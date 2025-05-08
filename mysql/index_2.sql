-- 索引是否失效

-- 创建简化版的用户表
CREATE TABLE customer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(20),
    age INT,
    birth_date DATE,
    income DECIMAL(10,2),
    mobile VARCHAR(15),
    is_active TINYINT(1),
    gender CHAR(1),
    
    -- 创建各种索引
    INDEX idx_username (username),
    INDEX idx_age (age),
    INDEX idx_birth_date (birth_date),
    INDEX idx_income (income),
    INDEX idx_mobile (mobile),
    INDEX idx_is_active (is_active),
    INDEX idx_gender (gender),
    INDEX idx_username_age (username, age),
    INDEX idx_composite (username, is_active, age)
);

-- 插入20条测试数据
INSERT INTO customer (username, age, birth_date, income, mobile, is_active, gender) VALUES
('john_doe', 30, '1990-05-15', 5000.00, '13800138000', 1, 'M'),
('jane_smith', 25, '1995-08-22', 6500.00, '13900139000', 1, 'F'),
('bob_johnson', 40, '1980-11-30', 8000.00, '13700137000', 0, 'M'),
('alice_wang', 28, '1992-03-18', 7200.00, '13600136000', 1, 'F'),
('tom_lee', 35, '1985-07-10', 9000.00, '13500135000', 1, 'M'),
('sara_chen', 22, '1998-09-05', 6000.00, '13400134000', 0, 'F'),
('mike_zhang', 45, '1975-12-20', 12000.00, '13300133000', 1, 'M'),
('lily_liu', 29, '1991-04-25', 7500.00, '13200132000', 1, 'F'),
('david_wu', 33, '1987-06-15', 8500.00, '13100131000', 0, 'M'),
('emma_zhao', 27, '1993-02-28', 6800.00, '13000130000', 1, 'F'),
('alex_sun', 31, '1989-10-12', 9500.00, '15900159000', 1, 'M'),
('olivia_hu', 24, '1996-07-30', 6200.00, '15800158000', 0, 'F'),
('kevin_ma', 38, '1982-05-22', 11000.00, '15700157000', 1, 'M'),
('sophia_lin', 26, '1994-11-08', 7000.00, '15600156000', 1, 'F'),
('ryan_zhou', 42, '1978-09-17', 10500.00, '15500155000', 0, 'M'),
('zoey_gao', 23, '1997-08-14', 5800.00, '15400154000', 1, 'F'),
('peter_cai', 36, '1984-04-03', 9800.00, '15300153000', 1, 'M'),
('mia_deng', 21, '1999-01-19', 5500.00, '15200152000', 0, 'F'),
('jack_huang', 39, '1981-03-27', 11500.00, '15100151000', 1, 'M'),
('ava_luo', 20, '2000-12-10', 5200.00, '15000150000', 1, 'F');


select * from customer;

-- 1. 使用 or 导致无法走索引（可以使用索引（索引合并））
EXPLAIN SELECT * FROM customer WHERE username = 'john_doe' OR age = 25;

-- 2. 对索引字段使用函数计算
EXPLAIN SELECT * FROM customer WHERE YEAR(birth_date) = 1990;
EXPLAIN SELECT * FROM customer WHERE income + 100 = 5100;

EXPLAIN SELECT * FROM customer WHERE mobile = 13800138000;

EXPLAIN SELECT * FROM customer WHERE is_active != 1;

EXPLAIN SELECT * FROM customer WHERE username LIKE '%doe';

EXPLAIN SELECT * FROM customer WHERE is_active = 1;

EXPLAIN SELECT * FROM customer WHERE gender = 'M';

EXPLAIN SELECT * FROM customer WHERE username > 'john' AND is_active = 1;