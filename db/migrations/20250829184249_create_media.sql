-- +goose Up
-- +goose StatementBegin
create table if not exists media(
    id bigint generated always as identity primary key,
    owner_id bigint not null references users(id) on delete cascade,
    kind varchar(20) not null,
    storage_key text not null,
    mime_type varchar(100) not null,
    size_bytes bigint not null,
    width int null,
    height int null,
    duration_ms int null,
    created_at timestamptz not null default now(),
    deleted_at timestamptz null
);

create unique index if not exists ux_media_storage_key on media(storage_key);
create index if not exists idx_media_owner on media(owner_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists media;
drop index if exists ux_media_storage_key;
drop index if exists idx_media_owner;
-- +goose StatementEnd
