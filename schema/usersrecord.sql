create table if not exists usersrecord(id uuid primary key default uuid_generate_v4() not null, username text, email text unique);
