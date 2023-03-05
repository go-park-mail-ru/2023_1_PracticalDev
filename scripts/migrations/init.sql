CREATE TYPE account_type AS ENUM ('personal', 'business');
CREATE TYPE privacy AS ENUM ('public', 'secret');

CREATE TABLE IF NOT EXISTS users
(
    id              serial       NOT NULL PRIMARY KEY,
    username        text         NOT NULL,
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
    board_id     int       NOT NULL
);

CREATE TABLE IF NOT EXISTS comments
(
    id          serial    NOT NULL PRIMARY KEY,
    description text,
    created_at  timestamp NOT NULL DEFAULT now(),
    pin_id      int       NOT NULL,
    user_id     int       NOT NULL
);

ALTER TABLE ONLY boards
    ADD CONSTRAINT fk_boards_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE ONLY pins
    ADD CONSTRAINT fk_pins_board_id
        FOREIGN KEY (board_id)
            REFERENCES boards (id);

ALTER TABLE ONLY comments
    ADD CONSTRAINT fk_pins_pin_id
        FOREIGN KEY (pin_id)
            REFERENCES pins (id);

ALTER TABLE ONLY comments
    ADD CONSTRAINT fk_pins_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);
