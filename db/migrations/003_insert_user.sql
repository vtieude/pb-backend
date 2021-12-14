-- +goose Up
INSERT into user(username, password, email)
values("admin", "qwe@123", "admin@gmail.com"), ("vu", "qweqwe", "vule@gmail.com");