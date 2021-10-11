CREATE TABLE IF NOT EXISTS  users (
    id bigint(50) not null PRIMARY KEY,
    name text not null,
    email text not null,
    password text not null,
    password_salt varchar(10) not null
);


