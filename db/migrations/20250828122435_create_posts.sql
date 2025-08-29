-- +goose Up
-- +goose StatementBegin
create table if not exists posts (
   id bigint generated always as identity primary key,
   title varchar(255) not null,
   description text not null default '',
   user_id bigint not null references users(id) on delete cascade,
   created_at timestamptz not null default now(),
   updated_at timestamptz not null default now(),
   deleted_at timestamptz null
);

create trigger trg_post_update_at
before update on posts
for each row 
execute function set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger if exists trg_post_update_at on posts;
drop table if exists posts;
-- +goose StatementEnd
