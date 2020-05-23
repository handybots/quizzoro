-- +goose Up

create table users (
    created_at      timestamp       not null default now()                      ,
    updated_at      timestamp       not null default now() on update now()      ,
    id              integer(10)     not null primary key                        ,
    state           varchar(32)     not null default 'default'                  ,
    last_message_id integer(10)     not null default 0                          ,
    last_category   varchar(16)     not null default ''
);

create table polls (
    created_at      timestamp       not null default now()                      ,
    updated_at      timestamp       not null default now() on update now()      ,
    id              varchar(20)     not null primary key                        ,
    message_id      varchar(10)     not null default ''                         ,
    chat_id         integer(10)     not null default 0                          ,
    category        varchar(32)     not null default ''                         ,
    difficulty      varchar(16)     not null default ''                         ,

    question        text            not null                                    ,
    question_eng    text            not null                                    ,
    correct         text            not null                                    ,
    correct_eng     text            not null                                    ,
    answers         json            not null                                    ,
    answers_eng     json            not null
);

create table passed_polls (
    created_at      timestamp       not null default now()                      ,
    updated_at      timestamp       not null default now() on update now()      ,
    user_id         int(10)         not null default 0                          ,
    poll_id         varchar(20)     not null default ''                         ,
    correct         bool            not null default false                      ,

    primary key (user_id, poll_id)
);

-- +goose Down

drop table users;
drop table polls;
drop table passed_polls;
