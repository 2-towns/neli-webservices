-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `video_content` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255),
  `description` TEXT,
  `path` TEXT,
  `duration` INT UNSIGNED,
  `leader_id` BIGINT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP NULL,
  `update_at` TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
  CONSTRAINT `fk_user_video_content_id` FOREIGN KEY (leader_id) REFERENCES user(id) ON DELETE CASCADE,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `video_content`;
