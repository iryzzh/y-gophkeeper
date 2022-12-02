-- noinspection SqlNoDataSourceInspectionForFile

create table if not exists items
(
    id         integer primary key autoincrement,
    user_id    text,
    meta       text,
    data_id    integer,
    data_type  text     default 'text',
    created_at datetime default current_timestamp,
    updated_at datetime default null,
    CONSTRAINT uniq UNIQUE (user_id, meta)
);

create table if not exists items_data
(
    id   integer primary key autoincrement,
    data blob not null
);