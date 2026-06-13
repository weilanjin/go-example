-- 创建后台操作日志表
CREATE TABLE `operate_log` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '日志主键',
  `operator` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '操作者名称',
  `operate_at` BIGINT NOT NULL COMMENT '操作时间，Unix时间戳',
  `ip_address` VARCHAR(64) DEFAULT '' COMMENT '操作者IP地址',
  `op_type` TINYINT NOT NULL DEFAULT 0 COMMENT '操作类型：1-创建，2-更新，3-删除',
  `module_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '模块名称',
  `target_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '目标ID',
  `content` TEXT COMMENT '操作内容，JSON格式',
  `remark` VARCHAR(255) DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`),
  KEY `idx_module_target` (`module_name`, `target_id`),
  KEY `idx_operate_at` (`operate_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='操作日志表';