-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS users
(
    user_id  text unique,
    login    text unique,
    password text,
    CONSTRAINT uid_pkey PRIMARY KEY (user_id)
);