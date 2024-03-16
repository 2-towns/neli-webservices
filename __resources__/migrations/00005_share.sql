-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `share` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `url` TEXT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `content_id` BIGINT UNSIGNED NOT NULL,
  `message` TEXT NOT NULL,
  `expiration_date` TIMESTAMP NULL,
  `created_at` TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `settings`;
