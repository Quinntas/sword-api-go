create table users
(
    id         serial primary key,
    pid        UUID         not null default uuid_generate_v4(),
    email      varchar(255) not null,
    password   varchar(255) not null,
    created_at timestamptz  not null default now()
);