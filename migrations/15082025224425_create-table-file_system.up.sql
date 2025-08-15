create table if not exists file_system (
    id UUID primary key not null default gen_random_uuid (),
    owner_id UUID not null references users (id),
    parent_id UUID references file_system (id),
    rights varchar(8) not null default 'rw------',
    is_file boolean not null,
    name text not null,
    path text not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted boolean not null default false
);