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
    name varchar(100) not null,
    alternative_name varchar(100),
    views integer default 0 not null,
    preview_url varchar(255) not null,
    uploader_id integer references users(id) on delete cascade
);

CREATE TABLE manga_chapter
(
    id SERIAL PRIMARY KEY,
    preview_url varchar(255) not null,
    manga_id integer references manga(id) on delete cascade,
    uploader_id integer references users(id) on delete cascade,
    number integer default 0
);

CREATE TABLE chapter_image
(
    id serial primary key,
    url varchar(255) not null,
    chapter_id integer references manga_chapter(id) on delete cascade
);

CREATE TABLE likes
(
    manga_id integer references manga(id) on delete cascade,
    user_id integer references users(id) on delete cascade
);

CREATE TABLE read_manga
(
    manga_id integer references manga(id) on delete cascade,
    user_id integer references users(id) on delete cascade
);