-- +goose Up
alter table polls change id poll_id varchar(20) not null default '';
alter table polls drop primary key;
alter table polls add column id integer primary key auto_increment after updated_at;
alter table polls add column deleted tinyint(1) not null default false after updated_at;

-- +goose Down
