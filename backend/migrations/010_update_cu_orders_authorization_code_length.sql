-- 更新cu_orders表authorization_code字段长度
-- 支持新的HMAC签名授权码格式（约150-200字符）

-- MySQL
-- 修改authorization_code字段长度从50到500字符
ALTER TABLE cu_orders MODIFY COLUMN authorization_code varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- 注意事项：
-- 1. 此迁移会更改现有数据结构，请在生产环境执行前做好备份
-- 2. 500字符长度足以支持当前和未来的授权码格式
-- 3. 如果将来需要更长，仍然可以进一步扩展
