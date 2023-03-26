-- +goose Up
ALTER TABLE `product` 
ADD COLUMN `image_base64` LONGBLOB NULL;
