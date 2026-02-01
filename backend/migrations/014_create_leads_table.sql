-- 企业线索表
CREATE TABLE leads (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    lead_no VARCHAR(50) NOT NULL COMMENT '线索编号，系统生成',
    company_name VARCHAR(200) NOT NULL COMMENT '企业/机构名称',
    contact_name VARCHAR(100) NOT NULL COMMENT '联系人',
    contact_phone VARCHAR(20) NOT NULL COMMENT '联系电话',
    contact_email VARCHAR(100) COMMENT '邮箱地址',
    requirement TEXT NOT NULL COMMENT '需求描述',
    extra_info TEXT COMMENT '补充信息',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待联系, contacted-已联系, converted-已成交, invalid-已失效',
    follow_up_date DATETIME COMMENT '跟进日期',
    follow_up_record TEXT COMMENT '跟进记录',
    internal_note TEXT COMMENT '内部备注',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    INDEX idx_leads_status (status),
    INDEX idx_leads_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业线索表';
