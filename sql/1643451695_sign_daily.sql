-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `sign_daily`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `aid`        bigint          NOT NULL COMMENT '角色id',
    `sign_at`    int             NOT NULL COMMENT '签到那天,20210101',
    `created_at` timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY sign_aid_time_index (`aid`, `sign_at`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='签到信息表';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS sign_daily;
