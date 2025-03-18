-- name: GetAllUsers :many
select *
from users
order by id desc;

-- name: GetUserWithUsername :one
select *
from users
where username = ?
order by id desc;

-- name: GetUserWithPid :one
select *
from users
where pid = ?
order by id desc;

-- name: CreateUser :execresult
insert into users (pid, username, password, role)
values (?, ?, ?, ?);
