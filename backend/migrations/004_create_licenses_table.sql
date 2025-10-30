-- 许可证表结构
-- 基于授权模块属性设计创建
-- MySQL版本

-- 创建许可证表
CREATE TABLE licenses (
    -- 基础信息属性
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    license_key VARCHAR(200) NOT NULL UNIQUE COMMENT '许可证密钥，设备特定',
    authorization_code_id VARCHAR(36) NOT NULL COMMENT '关联的授权码ID',
    customer_id VARCHAR(36) NOT NULL COMMENT '客户ID，冗余字段便于查询',
    
    -- 硬件绑定属性
    hardware_fingerprint VARCHAR(200) NOT NULL COMMENT '绑定的硬件指纹',
    device_info JSON COMMENT '设备信息，包含CPU、内存等',
    activation_ip VARCHAR(45) COMMENT '激活时的IP地址',
    
    -- 激活状态属性
    status VARCHAR(20) NOT NULL DEFAULT 'inactive' COMMENT '许可证状态: active激活/inactive未激活/revoked已撤销',
    activated_at DATETIME(3) NULL COMMENT '激活时间',
    last_heartbeat DATETIME(3) NULL COMMENT '最后心跳时间',
    last_online_ip VARCHAR(45) COMMENT '最后在线IP',
    
    -- 配置同步属性
    config_updated_at DATETIME(3) NULL COMMENT '配置更新时间',
    
    -- 使用统计属性
    usage_data JSON COMMENT '使用数据，软件上报的统计信息',
    
    -- 时间字段 (使用DATETIME(3)匹配Go的time.Time和GORM的精度要求)
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    
    -- 外键约束
    FOREIGN KEY (authorization_code_id) REFERENCES authorization_codes(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
    
    -- 唯一约束：同一授权码下的硬件指纹应该唯一
    UNIQUE KEY uk_licenses_auth_hardware (authorization_code_id, hardware_fingerprint)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='许可证表 - 硬件激活凭证';

-- 创建索引
CREATE INDEX idx_licenses_license_key ON licenses(license_key);
CREATE INDEX idx_licenses_authorization_code_id ON licenses(authorization_code_id);
CREATE INDEX idx_licenses_customer_id ON licenses(customer_id);
CREATE INDEX idx_licenses_hardware_fingerprint ON licenses(hardware_fingerprint);
CREATE INDEX idx_licenses_status ON licenses(status);
CREATE INDEX idx_licenses_activated_at ON licenses(activated_at);
CREATE INDEX idx_licenses_last_heartbeat ON licenses(last_heartbeat);
CREATE INDEX idx_licenses_config_updated_at ON licenses(config_updated_at);
CREATE INDEX idx_licenses_created_at ON licenses(created_at);

-- 复合索引用于常见查询
CREATE INDEX idx_licenses_customer_status ON licenses(customer_id, status);
CREATE INDEX idx_licenses_auth_status ON licenses(authorization_code_id, status);
CREATE INDEX idx_licenses_status_activated ON licenses(status, activated_at);
CREATE INDEX idx_licenses_heartbeat_status ON licenses(last_heartbeat, status);

-- 注意：时间戳由Go应用程序管理，不依赖数据库默认值
-- 许可证状态说明：
-- active激活: 设备已激活并可正常使用
-- inactive未激活: 许可证已生成但设备尚未激活
-- revoked已撤销: 由于硬件更换、设备丢失或违规使用而被撤销
-- 在线状态(is_online)是虚字段，通过last_heartbeat时间计算：
-- 在线: last_heartbeat在5分钟内
-- 离线: last_heartbeat超过5分钟但在24小时内
-- 异常: last_heartbeat超过24小时