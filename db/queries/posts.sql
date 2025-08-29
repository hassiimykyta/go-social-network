-- name: CreatePost :one
insert into posts (title, description, user_id)
values ($1, $2, $3)
returning posts.*;

-- name: GetAllPostsByUser :many
select posts.*
from posts
where user_id = $1
and deleted_at is null
order by created_at desc, id desc;

-- name: SoftDeletePost :exec
update posts
set deleted_at = now()
where id = $1 and deleted_at is null;

-- name: UpdatePostPartial :one
update posts
set title = coalesce(sqlc.narg(title), title),
description = coalesce (sqlc.narg(description), description)
where id = sqlc.arg(id)
and deleted_at is null
returning posts.*;

-- name: ListPostsPaginated :many
select posts.*
from posts
where deleted_at is null
order by created_at desc, id desc
limit $1 offset $2;
