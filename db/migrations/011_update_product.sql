-- +goose Up

ALTER TABLE `product` 
ADD COLUMN `description` VARCHAR(250) NULL;
ALTER TABLE `product` 
ADD COLUMN `image_prefix` VARCHAR(250) NULL;
