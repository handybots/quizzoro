-- +goose Up
alter table users modify id bigint;
alter table passed_polls modify user_id bigint;

-- +goose Down
alter table users modify id integer(10);
alter table passed_polls modify user_id integer(10);
