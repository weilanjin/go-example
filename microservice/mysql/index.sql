-- Create a test table with proper indexes
CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY,
    `name` VARCHAR(100),
    email VARCHAR(100),
    age INT,
    city VARCHAR(50),
    created_at DATETIME
);

-- Add single-column indexes
CREATE INDEX idx_name ON users(`name`);
CREATE INDEX idx_age ON users(age);

-- Add composite index for multiple columns
CREATE INDEX idx_city_age_email ON users(city, age, email);


show index from users;

-- Insert some sample data
INSERT INTO users VALUES
(1, 'John Doe', 'john@example.com', 25, 'New York', '2023-01-01'),
(2, 'Jane Smith', 'jane@example.com', 30, 'Boston', '2023-01-02'),
(3, 'Bob Johnson', 'bob@example.com', 35, 'Chicago', '2023-01-03'),
(4, 'Alice Williams', 'alice@example.com', 28, 'New York', '2023-01-04'),
(5, 'Charlie Brown', 'charlie@example.com', 40, 'Boston', '2023-01-05');

-- 4. Test LIKE '%xxx%' pattern
-- This won't use the index on 'name'
EXPLAIN SELECT * FROM users WHERE name LIKE '%oh%';

-- But LIKE 'xxx%' (prefix match) can use the index
EXPLAIN SELECT * FROM users WHERE name LIKE 'Jo%';

-- 5. Test composite index without leftmost column
-- This won't use the composite index because 'city' is not in the WHERE clause
EXPLAIN SELECT * FROM users WHERE age = 30 AND email = 'jane@example.com';

-- This will use the composite index (follows leftmost principle)
EXPLAIN SELECT * FROM users WHERE city = 'Boston' AND age = 30;

-- This is also fine (uses just the first part of the index)
EXPLAIN SELECT * FROM users WHERE city = 'New York';

-- 6. Test type mismatch
-- This will not use the index efficiently due to type conversion
EXPLAIN SELECT * FROM users WHERE id = '3';

-- For comparison, this will use the index properly
EXPLAIN SELECT * FROM users WHERE id = 3;

-- Similarly with the age index
EXPLAIN SELECT * FROM users WHERE age = '30'; -- Type conversion, may not use index efficiently
EXPLAIN SELECT * FROM users WHERE age = 30;   -- Correct type, should use index

-- 强制使用复合索引
EXPLAIN SELECT * 
FROM users FORCE INDEX(idx_city_age_email) 
WHERE city = 'Boston' 
AND age = 30;

-- 提示使用复合索引
EXPLAIN SELECT * FROM users USE INDEX(idx_city_age_email)
WHERE city = 'Boston' AND age = 30;