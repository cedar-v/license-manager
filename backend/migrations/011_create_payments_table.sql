-- 支付订单表
-- 支持多种支付方式和业务场景的通用支付表
-- MySQL版本

-- 创建支付订单表
CREATE TABLE payments (
    -- 基础信息
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    payment_no VARCHAR(64) UNIQUE NOT NULL COMMENT '支付单号',
    business_type VARCHAR(50) NOT NULL COMMENT '业务类型 (package_order)',
    business_id VARCHAR(36) NULL COMMENT '业务ID (关联cu_orders.id)',

    -- 用户信息
    customer_id VARCHAR(36) NOT NULL COMMENT '客户ID',
    cu_user_id VARCHAR(36) NOT NULL COMMENT '客户用户ID',

    -- 支付信息
    amount DECIMAL(10,2) NOT NULL COMMENT '支付金额',
    currency VARCHAR(3) DEFAULT 'CNY' COMMENT '货币类型',
    payment_method VARCHAR(20) DEFAULT 'alipay' COMMENT '支付方式',
    payment_provider VARCHAR(20) DEFAULT 'alipay' COMMENT '支付提供商',

    -- 状态信息
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '支付状态: pending-待支付, paid-已支付, cancelled-已取消, expired-已过期, failed-支付失败',
    trade_no VARCHAR(64) NULL COMMENT '第三方交易号',
    payment_time DATETIME(3) NULL COMMENT '支付完成时间',
    expire_time DATETIME(3) NOT NULL COMMENT '支付过期时间',

    -- 支付链接和回调
    payment_url TEXT NULL COMMENT '支付链接',
    notify_data JSON NULL COMMENT '回调数据',
    extra_data JSON NULL COMMENT '扩展数据',

    -- 时间信息
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付订单表 - 支持多种支付方式和业务场景';

-- 创建索引
CREATE INDEX idx_payments_payment_no ON payments(payment_no);
CREATE INDEX idx_payments_business_type_id ON payments(business_type, business_id);
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_cu_user_id ON payments(cu_user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- 外键约束
ALTER TABLE payments ADD CONSTRAINT fk_payments_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE;
ALTER TABLE payments ADD CONSTRAINT fk_payments_cu_user_id FOREIGN KEY (cu_user_id) REFERENCES cu_users(id) ON DELETE CASCADE;

-- 注意：
-- 1. 时间戳由Go应用程序管理，不依赖数据库默认值
-- 2. payment_no全局唯一，作为支付单的主要标识
-- 3. business_type + business_id 支持关联不同业务场景
-- 4. extra_data字段用于存储支付方式特定的扩展信息
-- 5. 支持软删除（如果需要可以通过deleted_at字段扩展）
