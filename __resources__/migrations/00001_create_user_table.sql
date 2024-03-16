-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `firstname` VARCHAR(255) NOT NULL,
  `lastname` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL UNIQUE,
  `role` VARCHAR(255) NOT NULL,
  `password` TEXT,
  `password_reset_token` TEXT,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

ALTER TABLE `user` ADD CONSTRAINT `user_email_unique` UNIQUE (`email`);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `user`;
