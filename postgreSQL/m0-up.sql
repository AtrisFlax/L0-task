create table if not exists task.items
(
    uuid    uuid  not null
        constraint items_1_pkey
            primary key,
    payload jsonb not null
);

alter table task.items
    owner to my_user_name;

create unique index if not exists items_uuid_uindex
    on task.items (uuid);

