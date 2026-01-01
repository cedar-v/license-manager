-- 更新授权码表code字段长度以支持自包含配置
-- 使用HMAC签名的极简格式，预计长度80-200字符，预留足够空间

-- MySQL
-- 删除现有的UNIQUE索引（基于code字段的唯一索引）
-- 注意：如果索引不存在会报错但不影响迁移，可以忽略
ALTER TABLE authorization_codes DROP INDEX IF EXISTS idx_authorization_codes_code;

-- 修改字段长度为1000字符以支持更长的授权码
ALTER TABLE authorization_codes MODIFY COLUMN code varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL;

-- 创建前缀UNIQUE索引（只索引前255个字符，足够区分授权码）
-- 这样既保证唯一性，又避免了InnoDB的索引长度限制（3072字节）
CREATE UNIQUE INDEX idx_authorization_codes_code_prefix ON authorization_codes(code(255));

-- 注意事项：
-- 1. 此迁移会更改现有数据结构，请在生产环境执行前做好备份
-- 2. 前缀索引保证了唯一性，但只检查前255个字符的唯一性
-- 3. 对于LIC-{随机}-{Base64数据}格式，前255字符的碰撞概率极低
