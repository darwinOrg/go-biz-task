create table biz_task
(
    `id`                 bigint auto_increment primary key,
    `type`               int     NOT NULL COMMENT '任务类型',
    `channel`            int     NOT NULL COMMENT '渠道',
    `content`            longtext COMMENT '任务内容',
    `status`             tinyint NOT NULL DEFAULT 0 COMMENT '状态（0：初始化，1：处理中，2：成功，3：失败）',
    `reason`             varchar(64) COMMENT '失败原因',
    `scheduled_start_at` datetime(3) COMMENT '计划开始时间',
    `scheduled_end_at`   datetime(3) COMMENT '计划结束时间',
    `end_at`             datetime(3) COMMENT '执行结束时间',
    `priority`           int     NOT NULL default 100 COMMENT '优先级',
    `processed_count`    int     NOT NULL default 0 COMMENT '已处理次数',
    `locked_at`          datetime(3) COMMENT '锁定时间',
    `lock_until`         datetime(3) COMMENT '锁定到',
    `locked_by`          varchar(64) COMMENT '锁定者',
    `created_at`         datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3),
    `modified_at`        datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3)
) engine = innodb
  character set = utf8mb4 comment '业务任务';

CREATE TABLE `biz_task_result`
(
    `id`         bigint   NOT NULL AUTO_INCREMENT,
    `task_id`    bigint   NOT NULL COMMENT '任务id',
    `content`    longtext NOT NULL COMMENT '结果内容',
    `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_task_id` (`task_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='业务任务结果';