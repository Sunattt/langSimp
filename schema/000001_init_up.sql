CREATE TABLE profiles
(
    id serial not null unique ,
    user_id integer references users (id) on delete cascade not null,
    photo_profile text not null
);

CREATE TABLE topics
(
    id serial not null unique ,
    title text not null ,
    description text not null,
    priority integer not null,
    image text not null,
    count_article integer not null ,
    created_at timestamp not null default current_timestamp,
    updated_at  timestamp not null default current_timestamp,
    deleted_at  timestamp not null default current_timestamp
);

CREATE TABLE saved_chapter
(
    id serial not null unique ,
    profile_id int references profile_id (id ) on delete cascade not null,
    chapter_id int references chapter_id (id ) on delete cascade not null
);


CREATE TABLE chapters
(
    id serial not null unique ,
    title text not null ,
    description text not null,
    priority integer not null,
    image text not null,
    chapter_id int references chapter (id) on delete cascade not null,
    beginner boolean not null,
    intermediate boolean not null,
    advanced  boolean not null ,
    created_at timestamp not null default current_timestamp,
    updated_at  timestamp not null default current_timestamp,
    deleted_at  timestamp not null default current_timestamp
);

CREATE TABLE saved_chapter
(
    id serial not null unique ,
    profile_id int references profile_id (id ) on delete cascade not null,
    article_id int references chapter_id (id ) on delete cascade not null
);

CREATE TABLE topics
(
    id serial not null unique ,
    title text not null ,
    description text not null,
    priority integer not null,
    image text not null,
    comment
        created_at timestamp not null default current_timestamp,
    updated_at  timestamp not null default current_timestamp,
    deleted_at  timestamp not null default current_timestamp
);

CREATE TABLE commets
(
    id serial not null unique ,
    article_id int references article (id) on delete cascade no null ,
    profile_id int references users (id) on delete cascade no null,
    count_likes int ,
    content text ,
    created_at timestamp not null default current_timestamp,
    updated_at  timestamp not null default current_timestamp,
    deleted_at  timestamp not null default current_timestamp
);

CREATE TABLE users
(
    id int not null primary key ,
    username varchar(100) not null unique,
    email varchar(100) not null ,
    gender integer ,
    password_hash varchar(200) not null ,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,
    delated_at timestamp not null default CURRENT_TIMESTAMP
);