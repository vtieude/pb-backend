-- +goose Up

ALTER TABLE `product` 
ADD COLUMN `description` VARCHAR(250) NULL;
ALTER TABLE `product` 
ADD COLUMN `image_url` VARCHAR(250) NULL;
