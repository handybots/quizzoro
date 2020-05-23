-- +goose Up
alter table users
    add column last_poll_id varchar(20) not null default ''
after state;

-- +goose Down
alter table users drop column last_poll_id;
