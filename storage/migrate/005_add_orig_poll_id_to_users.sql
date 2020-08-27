-- +goose Up
alter table users
add column orig_poll_id varchar(20) not null default ''
after privacy;

-- +goose Down
alter table users drop column orig_poll_id;
