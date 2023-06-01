-- +goose Up

ALTER TABLE `user` 
ADD COLUMN `role` VARCHAR(64) NOT NULL DEFAULT 'staff' AFTER `email`;
ALTER TABLE `user` 
ADD COLUMN `role_label` VARCHAR(64) NOT NULL DEFAULT 'Nhân Viên' AFTER `email`;

DROP TABLE user_role;
DROP TABLE role;

update user set `role` = "super_admin", role_label = "Super Admin" where 
id = (select userId.id from (select id from user where email = "admin@gmail.com")as userId limit 1);

update user set `role` = "admin", role_label = "Quản lí" where 
id = (select userId.id from (select id from user where email = "qa@gmail.com")as userId limit 1);

update user set `role` = "user", role_label = "Người dùng" where 
id = (select userId.id from (select id from user where email = "qa2@gmail.com")as userId limit 1);