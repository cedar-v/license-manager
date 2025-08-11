-- 客户表结构
-- 基于客户属性设计与Web显示方案创建
-- MySQL版本

-- 创建客户表
CREATE TABLE customers (
    -- 基础标识
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    customer_code VARCHAR(20) NOT NULL UNIQUE,
    customer_name VARCHAR(200) NOT NULL,
    customer_type VARCHAR(20) NOT NULL DEFAULT 'enterprise' COMMENT '客户类型: individual, enterprise, government, education',
    
    -- 联系信息
    contact_person VARCHAR(100) NOT NULL,
    contact_title VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    
    -- 地址信息
    address TEXT,
    
    -- 商业属性
    company_size VARCHAR(20) COMMENT '企业规模: small, medium, large, enterprise',
    
    -- 客户分级
    customer_level VARCHAR(20) NOT NULL DEFAULT 'normal' COMMENT '客户等级: normal, vip, enterprise, strategic',
    
    -- 状态管理
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active, disabled',
    description TEXT,
    
    -- 时间字段 (使用DATETIME(3)匹配Go的time.Time和GORM的精度要求)
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    updated_by VARCHAR(36),
    deleted_at DATETIME(3) NULL  -- 软删除字段，匹配gorm.DeletedAt
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户信息表';

-- MySQL的列注释已在CREATE TABLE语句中添加

-- 创建索引
CREATE INDEX idx_customers_customer_name ON customers(customer_name);
CREATE INDEX idx_customers_customer_type ON customers(customer_type);
CREATE INDEX idx_customers_customer_level ON customers(customer_level);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_customers_created_at ON customers(created_at);
CREATE INDEX idx_customers_created_by ON customers(created_by);
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);  -- 软删除索引，提高查询性能

-- 注意：时间戳由Go应用程序管理，不依赖数据库默认值
-- 这样可以确保跨数据库兼容性并符合失败快速原则

-- 创建客户编码序列表
CREATE TABLE customer_code_sequence (
    year INT PRIMARY KEY,
    sequence_number INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户编码序列表';

-- 插入当前年份的序列记录
INSERT INTO customer_code_sequence (year, sequence_number) 
VALUES (YEAR(NOW()), 0) 
ON DUPLICATE KEY UPDATE sequence_number = sequence_number;

-- 注意：MySQL版本简化了客户编码生成逻辑
-- 客户编码将在应用层生成，格式：CUS-YYYY-NNNN
-- 可以通过以下SQL获取下一个序列号：
-- UPDATE customer_code_sequence SET sequence_number = sequence_number + 1 WHERE year = YEAR(NOW());
-- SELECT sequence_number FROM customer_code_sequence WHERE year = YEAR(NOW());