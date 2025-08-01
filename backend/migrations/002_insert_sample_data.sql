-- 示例数据插入
-- 基于客户属性设计文档的示例数据
-- MySQL版本

-- 先更新客户编码序列表
INSERT INTO customer_code_sequence (year, sequence_number) 
VALUES (YEAR(NOW()), 5) 
ON DUPLICATE KEY UPDATE sequence_number = 5;

-- 插入示例客户数据
INSERT INTO customers (
    customer_code, customer_name, customer_type, 
    contact_person, contact_title, email, phone, 
    address, company_size, preferred_license_type, 
    customer_level, status, description, created_by
) VALUES 
-- 示例1：VIP企业客户
(
    CONCAT('CUS-', YEAR(NOW()), '-0001'), 
    '北京科技有限公司', 
    'enterprise',
    '张三', 
    '技术总监', 
    'zhangsan@example.com', 
    '13800138000',
    '北京市朝阳区科技园区100号', 
    'medium', 
    'online',
    'vip', 
    'active', 
    '重要客户，IoT平台采购方，需要优先技术支持', 
    UUID()
),
-- 示例2：普通个人客户
(
    CONCAT('CUS-', YEAR(NOW()), '-0002'), 
    '李四', 
    'individual',
    '李四', 
    NULL, 
    'lisi@personal.com', 
    '13900139000',
    '上海市浦东新区张江高科技园区200号', 
    NULL, 
    'offline',
    'normal', 
    'active', 
    '个人开发者，主要用于学习和小型项目', 
    UUID()
),
-- 示例3：政府客户
(
    CONCAT('CUS-', YEAR(NOW()), '-0003'), 
    '某市政府信息化中心', 
    'government',
    '王五', 
    '信息化处处长', 
    'wangwu@gov.cn', 
    '010-12345678',
    '某市政府大楼信息化中心', 
    'large', 
    'hybrid',
    'enterprise', 
    'active', 
    '政府客户，智慧城市项目，需要专项技术支持', 
    UUID()
),
-- 示例4：教育客户
(
    CONCAT('CUS-', YEAR(NOW()), '-0004'), 
    '清华大学计算机学院', 
    'education',
    '赵六', 
    '实验室主任', 
    'zhaoliu@tsinghua.edu.cn', 
    '010-62785678',
    '北京市海淀区清华大学计算机学院', 
    'large', 
    'online',
    'strategic', 
    'active', 
    '重点高校合作伙伴，科研项目合作，享受教育优惠政策', 
    UUID()
),
-- 示例5：禁用状态的客户
(
    CONCAT('CUS-', YEAR(NOW()), '-0005'), 
    '测试公司', 
    'enterprise',
    '测试用户', 
    '测试经理', 
    'test@disabled.com', 
    '13700137000',
    '测试地址', 
    'small', 
    'online',
    'normal', 
    'disabled', 
    '测试客户，已禁用，用于演示禁用状态', 
    UUID()
);