-- +goose Up

INSERT into user(username, password, email)
values("admin", "$2a$14$uAMQMJTteRe5oPvgoaGIgeWMujBvgUdud2kbsS1l5yU0.AO6o/qDO", "admin@gmail.com");

INSERT into `role`(role_name, label)
values("admin", "Admin");

INSERT into user_role(fk_role, fk_user)
values((select id from user where email = "admin@gmail.com" limit 1), (select id from role where role_name = "admin" limit 1));
INSERT into `role`(role_name, label)
values("staff", "Nhân Viên"), ("user", "Người dùng");
