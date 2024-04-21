CREATE DATABASE wg_zhihu_user;
USE wg_zhihu_user;

CREATE TABLE `user` (
                        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                        `username` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '用户名',
                        `avatar` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '头像',
                        `mobile` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '手机号',
                        `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                        PRIMARY KEY (`id`),
                        KEY `ix_update_time` (`update_time`),
                        UNIQUE KEY `uk_mobile` (`mobile`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户表';

INSERT INTO USER(username, avatar, mobile) VALUES ('张三', 'https://beyond-blog.oss-cn-beijing.aliyuncs.com/avatar/2021/01/01/1609488000.jpg', '13800138000');