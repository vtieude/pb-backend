-- +goose Up

CREATE TABLE IF NOT EXISTS `customer` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `customer_name` int(10) Not NULL,
        `phone` int(10) NULL,
        `address` DECIMAL(9,1) NULL,
        `active`  TINYINT(1) NOT NULL DEFAULT 1,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
        PRIMARY KEY (`id`)
);
CREATE TABLE IF NOT EXISTS `sale` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `fk_user` int(10) Not NULL,
        `fk_customer` int(10) NULL,
        `fk_product` int(10) NULL,
        `price` DECIMAL(9,1) Not NULL,
        `note` VARCHAR(100)  NULL,
        `status` enum('pending','saled', 'failed') default 'pending',
        `saled_date`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `active`  TINYINT(1) NOT NULL DEFAULT 1,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
        PRIMARY KEY (`id`)
);

ALTER TABLE sale ADD FOREIGN KEY (fk_user) REFERENCES user (id);
ALTER TABLE sale ADD FOREIGN KEY (fk_product) REFERENCES `product` (id);
ALTER TABLE sale ADD FOREIGN KEY (fk_customer) REFERENCES `customer` (id);
ALTER TABLE user MODIFY `username`  VARCHAR(64) Not NULL;
ALTER TABLE `role` MODIFY `label`  VARCHAR(64) Not NULL;
ALTER TABLE product ADD COLUMN `price` DECIMAL(9,1) Not NULL;