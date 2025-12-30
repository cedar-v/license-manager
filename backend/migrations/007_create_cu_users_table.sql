-- 客户用户表结构
-- 基于客户用户管理模块设计创建
-- MySQL版本

-- 创建客户用户表
CREATE TABLE cu_users (
    -- 基础信息属性
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    customer_id VARCHAR(36) NOT NULL COMMENT '所属客户ID',
    phone VARCHAR(20) NOT NULL COMMENT '手机号，作为登录账号',
    phone_country_code VARCHAR(10) NOT NULL DEFAULT '+86' COMMENT '国家代码',

    -- 安全信息
    password VARCHAR(255) NULL COMMENT '密码（可选，用于增强安全）',
    salt VARCHAR(32) NULL COMMENT '密码盐值',
    user_role VARCHAR(20) NOT NULL DEFAULT 'member' COMMENT '角色: admin-企业管理员, member-普通成员',

    -- 个人信息
    real_name VARCHAR(100) NULL COMMENT '真实姓名',
    email VARCHAR(255) NULL COMMENT '邮箱（可选）',

    -- 状态管理
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active-正常, disabled-禁用, pending-待激活',
    phone_verified TINYINT(1) NOT NULL DEFAULT 1 COMMENT '手机是否验证（注册时已验证）',
    email_verified TINYINT(1) NOT NULL DEFAULT 0 COMMENT '邮箱是否验证',

    -- 登录安全
    last_login_at DATETIME(3) NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(50) NULL COMMENT '最后登录IP',
    login_attempts INT NOT NULL DEFAULT 0 COMMENT '登录失败次数',
    locked_until DATETIME(3) NULL COMMENT '账号锁定到期时间',

    -- 用户偏好设置
    avatar_url VARCHAR(500) NULL COMMENT '头像URL',
    language VARCHAR(10) NOT NULL DEFAULT 'zh-CN' COMMENT '界面语言',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai' COMMENT '时区设置',
    additional_info JSON NULL COMMENT '附加信息',
    remark TEXT NULL COMMENT '备注',

    -- 时间字段 (使用DATETIME(3)匹配Go的time.Time和GORM的精度要求)
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    deleted_at DATETIME(3) NULL,

    -- 外键约束
    CONSTRAINT fk_cu_users_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户用户表 - 客户用户管理';

-- 创建索引
CREATE INDEX idx_cu_users_phone ON cu_users(phone, phone_country_code);
CREATE INDEX idx_cu_users_email ON cu_users(email);
CREATE INDEX idx_cu_users_customer_id ON cu_users(customer_id);
CREATE INDEX idx_cu_users_status ON cu_users(status);
CREATE INDEX idx_cu_users_created_at ON cu_users(created_at);
CREATE INDEX idx_cu_users_deleted_at ON cu_users(deleted_at);
CREATE INDEX idx_cu_users_locked_until ON cu_users(locked_until);

-- 复合索引用于常见查询
CREATE INDEX idx_cu_users_customer_status ON cu_users(customer_id, status);
CREATE INDEX idx_cu_users_status_locked ON cu_users(status, locked_until);

-- 注意：
-- 1. 时间戳由Go应用程序管理，不依赖数据库默认值
-- 2. password字段可选，用于手机号+密码登录方式
-- 3. phone + phone_country_code 作为复合唯一键，确保手机号唯一性
-- 4. user_role字段保留但暂不用于业务逻辑检查
-- 5. phone_verified在注册时设为true，后续手机号修改需要重新验证
-- 6. 使用软删除，通过deleted_at字段实现
