-- 授权码表结构
-- 基于授权模块属性设计创建
-- MySQL版本

-- 创建授权码表
CREATE TABLE authorization_codes (
    -- 基础信息属性
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '唯一授权码，软件身份凭证',
    customer_id VARCHAR(36) NOT NULL COMMENT '关联客户ID',
    created_by VARCHAR(36) NOT NULL COMMENT '创建人ID',
    software_id VARCHAR(50) COMMENT '目标软件产品ID',
    description TEXT COMMENT '授权码描述',
    
    -- 授权核心属性
    start_date DATETIME(3) NOT NULL COMMENT '授权开始时间',
    end_date DATETIME(3) NOT NULL COMMENT '授权结束时间',
    deployment_type VARCHAR(20) NOT NULL DEFAULT 'standalone' COMMENT '部署类型: standalone单机版/cloud云端版/hybrid混合版',
    encryption_type VARCHAR(20) DEFAULT 'standard' COMMENT '加密类型: standard标准加密/advanced高级加密',
    software_version VARCHAR(50) COMMENT '支持的软件版本范围',
    
    -- 激活控制属性
    max_activations INT NOT NULL DEFAULT 1 COMMENT '最大许可证数量',
    
    -- 状态管理属性
    is_locked BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否被锁定',
    lock_reason TEXT COMMENT '锁定原因，自定义文本输入',
    locked_at DATETIME(3) NULL COMMENT '锁定时间',
    locked_by VARCHAR(36) COMMENT '锁定操作人',
    
    -- 通用功能控制属性
    feature_config JSON COMMENT '功能配置，软件产品基础功能',
    usage_limits JSON COMMENT '使用量限制，如用户数、API调用次数等',
    custom_parameters JSON COMMENT '自定义参数，软件特定功能或用户定制功能',
    
    -- 时间字段 (使用DATETIME(3)匹配Go的time.Time和GORM的精度要求)
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    
    -- 外键约束
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='授权码表 - 业务配置容器';

-- 创建索引
CREATE INDEX idx_authorization_codes_code ON authorization_codes(code);
CREATE INDEX idx_authorization_codes_customer_id ON authorization_codes(customer_id);
CREATE INDEX idx_authorization_codes_software_id ON authorization_codes(software_id);
CREATE INDEX idx_authorization_codes_deployment_type ON authorization_codes(deployment_type);
CREATE INDEX idx_authorization_codes_start_date ON authorization_codes(start_date);
CREATE INDEX idx_authorization_codes_end_date ON authorization_codes(end_date);
CREATE INDEX idx_authorization_codes_created_by ON authorization_codes(created_by);
CREATE INDEX idx_authorization_codes_created_at ON authorization_codes(created_at);
CREATE INDEX idx_authorization_codes_is_locked ON authorization_codes(is_locked);
CREATE INDEX idx_authorization_codes_locked_at ON authorization_codes(locked_at);

-- 复合索引用于常见查询
CREATE INDEX idx_authorization_codes_customer_status ON authorization_codes(customer_id, start_date, end_date);
CREATE INDEX idx_authorization_codes_software_status ON authorization_codes(software_id, deployment_type);

-- 注意：时间戳由Go应用程序管理，不依赖数据库默认值
-- 状态字段(status)是虚字段，通过应用逻辑计算得出：
-- normal正常: 当前时间在start_date和end_date之间且is_locked=false
-- locked已锁定: is_locked=true
-- expired已过期: 当前时间超过end_date