CREATE TABLE `tb_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `hash` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `filename` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `filesize` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `filepath` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `last_edit_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_hash` (`hash`),
  KEY `idx_status` (`status`),
  KEY `idx_filename` (`filename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 

-- CREATE TABLE `tb_user` (
--   `id` int(11) NOT NULL AUTO_INCREMENT,
--   `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
--   `password` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
--   `email` varchar(64) DEFAULT '' COMMENT '邮箱',
--   `phone` varchar(128) DEFAULT '' COMMENT '手机号',
--   `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
--   `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机号是否已验证',
--   `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
--   `last_edit_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
--   `profile` text COMMENT '用户属性',
--   `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
--   PRIMARY KEY (`id`),
--   UNIQUE KEY `idx_phone` (`phone`),
--   KEY `idx_status` (`status`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `tb_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_edit_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tb_user_token` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户登录token',
    PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- 创建用户文件表
CREATE TABLE `tb_user_file` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `filename` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `filesize` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `hash` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `last_edit_time` datetime DEFAULT CURRENT_TIMESTAMP 
          ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
  KEY `idx_user_id` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;