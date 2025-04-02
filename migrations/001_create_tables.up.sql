-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE IF NOT EXISTS Posts (
    post_id  BIGSERIAL  PRIMARY KEY,
    author_id UUID NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    comments_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    create_date TIMESTAMP  WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Comments (
    comment_id BIGSERIAL  PRIMARY KEY,
    author_id UUID NOT NULL,
    post_id BIGINT NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES comments(comment_id) ON DELETE CASCADE,
    path LTREE UNIQUE NOT NULL,
    replies_level INTEGER NOT NULL,
    text TEXT NOT NULL,
    create_date TIMESTAMP  WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX path_gist_idx ON Comments USING GIST (path);

CREATE INDEX create_date_idx ON Comments (create_date);

CREATE INDEX post_idx ON Comments (post_id);

-- +goose StatementEnd

-- test data
-- -- +goose StatementBegin
-- INSERT INTO posts ( title, text, author_id, create_date)
-- VALUES (
--     'Тестовый пост про SQL', 
--     'Этот пост создан для проверки API комментариев', 
--     'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 
--     NOW()
-- );


-- INSERT INTO comments ( post_id, author_id, text, path, replies_level, create_date)
-- VALUES (
--     1,
--     'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
--     'Первый комментарий к посту!',
--     '1',
--     1,
--     NOW() - INTERVAL '1 hour'
-- );


-- INSERT INTO comments ( post_id, author_id, text, path, replies_level, create_date)
-- VALUES (
--     1,
--     'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
--     'Второй корневой комментарий',
--     '2',
--     1,
--     NOW() - INTERVAL '45 minutes'
-- );

-- INSERT INTO comments ( post_id, parent_id, author_id, text, path, replies_level, create_date)
-- VALUES (
--     1,
--     1,
--     'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
--     'Это ответ на первый комментарий',
--     '1.3',
--     2,
--     NOW() - INTERVAL '30 minutes'
-- );


-- INSERT INTO comments ( post_id, parent_id, author_id, text, path, replies_level, create_date)
-- VALUES (
--     1,
--     3,
--     'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55',
--     'А это ответ на ответ (уровень 3)',
--     '1.3.4',
--     3,
--     NOW() - INTERVAL '15 minutes'
-- );


-- INSERT INTO comments ( post_id, parent_id, author_id, text, path, replies_level, create_date)
-- VALUES (
--     1,
--     1,
--     'f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66',
--     'Ещё один ответ к первому комментарию',
--     '1.5',
--     2,
--     NOW() - INTERVAL '5 minutes'
-- );


-- INSERT INTO comments (post_id, parent_id, author_id, text, path, replies_level, create_date)
-- VALUES (
--     1,
--     2,
--     'a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77',
--     'Ответ на второй корневой комментарий',
--     '2.6',
--     2,
--     NOW()
-- );


-- INSERT INTO comments ( post_id, parent_id, author_id, text, path, replies_level, create_date)
-- VALUES (
--     1,
--     1,
--     'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
--     'Это ответ на первый комментарий',
--     '1.7',
--     2,
--     NOW() - INTERVAL '30 minutes'
-- );
-- -- +goose StatementEnd