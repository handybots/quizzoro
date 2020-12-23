CREATE DATABASE IF NOT EXISTS quizzoro;

CREATE TABLE IF NOT EXISTS quizzoro.logs (
    date Date,
    time DateTime,
    level String,
    event String,
    message String,
    user_id UInt32
) ENGINE = MergeTree(date, (level, event, user_id), 8192);
