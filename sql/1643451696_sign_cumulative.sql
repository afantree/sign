-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `sign_cumulative`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `aid`        bigint          NOT NULL COMMENT '角色id',
    `num`        tinyint         NOT NULL DEFAULT 0 COMMENT '领取次数',
    `sign_month` int             NOT NULL COMMENT '签到月份,202101',
    `created_at` timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY sign_aid_num_index (`aid`, `sign_month`, `num`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='签到累计领取信息表';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS sign_cumulative;