
-- +goose Up
CREATE TABLE IF NOT EXISTS `user` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `username` VARCHAR(256) Not NULL,
        `password` VARCHAR(256) Not NULL,
        `email` VARCHAR(64) NOT NULL UNIQUE,
        `active`  TINYINT(1) NOT NULL DEFAULT 1,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
        PRIMARY KEY (`id`)
    );
    CREATE TABLE IF NOT EXISTS `role` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `role_name` VARCHAR(64) Not NULL,
        `label` VARCHAR(256) NOT NULL,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP  ON UPDATE CURRENT_TIMESTAMP,
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
        PRIMARY KEY (`id`)
    );
CREATE TABLE IF NOT EXISTS `user_role` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `fk_user` int(10) Not NULL,
        `fk_role` int(10) Not NULL,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP  ON UPDATE CURRENT_TIMESTAMP,
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
         `active`  TINYINT(1) NOT NULL DEFAULT 1,
        PRIMARY KEY (`id`)
    ) ;
ALTER TABLE user_role ADD FOREIGN KEY (fk_user) REFERENCES user (id);
ALTER TABLE user_role ADD FOREIGN KEY (fk_role) REFERENCES `role` (id);
