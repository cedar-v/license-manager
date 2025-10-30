-- 授权变更历史表结构
-- 基于授权模块属性设计创建
-- MySQL版本

-- 创建授权变更历史表
CREATE TABLE authorization_changes (
    -- 基础属性
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    authorization_code_id VARCHAR(36) NOT NULL COMMENT '授权码ID',
    change_type VARCHAR(30) NOT NULL COMMENT '变更类型: renewal续费/upgrade升级/limit_change限制调整/feature_toggle功能开关/lock锁定/unlock解锁/other其他',
    old_config JSON COMMENT '变更前配置',
    new_config JSON COMMENT '变更后配置',
    operator_id VARCHAR(36) NOT NULL COMMENT '操作人ID',
    reason TEXT COMMENT '变更原因',
    
    -- 时间字段 (使用DATETIME(3)匹配Go的time.Time和GORM的精度要求)
    created_at DATETIME(3) NOT NULL COMMENT '记录创建时间',
    
    -- 外键约束 - 修改为级联删除
    FOREIGN KEY (authorization_code_id) REFERENCES authorization_codes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='授权变更历史表 - 记录授权码的所有变更操作';

-- 创建索引
CREATE INDEX idx_authorization_changes_authorization_code_id ON authorization_changes(authorization_code_id);
CREATE INDEX idx_authorization_changes_change_type ON authorization_changes(change_type);
CREATE INDEX idx_authorization_changes_operator_id ON authorization_changes(operator_id);
CREATE INDEX idx_authorization_changes_created_at ON authorization_changes(created_at);

-- 复合索引用于常见查询
CREATE INDEX idx_authorization_changes_auth_type ON authorization_changes(authorization_code_id, change_type);
CREATE INDEX idx_authorization_changes_auth_time ON authorization_changes(authorization_code_id, created_at);
CREATE INDEX idx_authorization_changes_operator_time ON authorization_changes(operator_id, created_at);

-- 注意：
-- 1. 该表只有created_at字段，没有updated_at和deleted_at，因为历史记录不应被修改或删除
-- 2. 变更类型包括：
--    - renewal: 续费操作
--    - upgrade: 升级操作  
--    - limit_change: 限制调整（如修改max_activations）
--    - feature_toggle: 功能开关（修改feature_config）
--    - lock: 锁定操作
--    - unlock: 解锁操作
--    - other: 其他类型的变更
-- 3. old_config和new_config存储变更前后的完整配置快照，便于审计和回滚