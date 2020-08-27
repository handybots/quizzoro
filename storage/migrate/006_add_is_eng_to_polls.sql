-- +goose Up
alter table polls
add column is_eng boolean not null default false
after difficulty;

-- +goose Down
alter table polls drop column is_eng;
