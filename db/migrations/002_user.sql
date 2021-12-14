
-- +goose Up
CREATE TABLE IF NOT EXISTS `user` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `username` VARCHAR(64) NULL DEFAULT NULL,
        `password` VARCHAR(64) NULL DEFAULT NULL,
        `email` VARCHAR(64) NULL DEFAULT NULL,
        PRIMARY KEY (`id`)
    );