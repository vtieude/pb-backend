-- +goose Up

ALTER TABLE sale
CHANGE COLUMN status sale_status  enum('pending','saled', 'failed') default 'pending';