-- +goose Up

ALTER TABLE `product` 
ADD COLUMN `category` VARCHAR(50) NULL;
ALTER TABLE `product` 
ADD COLUMN `selling_price` DECIMAL(9,1) Not NULL;
ALTER TABLE `product` 
ADD COLUMN `quantity` INT(10) NOT NULL default 0;

CREATE TABLE IF NOT EXISTS `product_check_point` (
        `id` INT(10) NOT NULL AUTO_INCREMENT,
        `fk_product` int(10) Not NULL,
        `fk_sale` int(10) NULL,
        `type`  enum('saled','added') default 'added',
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP  ON UPDATE CURRENT_TIMESTAMP,
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
         `active`  TINYINT(1) NOT NULL DEFAULT 1,
        PRIMARY KEY (`id`)
    ) ;
ALTER TABLE product_check_point ADD FOREIGN KEY (fk_product) REFERENCES product (id);
ALTER TABLE product_check_point ADD FOREIGN KEY (fk_sale) REFERENCES sale (id);

