CREATE TABLE users
(
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE
);

CREATE TABLE tags
(
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE users_tags(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    tag_id INTEGER REFERENCES tags(id)
);

CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    description VARCHAR(100)
);

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id),
    token VARCHAR(500) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO users (name, is_admin) VALUES
('user1', false),
('user2', true),
('user3', false),
('user4', false),
('user5', true);

INSERT INTO tags (name) VALUES 
('male'),
('female'),
('city'),
('village'),
('young'),
('old'),
('one more tag'),
('yet one more tag');

INSERT INTO users_tags (user_id, tag_id) VALUES 
(1, 1),
(1, 3),
(1, 5),
(2, 2),
(2, 7),
(3, 1),
(4, 2),
(5, 2),
(5, 6),
(5, 8);

INSERT INTO features (description) VALUES
('desc1'),
('some desc'),
('desc2'),
('lol'),
('desc3'),
('kek');