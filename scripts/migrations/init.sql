CREATE TABLE IF NOT EXISTS users (
    "user_id" serial NOT NULL PRIMARY KEY,
    "username" text NOT NULL,
    "email" text NOT NULL,
    "password" bytea NOT NULL
);
