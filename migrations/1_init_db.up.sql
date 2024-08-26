create table if not exists public.user
(
    id         bigserial primary key,
    username   varchar unique not null,
    password   varchar        not null,
    created_at timestamp      not null default now()
);

create table if not exists note
(
    id         bigserial primary key,
    username   varchar references public.user (username),
    title      varchar   not null,
    text       text               default '',
    created_at timestamp not null default now()
)