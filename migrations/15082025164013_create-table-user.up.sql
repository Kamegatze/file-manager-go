create table if not exists users (
    id UUID primary key not null default gen_random_uuid (),
    last_name text not null,
    first_name text not null,
    username text not null
);