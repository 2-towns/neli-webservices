-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user CHANGE firstname firstname VARCHAR(255) NULL;
ALTER TABLE user CHANGE lastname lastname VARCHAR(255) NULL;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user CHANGE firstname firstname VARCHAR(255) NOT NULL;
ALTER TABLE user CHANGE lastname lastname VARCHAR(255) NOT NULL;