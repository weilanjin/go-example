CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    department VARCHAR(50),
    salary INT,
    hire_date DATE,
    INDEX idx_department_salary (department, salary),
    INDEX idx_last_name (last_name)
) ENGINE=InnoDB;