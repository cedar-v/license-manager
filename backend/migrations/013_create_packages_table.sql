-- 套餐表
-- 用于管理员配置产品套餐，用户端从数据库查询

CREATE TABLE packages (
    -- 基础信息
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(100) NOT NULL COMMENT '套餐名称',
    type VARCHAR(20) NOT NULL COMMENT '套餐类型: trial-试用版, basic-基础版, professional-专业版, custom-定制版',
    price DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '价格',
    price_description VARCHAR(100) DEFAULT '' COMMENT '价格描述，如"定制报价"',
    duration_description VARCHAR(200) DEFAULT '' COMMENT '期限描述，如"永久有效"、"当月25日到期"',
    description VARCHAR(500) DEFAULT '' COMMENT '套餐说明',
    features TEXT DEFAULT '' COMMENT '功能项，JSON格式: ["功能1", "功能2"]',
    
    -- 状态信息
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1-启用, 0-禁用',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序，数字越大越靠前',
    remark VARCHAR(500) DEFAULT '' COMMENT '备注',
    
    -- 时间信息
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    
    -- 索引
    INDEX idx_packages_type (type),
    INDEX idx_packages_status (status),
    INDEX idx_packages_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品套餐表';

-- 插入初始数据
INSERT INTO packages (id, name, type, price, price_description, duration_description, description, features, status, sort_order, created_at, updated_at) VALUES
("DEMO-TRIAL", '试用版', 'trial', 0, '￥0', '当月25日到期', '免费体验全部功能', '["全部功能", "1个许可"]', 1, 1, NOW(), NOW()),
("DEMO-BASIC", '基础版', 'basic', 300, '300元/许可', '永久有效', '小型企业最佳选择', '["基础功能", "批量许可购买"]', 1, 2, NOW(), NOW()),
("DEMO-PRO", '专业版', 'professional', 2000, '2000元/许可', '永久有效', '企业级完整解决方案', '["全部功能", "技术支持", "数据分析"]', 1, 3, NOW(), NOW()),
("DEMO-CUSTOM", '企业定制版', 'custom', 0, '定制报价', '按需定制', '按需定制解决方案', '["按需定制", "专属服务"]', 1, 4, NOW(), NOW());
