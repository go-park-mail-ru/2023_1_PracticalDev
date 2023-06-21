CREATE TYPE account_type AS ENUM ('personal', 'business');
CREATE TYPE privacy AS ENUM ('public', 'secret');
CREATE TYPE notification_type AS ENUM ('new_pin', 'new_like', 'new_comment', 'new_follower');


-- В таблцие хранится информация о пользователях
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS users
(
    id              serial       NOT NULL PRIMARY KEY,
    username        text         NOT NULL,
    email           text         NOT NULL UNIQUE, -- идентификация пользователя происходит по email
    hashed_password bytea        NOT NULL,
    name            varchar(256) NOT NULL,
    profile_image   varchar,
    website_url     varchar,
    account_type    account_type NOT NULL
);

-- Таблица хранит информацию о досках (коллекции пинов)
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS boards
(
    id          serial       NOT NULL PRIMARY KEY,
    name        varchar(256) NOT NULL,
    description varchar(500),
    privacy     privacy      NOT NULL,
    user_id     int          NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- Таблциа хранит информацию о пинах
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов      +
-- 2NF: В таблице отсутствуют частичные зависимости     +
-- 3NF: В таблице отсутствуют транзитивные зависимости  +
-- Таблица соответствует 3НФ

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


-- Таблциа хранит информацию о лайках
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости + 
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS pin_likes
(
    pin_id     int       NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    author_id  int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (pin_id, author_id)
);

-- Таблциа хранит информацию о том, какие пины относятся к доске
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS boards_pins
(
    board_id int NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
    pin_id   int NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    PRIMARY KEY (board_id, pin_id)
);

-- Таблциа хранит информацию о комментариях к пинам
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS comments
(
    id         serial    NOT NULL PRIMARY KEY,
    author_id  int       NOT NULL REFERENCES users (id),
    pin_id     int       NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    text       text      NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);

-- Таблциа хранит информацию о подписках на пользователей
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS followings
(
    follower_id int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    followee_id int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at  timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (followee_id, follower_id)
);

-- Таблциа хранит информацию о чатах
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS chats
(
    id         serial    NOT NULL PRIMARY KEY,
    user1_id   int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    user2_id   int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT chats_user_pair UNIQUE (user1_id, user2_id) -- чат между двумя пользователями только один
);

-- Таблциа хранит информацию о сообщениях
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS messages
(
    id         serial    NOT NULL PRIMARY KEY,
    author_id  int       NOT NULL REFERENCES users (id),
    chat_id    int       NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    text       text      NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);

-- Таблциа хранит информацию о уведомлениях всех типов
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS notifications
(
    id         serial            NOT NULL PRIMARY KEY,
    user_id    int               NOT NULL REFERENCES users (id),
    type       notification_type NOT NULL,
    is_read    boolean           NOT NULL DEFAULT false,
    created_at timestamp         NOT NULL DEFAULT now()
);

-- Таблциа хранит информацию о уведомлениях на публикацию новых пинов
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS new_pin_notifications
(
    notification_id int NOT NULL REFERENCES notifications (id) ON DELETE CASCADE,
    pin_id          int NOT NULL REFERENCES pins (id) ON DELETE CASCADE
);

-- Таблциа хранит информацию о уведомлениях на лайк пина
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS new_like_notifications
(
    notification_id int NOT NULL REFERENCES notifications (id) ON DELETE CASCADE,
    pin_id          int NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    author_id       int NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- Таблциа хранит информацию о уведомлениях на новые комментарии
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS new_comment_notifications
(
    notification_id int  NOT NULL REFERENCES notifications (id) ON DELETE CASCADE,
    pin_id          int  NOT NULL REFERENCES pins (id) ON DELETE CASCADE,
    author_id       int  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    text            text NOT NULL
);


-- Таблциа хранит информацию о уведомлениях на новые подписки
-- Нормализация
-- 1NF: Таблица не содержит многозначных атрибутов     +
-- 2NF: В таблице отсутствуют частичные зависимости    +
-- 3NF: В таблице отсутствуют транзитивные зависимости +
-- Таблица соответствует 3НФ

CREATE TABLE IF NOT EXISTS new_follower_notifications
(
    notification_id int NOT NULL REFERENCES notifications (id) ON DELETE CASCADE,
    follower_id     int NOT NULL REFERENCES users (id) ON DELETE CASCADE
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

-- Удаление уведомления из-за удаления сущностей на которые ссылается уведомление
CREATE OR REPLACE FUNCTION on_specific_notification_delete() RETURNS TRIGGER AS
$$
BEGIN
    DELETE
    FROM notifications
    WHERE id = old.notification_id;

    RETURN NULL;
END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER new_pin_notification_delete
    AFTER DELETE
    ON new_pin_notifications
    FOR EACH ROW
EXECUTE PROCEDURE on_specific_notification_delete();

CREATE OR REPLACE TRIGGER new_comment_notification_delete
    AFTER DELETE
    ON new_comment_notifications
    FOR EACH ROW
EXECUTE PROCEDURE on_specific_notification_delete();

CREATE OR REPLACE TRIGGER new_like_notification_delete
    AFTER DELETE
    ON new_like_notifications
    FOR EACH ROW
EXECUTE PROCEDURE on_specific_notification_delete();

CREATE OR REPLACE TRIGGER new_follower_notification_delete
    AFTER DELETE
    ON new_follower_notifications
    FOR EACH ROW
EXECUTE PROCEDURE on_specific_notification_delete();


-- Indexes

-- Pins
CREATE INDEX IF NOT EXISTS pins_author_id_fk ON pins(author_id);

-- Boards
CREATE INDEX IF NOT EXISTS boards_user_id_fk ON boards(user_id);

-- Followings
CREATE INDEX IF NOT EXISTS followings_follower_id_fk ON followings(follower_id);
CREATE INDEX IF NOT EXISTS followings_followee_id_fk ON followings(followee_id);

-- Comments
CREATE INDEX IF NOT EXISTS comments_pin_id_fk ON comments(pin_id);

-- Likes
CREATE INDEX IF NOT EXISTS likes_author_id_fk on pin_likes(author_id);

-- Notifications
CREATE INDEX IF NOT EXISTS notifications_user_id_fk on notifications(user_id);

-- Chats
CREATE INDEX IF NOT EXISTS messages_chat_id_fk on messages(chat_id);

-- User search
CREATE INDEX IF NOT EXISTS user_search_by_username ON users(lower(username));
CREATE INDEX IF NOT EXISTS user_search_by_name ON users(lower(name));

-- Board search
CREATE INDEX IF NOT EXISTS board_search ON boards(lower(name));

-- Pin search
CREATE INDEX IF NOT EXISTS pin_search ON pins(lower(title));
