-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `settings` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `max_duration` VARCHAR(255),
  `update_at` TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `settings`;
