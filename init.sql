-- Create database if not exists, database name passed from config
CREATE DATABASE IF NOT EXISTS `{db_name}`;
USE `{db_name}`;

-- Table users
-- Table checkins
-- Table tasks
-- Table points
-- Table moments
-- Table comments
-- Table likes
-- Table reward_logs

-- Table users (updated)
CREATE TABLE IF NOT EXISTS `users`
(
    `id`                BIGINT       NOT NULL,
    `created_at`        BIGINT       NOT NULL,
    `updated_at`        BIGINT       NOT NULL,
    `username`          VARCHAR(255) NOT NULL,
    `is_premium`        BOOL         NOT NULL DEFAULT false,
    `referral_code`     VARCHAR(255) NOT NULL,
    `language_code`     VARCHAR(255),
    `referral_uid`      BIGINT,
    `reward_points`     INT          NOT NULL DEFAULT 200,
    `last_points_refresh` TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX RCODE (referral_code),
    INDEX RUID (referral_uid),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table checkins
CREATE TABLE IF NOT EXISTS `checkins`
(
    `id`                    VARCHAR(255) NOT NULL,
    `uid`                   BIGINT       NOT NULL,
    `created_at`            BIGINT       NOT NULL,
    `updated_at`            BIGINT       NOT NULL,
    `continuous_days`       INT          NOT NULL DEFAULT 1,
    `reward_point`          BIGINT       NOT NULL DEFAULT 0,
    `status`                INT          NOT NULL DEFAULT 0,
    INDEX UID (uid),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table tasks
CREATE TABLE IF NOT EXISTS `tasks`
(
    `id`                    VARCHAR(255) NOT NULL,
    `uid`                   BIGINT       NOT NULL,
    `created_at`            BIGINT       NOT NULL,
    `updated_at`            BIGINT       NOT NULL,
    `task_group`            VARCHAR(255) NOT NULL,
    `task_type`             VARCHAR(255) NOT NULL,
    `status`                INT          NOT NULL DEFAULT 0,
    `reward_point`          BIGINT       NOT NULL DEFAULT 0,
    INDEX UID (uid),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table points
CREATE TABLE IF NOT EXISTS `points`
(
    `uid`                        BIGINT       NOT NULL,
    `created_at`                 BIGINT       NOT NULL,
    `updated_at`                 BIGINT       NOT NULL,
    `last_claimed_point_value`   BIGINT       NOT NULL DEFAULT 0,
    `total_wallet_point_value`   BIGINT       NOT NULL DEFAULT 0,
    `last_invite_point_level`    BIGINT       NOT NULL DEFAULT 0,
    `total_invite_point_value`   BIGINT       NOT NULL DEFAULT 0,
    `total_point_value`          BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table moments
CREATE TABLE IF NOT EXISTS `moments` (
    `id`             BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`        BIGINT       NOT NULL,
    `content`        TEXT         NOT NULL,
    `image_url`      VARCHAR(255),
    `likes_count`    INT          NOT NULL DEFAULT 0,
    `comments_count` INT          NOT NULL DEFAULT 0,
    `rewards_amount` INT          NOT NULL DEFAULT 0,
    `created_at`     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX USER_ID (user_id),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table comments
CREATE TABLE IF NOT EXISTS `comments` (
    `id`               BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`          BIGINT       NOT NULL,
    `moment_id`        BIGINT       NOT NULL,
    `content`          TEXT         NOT NULL,
    `reply_to_user_id` BIGINT,
    `created_at`       timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX USER_ID (user_id),
    INDEX MOMENT_ID (moment_id),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table likes
CREATE TABLE IF NOT EXISTS `likes` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`       BIGINT       NOT NULL,
    `moment_id`     BIGINT       NOT NULL,
    `created_at`    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `user_moment` (`user_id`, `moment_id`),
    INDEX MOMENT_ID (moment_id),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table reward_logs
CREATE TABLE IF NOT EXISTS `reward_logs` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `moment_id`     BIGINT       NOT NULL,
    `from_user_id`  BIGINT       NOT NULL,
    `to_user_id`    BIGINT       NOT NULL,
    `amount`        INT          NOT NULL,
    `descr`         VARCHAR(255) NULL,
    `created_at`    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_from_user` (`from_user_id`),
    INDEX `idx_to_user` (`to_user_id`),
    INDEX `idx_moment` (`moment_id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;