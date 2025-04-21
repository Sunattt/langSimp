-- First create tables without dependencies
CREATE TABLE roles(
    id serial primary key ,
    role_type VARCHAR(20) not null unique,
    description varchar(35) not null,
    created_at timestamptz default current_timestamp
);

INSERT INTO roles (role_type, description) VALUES
    ( 'client', 'Обычный пользователь приложения'),
    ( 'moderator', 'Модератор контента'),
    ( 'admin', 'Администратор системы');

INSERT INTO user_photos (path, is_default) VALUES
    ('\image\profile_photo\default_photos\other.jpg', TRUE),
    ('\image\profile_photo\default_photos\female.jpg', TRUE),
    ('\image\profile_photo\default_photos\male.jpg', TRUE);


-- Add default photos before users table since users might reference them
CREATE TABLE user_photos (
    id SERIAL PRIMARY KEY,
    user_id BIGINT, -- Will be properly referenced after users table is created
    path TEXT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE languages (
    language_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(5) NOT NULL UNIQUE
);

INSERT INTO languages(name, code)
VALUES ('English', 'en');

-- Now create tables that depend on the above
CREATE TABLE users (
    id bigserial primary key ,
    username varchar(200) not null,
    email varchar(250) not null unique,
    password_hash varchar(255) not null,
    birthday DATE not null,
    gender int not null,
    photo_id INT REFERENCES user_photos(id) ON DELETE SET NULL DEFAULT 1,
    language_id int REFERENCES languages (language_id) ON DELETE SET NULL,
    role_id INT NOT NULL REFERENCES roles(id) DEFAULT 1,
    active boolean default true ,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);

-- Now update user_photos to add the proper foreign key
ALTER TABLE user_photos
    ADD CONSTRAINT fk_user_photos_users
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE TABLE chapters (
    chapter_id SERIAL PRIMARY KEY,
    language_id INTEGER REFERENCES languages(language_id) ON DELETE CASCADE,
    count_articles INT default 0,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    image_alt TEXT,
    active boolean default true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz default current_timestamp
);

CREATE TABLE articles (
    article_id SERIAL PRIMARY KEY,
    chapter_id INTEGER REFERENCES chapters(chapter_id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    level_id int REFERENCES grammar_levels (id) not null,
    description TEXT,
    image_url VARCHAR(255),
    image_alt TEXT,
    active boolean default true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz default current_timestamp
);

CREATE TABLE grammar_levels (
    id SERIAL PRIMARY KEY,
    level VARCHAR(50) NOT NULL
);

INSERT INTO grammar_levels(level) VALUES
    ('Beginner'),
    ('Intermediate'),
    ('Advanced');

CREATE TABLE grammar_contents (
    id SERIAL PRIMARY KEY,
    article_id INTEGER REFERENCES articles(article_id) ON DELETE CASCADE not null,
    level_id INTEGER REFERENCES grammar_levels(id) not null,
    explanation TEXT NOT NULL,
    structure TEXT,
    examples JSONB,
    tips TEXT,
    picture VARCHAR(255),
    active boolean default true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(article_id, level_id)
);

CREATE TABLE pronunciation_tables (
    id SERIAL PRIMARY KEY,
    grammar_content_id INTEGER REFERENCES grammar_contents(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    audio_url VARCHAR(512),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pronunciation_items (
    id SERIAL PRIMARY KEY,
    table_id INTEGER REFERENCES pronunciation_tables(id) ON DELETE CASCADE,
    word_or_phrase VARCHAR(255) NOT NULL,
    phonetic_transcription VARCHAR(255),
    audio_url VARCHAR(512),
    explanation TEXT,
    order_position INTEGER NOT NULL
);

CREATE TABLE grammar_exercises (
    id SERIAL PRIMARY KEY,
    grammar_content_id INTEGER REFERENCES grammar_contents(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL,
    options JSONB,
    correct_answer TEXT NOT NULL,
    explanation TEXT,
    difficulty INTEGER CHECK (difficulty BETWEEN 1 AND 3),
    help TEXT,
    active boolean default true
);

CREATE TABLE grammar_comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    grammar_content_id INTEGER REFERENCES grammar_contents(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    likes_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comment_likes (
    id SERIAL PRIMARY KEY,
    comment_id INTEGER NOT NULL REFERENCES grammar_comments(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(comment_id, user_id)
);

-- Finally create the saved chapters table
CREATE TABLE user_saved_chapters (
    save_id BIGSERIAL PRIMARY KEY ,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ,
    chapter_id BIGINT REFERENCES chapters(chapter_id) ON DELETE CASCADE,
    saved boolean default false,
    UNIQUE(user_id, chapter_id)
);

CREATE TABLE user_saved_articles (
    save_id BIGSERIAL PRIMARY KEY ,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ,
    article_id BIGINT REFERENCES articles(article_id) ON DELETE CASCADE,
    saved boolean default false,
    UNIQUE(user_id, article_id)
);
------------------------------------------------Vocabulary--------------------------------------------------------
/*

CREATE TABLE book_highlights (
   highlight_id BIGSERIAL PRIMARY KEY,
   book_id BIGINT REFERENCES books(book_id) ON DELETE CASCADE,
   user_id BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
   page_number INT,
   highlighted_text TEXT NOT NULL,
   note TEXT,
   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
*/
----------------------------------------------------------------------------------------------------------------------------

CREATE TABLE VocabularyWord (
    word_id BIGSERIAL primary key,

);





---------------------------------------Добавления индекса для ускорения запросов-------------------------------------------------------------

CREATE INDEX idx_grammar_articles_chapter_id ON articles(chapter_id);
CREATE INDEX idx_grammar_contents_article_id ON grammar_contents(article_id);
CREATE INDEX idx_grammar_contents_level_id ON grammar_contents(level_id);
CREATE INDEX idx_pronunciation_tables_content_id ON pronunciation_tables(grammar_content_id);
CREATE INDEX idx_grammar_comments_content_id ON grammar_comments(grammar_content_id);
CREATE INDEX idx_grammar_comments_user_id ON grammar_comments(user_id);
CREATE INDEX idx_grammar_comments_rating ON grammar_comments(rating);
CREATE INDEX idx_comment_likes_comment ON comment_likes(comment_id);
CREATE INDEX idx_comment_likes_user ON comment_likes(user_id);
CREATE INDEX idx_grammar_comments_likes ON grammar_comments(likes_count);