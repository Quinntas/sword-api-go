create table users
(
    id         int auto_increment
        primary key,
    pid        varchar(255)                       not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    username   varchar(255)                       not null,
    password   varchar(255)                       not null,
    role       varchar(255)                       not null,
    constraint users_pid_unique
        unique (pid),
    constraint users_username_unique
        unique (username)
);

create table tasks
(
    id            int auto_increment
        primary key,
    pid           varchar(255)                       not null,
    created_at    datetime default CURRENT_TIMESTAMP not null,
    technician_id int                                not null,
    summary       varchar(2500)                      not null,
    status        varchar(255)                       not null,
    completed_at  datetime                           ,
    constraint tasks_pid_unique
        unique (pid),
    constraint tasks_technician_id_users_id_fk
        foreign key (technician_id) references users (id)
);

