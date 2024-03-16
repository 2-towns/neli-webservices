-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `tribe` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `leader_id` BIGINT UNSIGNED NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_user_tribe_leader_id` FOREIGN KEY (leader_id) REFERENCES user(id) ON DELETE CASCADE,
  CONSTRAINT `fk_user_tribe_user_id` FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS `tribe`;
