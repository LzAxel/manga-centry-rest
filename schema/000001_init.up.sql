CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username varchar(30) not null unique,
    password_hash varchar(255) not null,
    about varchar(255) not null,
    image varchar(255) default '' not null
);

CREATE TABLE manga
(
    id SERIAL PRIMARY KEY,
    name varchar(100) not null unique,
    views integer default 0 not null,
    preview varchar(255) not null,
    uploader_id integer references users(id) on delete cascade
);

CREATE TABLE manga_image
(
    id serial primary key,
    url varchar(255) not null,
    manga_id integer references manga(id) on delete cascade
);

