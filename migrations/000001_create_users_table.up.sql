-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS users
(
    user_id text check (length(user_id) == 36),
    login    text check (length(login) > 1),
    password text check (length(login) > 1),
    CONSTRAINT login_unique UNIQUE (login),
    CONSTRAINT uid_pkey PRIMARY KEY (user_id)
);