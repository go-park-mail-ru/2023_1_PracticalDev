CREATE TYPE account_type AS ENUM ('personal', 'business');
CREATE TYPE privacy AS ENUM ('public', 'secret');

CREATE TABLE IF NOT EXISTS users
(
    id              serial       NOT NULL PRIMARY KEY,
    username        text         NOT NULL UNIQUE,
    email           text         NOT NULL,
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
    user_id     int          NOT NULL
);

CREATE TABLE IF NOT EXISTS pins
(
    id           serial    NOT NULL PRIMARY KEY,
    link         varchar(2048),
    title        varchar(100),
    description  varchar(500),
    created_at   timestamp NOT NULL DEFAULT now(),
    media_source varchar,
    n_likes      int       NOT NULL DEFAULT 0,
    author_id    int       NOT NULL
);

CREATE TABLE IF NOT EXISTS boards_pins
(
    board_id int REFERENCES boards (id) ON DELETE CASCADE,
    pin_id   int REFERENCES pins (id) ON DELETE CASCADE,
    PRIMARY KEY (board_id, pin_id)
);

CREATE TABLE IF NOT EXISTS comments
(
    id          serial    NOT NULL PRIMARY KEY,
    description text,
    created_at  timestamp NOT NULL DEFAULT now(),
    pin_id      int       NOT NULL,
    user_id     int       NOT NULL
);

CREATE TABLE IF NOT EXISTS pin_likes
(
    pin_id     int REFERENCES pins (id) ON DELETE CASCADE,
    author_id  int REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (pin_id, author_id)
);

ALTER TABLE ONLY boards
    ADD CONSTRAINT fk_boards_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE ONLY pins
    ADD CONSTRAINT fk_pins_author_id
        FOREIGN KEY (author_id)
            REFERENCES users (id);

ALTER TABLE ONLY comments
    ADD CONSTRAINT fk_pins_pin_id
        FOREIGN KEY (pin_id)
            REFERENCES pins (id);

ALTER TABLE ONLY comments
    ADD CONSTRAINT fk_pins_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE ONLY comments
    DROP CONSTRAINT fk_pins_pin_id,
    ADD CONSTRAINT fk_pins_pin_id
        FOREIGN KEY (pin_id) REFERENCES pins (id) ON DELETE CASCADE;

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
