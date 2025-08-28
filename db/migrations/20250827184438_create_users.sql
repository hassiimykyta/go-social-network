-- +goose Up
-- +goose StatementBegin
create extension if not exists citext;

create table if not exists users (
    id bigint generated always as identity primary key,
    username citext unique not null,
    email citext unique not null,
    password_hash text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz null
);

create trigger trg_user_update_at
before update on users
for each row 
execute function set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger if exists trg_user_update_at on users;
drop table if exists users;
-- +goose StatementEnd
