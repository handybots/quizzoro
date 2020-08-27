-- +goose Up
alter table polls convert to character set utf8mb4;

-- +goose Down
alter table polls convert to character set utf8;
