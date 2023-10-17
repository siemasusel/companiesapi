CREATE TABLE IF NOT EXISTS companies (
  `id` VARCHAR(36) PRIMARY KEY UNIQUE,
  `name` VARCHAR(15) UNIQUE NOT NULL,
  `description` VARCHAR(3000),
  `employees_count` INT NOT NULL,
  `registered` BOOLEAN NOT NULL,
  `type` VARCHAR(100) NOT NULL 
);

