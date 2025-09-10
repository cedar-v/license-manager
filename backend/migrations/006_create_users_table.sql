-- 用户表结构
-- 基于用户管理模块设计创建
-- MySQL版本

-- 创建用户表
CREATE TABLE users (
    -- 基础信息属性
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名，唯一',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱地址，唯一',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值',
    full_name VARCHAR(100) NOT NULL COMMENT '真实姓名',
    phone VARCHAR(20) COMMENT '手机号码',
    
    -- 角色权限属性
    role VARCHAR(20) NOT NULL DEFAULT 'viewer' COMMENT '用户角色: admin管理员/operator操作员/viewer查看者',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '账号状态: active活跃/disabled禁用/locked锁定',
    
    -- 登录相关属性
    last_login_at DATETIME(3) NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP',
    login_attempts INT NOT NULL DEFAULT 0 COMMENT '登录失败次数',
    locked_until DATETIME(3) NULL COMMENT '账号锁定到期时间',
    
    -- 时间字段 (使用DATETIME(3)匹配Go的time.Time和GORM的精度要求)
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表 - 系统用户管理';

-- 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_last_login_at ON users(last_login_at);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_locked_until ON users(locked_until);

-- 复合索引用于常见查询
CREATE INDEX idx_users_role_status ON users(role, status);
CREATE INDEX idx_users_status_locked ON users(status, locked_until);

-- 注意：
-- 1. 时间戳由Go应用程序管理，不依赖数据库默认值
-- 2. password_hash使用bcrypt算法，成本因子建议设置为12
-- 3. 账号锁定机制：连续5次登录失败锁定30分钟
-- 4. 角色权限通过中间件在应用层控制
-- 5. 默认管理员账号需要在应用初始化时创建