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

-- name: ListPostsWithMediaPaginated :many
select p.*,
  m.id          AS media_id,
  m.kind        AS media_kind,
  m.mime_type   AS media_mime_type,
  m.storage_key AS media_storage_key,
  m.width       AS media_width,
  m.height      AS media_height,
  m.duration_ms AS media_duration_ms,
  pm.position   AS media_position
from posts p 
left join post_media pm
on pm.post_id = p.id
left join media m
on m.id = pm.media_id
and m.deleted_at is null
where p.deleted_at is null
and  (sqlc.narg('user_id')::bigint IS NULL OR p.user_id = sqlc.narg('user_id')::bigint)
order by p.created_at desc, p.id desc,
pm.position asc, m.created_at asc, m.id asc
limit sqlc.arg('limit') 
offset sqlc.arg('offset');


