-- name: GetAllUsers :many
select *
from users
order by id desc;

-- name: CreateUser :one
insert into users (email, password)
values ($1, $2) returning *;
