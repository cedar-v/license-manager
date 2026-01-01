-- 客户用户订单表结构
-- 基于产品套餐体系设计创建
-- MySQL版本

-- 创建客户用户订单表
CREATE TABLE cu_orders (
    -- 基础信息
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    order_no VARCHAR(50) NOT NULL COMMENT '订单号',
    customer_id VARCHAR(36) NOT NULL COMMENT '客户ID',
    cu_user_id VARCHAR(36) NOT NULL COMMENT '客户用户ID',
    package_id VARCHAR(50) NOT NULL COMMENT '套餐ID (对应配置中的ID)',
    package_name VARCHAR(100) NOT NULL COMMENT '套餐名称(快照)',

    -- 购买信息
    license_count INT NOT NULL COMMENT '许可数量',
    unit_price DECIMAL(10,2) NOT NULL COMMENT '单价(每许可价格，快照)',
    discount_rate DECIMAL(3,2) NOT NULL DEFAULT 1.0 COMMENT '折扣率',
    total_amount DECIMAL(10,2) NOT NULL COMMENT '订单总金额',

    -- 状态信息
    status VARCHAR(20) NOT NULL DEFAULT 'paid' COMMENT '订单状态: paid-已支付',
    authorization_code VARCHAR(50) NULL COMMENT '生成的授权码',

    -- 时间信息
    expired_at DATETIME(3) NULL COMMENT '订单过期时间',
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    deleted_at DATETIME(3) NULL,

    -- 外键约束
    CONSTRAINT fk_cu_orders_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    CONSTRAINT fk_cu_orders_cu_user_id FOREIGN KEY (cu_user_id) REFERENCES cu_users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户用户订单表 - 产品套餐购买订单';

-- 创建索引
CREATE INDEX idx_cu_orders_customer_id ON cu_orders(customer_id);
CREATE INDEX idx_cu_orders_cu_user_id ON cu_orders(cu_user_id);
CREATE INDEX idx_cu_orders_status ON cu_orders(status);
CREATE INDEX idx_cu_orders_created_at ON cu_orders(created_at);
CREATE INDEX idx_cu_orders_package_id ON cu_orders(package_id);
CREATE INDEX idx_cu_orders_deleted_at ON cu_orders(deleted_at);

-- 注意：
-- 1. 时间戳由Go应用程序管理，不依赖数据库默认值
-- 2. license_count决定授权码可以激活多少次
-- 3. authorization_code在订单支付后生成
-- 4. 使用软删除，通过deleted_at字段实现
