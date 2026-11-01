CREATE TABLE IF NOT EXISTS `hg_addon_users_tree` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT 'uid',
  `pid` int NOT NULL COMMENT 'pid',
  `level` int NOT NULL COMMENT 'level',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `users_tree_pid_level_index` (`pid`,`level`) USING BTREE,
  KEY `users_tree_uid_index` (`uid`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户关系';