-- +goose Up
alter table users
    add column privacy bool not null default false
after state;

-- +goose Down
alter table users drop column privacy;
