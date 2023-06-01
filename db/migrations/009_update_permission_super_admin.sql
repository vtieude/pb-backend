-- +goose Up

update user set `permission` = 10 where id > 0 and `role` = "super_admin";