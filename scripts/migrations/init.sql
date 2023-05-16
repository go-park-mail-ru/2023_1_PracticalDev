CREATE TYPE account_type AS ENUM ('personal', 'business');
CREATE TYPE privacy AS ENUM ('public', 'secret');
CREATE TYPE notification_type AS ENUM ('new_pin', 'new_like', 'new_comment');

CREATE TABLE IF NOT EXISTS users
(
    id              serial       NOT NULL PRIMARY KEY,
    username        text         NOT NULL,
    email           text         NOT NULL UNIQUE,
    hashed_password bytea        NOT NULL,
    name            varchar(256) NOT NULL,
    profile_image   varchar,
    website_url     varchar,
    account_type    account_type NOT NULL
);

CREATE TABLE IF NOT EXISTS boards
(
    id          serial       NOT NULL PRIMARY KEY,
    name        varchar(256) NOT NULL,
    description varchar(500),
    privacy     privacy      NOT NULL,
    user_id     int          NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pins
(
    id                 serial    NOT NULL PRIMARY KEY,
    link               varchar(2048),
    title              varchar(100),
    description        varchar(500),
    created_at         timestamp NOT NULL DEFAULT now(),
    media_source       varchar   NOT NULL,
    media_source_color varchar   NOT NULL DEFAULT 'rgb(39, 102, 120)',
    n_likes            int       NOT NULL DEFAULT 0,
    author_id          int       NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pin_likes
(
    pin_id     int       NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    author_id  int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (pin_id, author_id)
);

CREATE TABLE IF NOT EXISTS boards_pins
(
    board_id int NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
    pin_id   int NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    PRIMARY KEY (board_id, pin_id)
);

CREATE TABLE IF NOT EXISTS comments
(
    id         serial    NOT NULL PRIMARY KEY,
    author_id  int       NOT NULL REFERENCES users (id),
    pin_id     int       NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    text       text      NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS comment_likes
(
    comment_id int       NOT NULL REFERENCES comments (id) ON DELETE CASCADE,
    author_id  int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (comment_id, author_id)
);

CREATE TABLE IF NOT EXISTS followings
(
    follower_id int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    followee_id int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at  timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (followee_id, follower_id)
);

CREATE TABLE IF NOT EXISTS chats
(
    id         serial    NOT NULL PRIMARY KEY,
    user1_id   int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    user2_id   int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT chats_user_pair UNIQUE (user1_id, user2_id)
);

CREATE TABLE IF NOT EXISTS messages
(
    id         serial    NOT NULL PRIMARY KEY,
    author_id  int       NOT NULL REFERENCES users (id),
    chat_id    int       NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    text       text      NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS notifications
(
    id         serial            NOT NULL PRIMARY KEY,
    user_id    int               NOT NULL REFERENCES users (id),
    type       notification_type NOT NULL,
    message    text,
    is_read    boolean           NOT NULL DEFAULT false,
    created_at timestamp         NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS new_pin_notifications
(
    notification_id int NOT NULL REFERENCES notifications (id),
    pin_id          int NOT NULL REFERENCES pins (id)
);

CREATE TABLE IF NOT EXISTS new_like_notifications
(
    notification_id int NOT NULL REFERENCES notifications (id),
    pin_id          int NOT NULL REFERENCES pins (id),
    author_id       int NOT NULL REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS new_comment_notifications
(
    notification_id int NOT NULL REFERENCES notifications (id),
    comment_id      int NOT NULL REFERENCES comments (id)
);

-- Обработка создания лайка
CREATE OR REPLACE FUNCTION on_pin_like() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE pins
    SET n_likes = n_likes + 1
    WHERE id = new.pin_id;

    RETURN NULL;
END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER pin_like
    AFTER INSERT
    ON pin_likes
    FOR EACH ROW
EXECUTE PROCEDURE on_pin_like();

-- Обработка удаления лайка
CREATE OR REPLACE FUNCTION on_pin_unlike() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE pins
    SET n_likes = n_likes - 1
    WHERE id = old.pin_id;

    RETURN NULL;
END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER pin_unlike
    AFTER DELETE
    ON pin_likes
    FOR EACH ROW
EXECUTE PROCEDURE on_pin_unlike();
