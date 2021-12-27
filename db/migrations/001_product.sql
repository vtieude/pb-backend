-- +goose Up
ALTER DATABASE app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE IF NOT EXISTS `product` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(64) NOT NULL,
        `product_key` VARCHAR(64) NOT NULL UNIQUE,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
        `active` TINYINT(1) NOT NULL DEFAULT 1,
        PRIMARY KEY (`id`)
    );