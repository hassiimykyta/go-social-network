-- +goose Up
-- +goose StatementBegin
create table if not exists post_media (
    media_id bigint not null references media(id) on delete cascade,
    post_id bigint not null references posts(id) on delete cascade,
    position int not null default 0,
    primary key (post_id, media_id)
);

create index if not exists idx_post_media_post on post_media(post_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists post_media;
drop index if exists idx_post_media_post;
-- +goose StatementEnd
