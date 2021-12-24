#!/usr/bin/env bash
echo 'Runing migrations...'
goose -dir db/migrations mysql "root:qweqwe@tcp(${MYSQL_HOST})/app-db?parseTime=true" up &
echo 'Deleting mysql-client...'
apk del mysql-client
./pb-backend
