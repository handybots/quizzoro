-- +goose Up
update users set orig_poll_id='0';
alter table users modify column orig_poll_id int not null default 0;

-- +goose Down
alter table users modify column orig_poll_id varchar(20) not null default '';
