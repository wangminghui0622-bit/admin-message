-- admin-message (mcms-rpc) 初始化脚本
-- 由 admin_zero 的 mysqldump DDL 提取并改写为可重复执行的形式：
-- 建表使用 CREATE TABLE IF NOT EXISTS，种子数据使用 INSERT IGNORE。

CREATE TABLE IF NOT EXISTS `mcms_sms_providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The SMS provider name | 短信服务的提供商',
  `secret_id` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The secret ID | 密钥 ID',
  `secret_key` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The secret key | 密钥 Key',
  `region` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The service region | 服务器所在地区',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Is it the default provider | 是否为默认提供商',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='SMS Provider Configuration Table | 短信服务商配置表';

CREATE TABLE IF NOT EXISTS `mcms_email_providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The email provider name | 电子邮件服务的提供商',
  `auth_type` tinyint unsigned NOT NULL COMMENT 'The auth type, supported plain, CRAMMD5 | 鉴权类型, 支持 plain, CRAMMD5',
  `email_addr` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The email address | 邮箱地址',
  `password` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'The email''s password | 电子邮件的密码',
  `host_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The host name is the email service''s host address | 电子邮箱服务的服务器地址',
  `identify` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'The identify info, for CRAMMD5 | 身份信息, 支持 CRAMMD5',
  `secret` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'The secret, for CRAMMD5 | 邮箱密钥, 用于 CRAMMD5',
  `port` int unsigned DEFAULT NULL COMMENT 'The port of the host | 服务器端口',
  `tls` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Whether to use TLS | 是否采用 tls 加密',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Is it the default provider | 是否为默认提供商',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Email Provider Configuration Table | 邮件服务商配置表';

CREATE TABLE IF NOT EXISTS `mcms_sms_logs` (
  `id` char(36) COLLATE utf8mb4_bin NOT NULL COMMENT 'UUID',
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `phone_number` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The target phone number | 目标电话',
  `content` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The content | 发送的内容',
  `send_status` tinyint unsigned NOT NULL COMMENT 'The send status, 0 unknown 1 success 2 failed | 发送的状态, 0 未知， 1 成功， 2 失败',
  `provider` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The sms service provider | 短信服务提供商',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='SMS Log Table | 短信日志表';

CREATE TABLE IF NOT EXISTS `mcms_email_logs` (
  `id` char(36) COLLATE utf8mb4_bin NOT NULL COMMENT 'UUID',
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `target` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The target email address | 目标邮箱地址',
  `subject` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The subject | 发送的标题',
  `content` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The content | 发送的内容',
  `send_status` tinyint unsigned NOT NULL COMMENT 'The send status, 0 unknown 1 success 2 failed | 发送的状态, 0 未知， 1 成功， 2 失败',
  `provider` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The sms service provider | 短信服务提供商',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Email Log Table | 邮件日志表';

-- 种子 provider（INSERT IGNORE 保证重复执行不冲突；密钥为占位值，需在后台修改）
INSERT IGNORE INTO `mcms_email_providers`
  (`created_at`, `updated_at`, `name`, `auth_type`, `email_addr`, `password`, `host_name`, `identify`, `secret`, `port`, `tls`, `is_default`)
VALUES
  (NOW(), NOW(), 'tencent', 1, 'input your email address', 'input your password', 'smtp.qq.com', NULL, NULL, 465, 1, 1);

INSERT IGNORE INTO `mcms_sms_providers`
  (`created_at`, `updated_at`, `name`, `secret_id`, `secret_key`, `region`, `is_default`)
VALUES
  (NOW(), NOW(), 'tencent', 'input your secret ID', 'input your secret key', 'ap-nanjing', 1);

-- API 注册（id 4001-5000 区间；路径带 /sys-api 网关前缀）
INSERT IGNORE INTO `sys_apis`
  (`id`, `created_at`, `updated_at`, `path`, `description`, `api_group`, `service_name`, `method`, `is_required`)
VALUES
  (4001, NOW(), NOW(), '/sys-api/email/send', 'apiDesc.sendEmail', 'message_sender', 'Mcms', 'POST', 0),
  (4002, NOW(), NOW(), '/sys-api/email_log', 'apiDesc.getEmailLogById', 'email_log', 'Mcms', 'POST', 0),
  (4003, NOW(), NOW(), '/sys-api/email_log/create', 'apiDesc.createEmailLog', 'email_log', 'Mcms', 'POST', 0),
  (4004, NOW(), NOW(), '/sys-api/email_log/update', 'apiDesc.updateEmailLog', 'email_log', 'Mcms', 'POST', 0),
  (4005, NOW(), NOW(), '/sys-api/email_log/delete', 'apiDesc.deleteEmailLog', 'email_log', 'Mcms', 'POST', 0),
  (4006, NOW(), NOW(), '/sys-api/email_log/list', 'apiDesc.getEmailLogList', 'email_log', 'Mcms', 'POST', 0),
  (4007, NOW(), NOW(), '/sys-api/email_provider', 'apiDesc.getEmailProviderById', 'email_provider', 'Mcms', 'POST', 0),
  (4008, NOW(), NOW(), '/sys-api/email_provider/create', 'apiDesc.createEmailProvider', 'email_provider', 'Mcms', 'POST', 0),
  (4009, NOW(), NOW(), '/sys-api/email_provider/update', 'apiDesc.updateEmailProvider', 'email_provider', 'Mcms', 'POST', 0),
  (4010, NOW(), NOW(), '/sys-api/email_provider/delete', 'apiDesc.deleteEmailProvider', 'email_provider', 'Mcms', 'POST', 0),
  (4011, NOW(), NOW(), '/sys-api/email_provider/list', 'apiDesc.getEmailProviderList', 'email_provider', 'Mcms', 'POST', 0),
  (4012, NOW(), NOW(), '/sys-api/sms/send', 'apiDesc.sendSms', 'message_sender', 'Mcms', 'POST', 0),
  (4013, NOW(), NOW(), '/sys-api/sms_log', 'apiDesc.getSmsLogById', 'sms_log', 'Mcms', 'POST', 0),
  (4014, NOW(), NOW(), '/sys-api/sms_log/create', 'apiDesc.createSmsLog', 'sms_log', 'Mcms', 'POST', 0),
  (4015, NOW(), NOW(), '/sys-api/sms_log/update', 'apiDesc.updateSmsLog', 'sms_log', 'Mcms', 'POST', 0),
  (4016, NOW(), NOW(), '/sys-api/sms_log/delete', 'apiDesc.deleteSmsLog', 'sms_log', 'Mcms', 'POST', 0),
  (4017, NOW(), NOW(), '/sys-api/sms_log/list', 'apiDesc.getSmsLogList', 'sms_log', 'Mcms', 'POST', 0),
  (4018, NOW(), NOW(), '/sys-api/sms_provider', 'apiDesc.getSmsProviderById', 'sms_provider', 'Mcms', 'POST', 0),
  (4019, NOW(), NOW(), '/sys-api/sms_provider/create', 'apiDesc.createSmsProvider', 'sms_provider', 'Mcms', 'POST', 0),
  (4020, NOW(), NOW(), '/sys-api/sms_provider/update', 'apiDesc.updateSmsProvider', 'sms_provider', 'Mcms', 'POST', 0),
  (4021, NOW(), NOW(), '/sys-api/sms_provider/delete', 'apiDesc.deleteSmsProvider', 'sms_provider', 'Mcms', 'POST', 0),
  (4022, NOW(), NOW(), '/sys-api/sms_provider/list', 'apiDesc.getSmsProviderList', 'sms_provider', 'Mcms', 'POST', 0);

-- 超级管理员角色（001）授权上述 22 条路径
INSERT IGNORE INTO `casbin_rules` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
VALUES
  (9001, 'p', '001', '/sys-api/email/send', 'POST', '', '', ''),
  (9002, 'p', '001', '/sys-api/email_log', 'POST', '', '', ''),
  (9003, 'p', '001', '/sys-api/email_log/create', 'POST', '', '', ''),
  (9004, 'p', '001', '/sys-api/email_log/update', 'POST', '', '', ''),
  (9005, 'p', '001', '/sys-api/email_log/delete', 'POST', '', '', ''),
  (9006, 'p', '001', '/sys-api/email_log/list', 'POST', '', '', ''),
  (9007, 'p', '001', '/sys-api/email_provider', 'POST', '', '', ''),
  (9008, 'p', '001', '/sys-api/email_provider/create', 'POST', '', '', ''),
  (9009, 'p', '001', '/sys-api/email_provider/update', 'POST', '', '', ''),
  (9010, 'p', '001', '/sys-api/email_provider/delete', 'POST', '', '', ''),
  (9011, 'p', '001', '/sys-api/email_provider/list', 'POST', '', '', ''),
  (9012, 'p', '001', '/sys-api/sms/send', 'POST', '', '', ''),
  (9013, 'p', '001', '/sys-api/sms_log', 'POST', '', '', ''),
  (9014, 'p', '001', '/sys-api/sms_log/create', 'POST', '', '', ''),
  (9015, 'p', '001', '/sys-api/sms_log/update', 'POST', '', '', ''),
  (9016, 'p', '001', '/sys-api/sms_log/delete', 'POST', '', '', ''),
  (9017, 'p', '001', '/sys-api/sms_log/list', 'POST', '', '', ''),
  (9018, 'p', '001', '/sys-api/sms_provider', 'POST', '', '', ''),
  (9019, 'p', '001', '/sys-api/sms_provider/create', 'POST', '', '', ''),
  (9020, 'p', '001', '/sys-api/sms_provider/update', 'POST', '', '', ''),
  (9021, 'p', '001', '/sys-api/sms_provider/delete', 'POST', '', '', ''),
  (9022, 'p', '001', '/sys-api/sms_provider/list', 'POST', '', '', '');
-- menu registration: message center management
INSERT IGNORE INTO `sys_menus` VALUES
(220,'2026-10-15 06:50:11','2026-10-15 06:50:11',4,1,1,'/mcms_dir','MessageCenterManagement','','LAYOUT',0,'Mcms',NULL,'route.messageCenterManagement','clarity:email-line',0,0,0,0,'',0,0,0,20,'',1000000),
(221,'2026-10-15 06:50:11','2026-10-15 06:50:11',1,2,2,'/mcms_email_provider','EmailProviderManagement','','/mcms/emailProvider/index',0,'Mcms',NULL,'route.emailProviderManagement','clarity:email-line',0,0,0,0,'',0,0,0,20,'',220),
(222,'2026-10-15 06:50:11','2026-10-15 06:50:11',2,2,2,'/mcms_sms_provider','SmsProviderManagement','','/mcms/smsProvider/index',0,'Mcms',NULL,'route.smsProviderManagement','clarity:mobile-line',0,0,0,0,'',0,0,0,20,'',220);

-- role 1 is authorized for the message center menus (resolved by path)
INSERT IGNORE INTO `role_menus` (`role_id`,`menu_id`)
SELECT 1, `id` FROM `sys_menus` WHERE `path` IN ('/mcms_dir','/mcms_email_provider','/mcms_sms_provider');


