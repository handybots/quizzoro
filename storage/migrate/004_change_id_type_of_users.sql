-- +goose Up
alter table users modify id bigint;

-- +goose Down
alter table users modify id integer(10);
