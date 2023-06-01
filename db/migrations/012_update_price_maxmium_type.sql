-- +goose Up
ALTER TABLE `product` MODIFY `selling_price`   DECIMAL(15,1) Not NULL;
ALTER TABLE `product` MODIFY `price`   DECIMAL(15,1) Not NULL;