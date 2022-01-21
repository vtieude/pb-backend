-- +goose Up

ALTER TABLE `user` 
ADD COLUMN `permission` int(2) NOT NULL DEFAULT 1 AFTER `role_label`;

update user set `permission` = 3 where `role` = "staff";
update user set `permission` = 7 where `role` in ("super_admin", "admin");
