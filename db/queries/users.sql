-- name: CreateUser :one
insert into users (email, username, password_hash)
values ($1, $2, $3)
returning users.*;

-- name: GetUserById :one
select users.* 
from users
where id = $1
limit 1;

-- name: GetUserByEmail :one
select users.*
from users
where email = $1 and deleted_at is null
limit 1;

-- name: GetUserByUsername :one
select users.*
from users
where username = $1 and deleted_at is null
limit 1;

-- name: ExistsUserByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL);

-- name: ExistsUserByUsername :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL);