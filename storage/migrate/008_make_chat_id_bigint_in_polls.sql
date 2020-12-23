-- +goose Up
alter table polls modify column chat_id bigint not null default 0;

-- +goose Down
alter table polls modify column chat_id integer(10) not null default 0;
