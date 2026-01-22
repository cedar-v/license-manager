-- 发票管理表
-- 支持订单开票申请、管理员处理、文件上传等完整发票管理流程
-- MySQL版本

-- 创建发票表
CREATE TABLE `invoices` (
  `id` varchar(36) NOT NULL,
  `invoice_no` varchar(50) NOT NULL COMMENT '发票申请号',
  `order_id` varchar(36) NOT NULL COMMENT '订单ID（唯一）',
  `order_no` varchar(50) NOT NULL COMMENT '订单号（冗余快照，便于检索）',
  `customer_id` varchar(36) NOT NULL COMMENT '客户ID',
  `cu_user_id` varchar(36) NOT NULL COMMENT '申请人（客户用户）ID',

  `amount` decimal(10,2) NOT NULL COMMENT '发票金额（取订单total_amount）',

  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT 'pending/issued/rejected',
  `invoice_type` varchar(20) NOT NULL COMMENT 'personal/enterprise/vat_special',
  `title` varchar(200) NOT NULL COMMENT '发票抬头',
  `taxpayer_id` varchar(50) DEFAULT NULL COMMENT '纳税人识别号',
  `content` varchar(200) NOT NULL COMMENT '开票内容',
  `receiver_email` varchar(255) NOT NULL COMMENT '收票邮箱',
  `remark` varchar(1000) DEFAULT NULL COMMENT '备注',

  `invoice_file_url` varchar(500) DEFAULT NULL COMMENT '发票文件URL（PDF）',
  `uploaded_at` datetime(3) DEFAULT NULL COMMENT '上传时间',
  `uploaded_by` varchar(36) DEFAULT NULL COMMENT '上传人（管理员）ID',
  `issued_at` datetime(3) DEFAULT NULL COMMENT '开票完成时间',

  `reject_reason` varchar(500) DEFAULT NULL COMMENT '驳回原因',
  `suggestion` varchar(500) DEFAULT NULL COMMENT '建议修改内容',
  `rejected_at` datetime(3) DEFAULT NULL COMMENT '驳回时间',
  `rejected_by` varchar(36) DEFAULT NULL COMMENT '驳回人（管理员）ID',

  `download_token` varchar(64) DEFAULT NULL COMMENT '下载token（可选）',

  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_invoices_invoice_no` (`invoice_no`),
  UNIQUE KEY `idx_invoices_order_id` (`order_id`),
  KEY `idx_invoices_customer_id` (`customer_id`),
  KEY `idx_invoices_cu_user_id` (`cu_user_id`),
  KEY `idx_invoices_status` (`status`),
  KEY `idx_invoices_created_at` (`created_at`),
  KEY `idx_invoices_deleted_at` (`deleted_at`),
  KEY `idx_invoices_order_no` (`order_no`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='发票管理表 - 支持订单开票申请和处理流程';

-- 外键约束
ALTER TABLE `invoices` ADD CONSTRAINT `fk_invoices_order_id` FOREIGN KEY (`order_id`) REFERENCES `cu_orders` (`id`) ON DELETE CASCADE;
ALTER TABLE `invoices` ADD CONSTRAINT `fk_invoices_customer_id` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE;
ALTER TABLE `invoices` ADD CONSTRAINT `fk_invoices_cu_user_id` FOREIGN KEY (`cu_user_id`) REFERENCES `cu_users` (`id`) ON DELETE CASCADE;
ALTER TABLE `invoices` ADD CONSTRAINT `fk_invoices_uploaded_by` FOREIGN KEY (`uploaded_by`) REFERENCES `users` (`id`) ON DELETE SET NULL;
ALTER TABLE `invoices` ADD CONSTRAINT `fk_invoices_rejected_by` FOREIGN KEY (`rejected_by`) REFERENCES `users` (`id`) ON DELETE SET NULL;

-- 注意：
-- 1. 时间戳由Go应用程序管理，不依赖数据库默认值
-- 2. invoice_no全局唯一，作为发票申请的主要标识
-- 3. order_id唯一，确保一个订单最多一张发票
-- 4. 使用软删除，通过deleted_at字段实现
-- 5. 外键约束确保数据一致性，管理员字段允许为NULL
