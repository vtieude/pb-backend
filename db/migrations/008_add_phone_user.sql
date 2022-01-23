-- +goose Up

ALTER TABLE `user` 
ADD COLUMN `phone_number` VARCHAR(20) NULL;