-- name: CreateTask :execresult
insert into tasks (pid, status, summary, technician_id)
values (?, ?, ?, ?);

-- name: GetTaskWithPid :one
select *
from tasks
where pid = ?
order by id desc;

-- name: DeleteTask :exec
DELETE
FROM tasks
WHERE pid = ?;

-- name: GetAllTasks :many
select *
from tasks
order by id desc;

-- name: GetTaskWithTechId :many
select *
from tasks
where technician_id = ?
order by id desc;

-- name: CompleteTask :exec
update tasks
set status       = 'COMPLETED',
    completed_at = now()
where pid = ?;
