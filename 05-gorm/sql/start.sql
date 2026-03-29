CREATE DATABASE IF NOT EXISTS test_gorm_db DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 检查用 
use test_gorm_db;
SELECT * FROM umbrellas LIMIT 10;

# 一次性把整张表数据全部清空
# 特点：速度极快，会重置自增 ID 下次插入 ID 又从 1 开始
TRUNCATE TABLE umbrellas;

# 删表
DROP TABLE umbrellas;