create database wgZhihu_article;
use wgZhihu_article;

CREATE TABLE `article` (
       `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
       `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
       `content` text COLLATE utf8_unicode_ci NOT NULL COMMENT '内容',
       `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '封面',
       `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
       `author_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '作者ID',
       `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:待审核 1:审核不通过 2:可见',
        -- 在业务端冗余一个计数器字段，以减少对应服务（例如评论服务）的压力。
        -- canal组件监听`评论计数表`的变更，投递到kafka中，文章mq服务消费kafka数据，将计数更新到本表中
       `comment_num` int(11) NOT NULL DEFAULT '0' COMMENT '评论数',
       `like_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
       `collect_num` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
       `view_num` int(11) NOT NULL DEFAULT '0' COMMENT '浏览数',
       `share_num` int(11) NOT NULL DEFAULT '0' COMMENT '分享数',
       `tag_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '标签ID',
       `publish_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
       `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
       `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
       PRIMARY KEY (`id`),
       KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='文章表';