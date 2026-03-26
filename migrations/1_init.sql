-- +goose Up

CREATE TABLE IF NOT EXISTS user_type
(
    code     VARCHAR PRIMARY KEY,
    name     VARCHAR NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false
);

INSERT INTO user_type (code, name, is_admin) VALUES ('DEFAULT', 'Default User Type', false);
INSERT INTO user_type (code, name, is_admin) VALUES ('ADMIN', 'Admin User Type', true);

CREATE TABLE IF NOT EXISTS users
(
    id           UUID                      DEFAULT uuid_generate_v4() PRIMARY KEY,
    name         VARCHAR          NOT NULL,
    username     VARCHAR          NOT NULL
        CONSTRAINT user_username_unique UNIQUE,
    password     VARCHAR          NOT NULL,
    public_name  VARCHAR,
    description  VARCHAR,
    user_type    VARCHAR          NOT NULL REFERENCES user_type (code),
    age          INTEGER          NOT NULL DEFAULT 1,
    initial_age  INTEGER,
    rate         DOUBLE PRECISION NOT NULL DEFAULT 0,
    last_rate    DOUBLE PRECISION,
    is_active    BOOLEAN          NOT NULL DEFAULT true,
    read_message BOOLEAN,
    balance      DECIMAL          NOT NULL DEFAULT 0,
    lock_balance DECIMAL,
    last_login   TIMESTAMP,
    created_at   TIMESTAMP        NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at   TIMESTAMP        NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE INDEX IF NOT EXISTS users_age_index ON users (age);
CREATE INDEX IF NOT EXISTS users_name_balance_index ON users (name, balance);

-- +goose Down
DROP TABLE IF EXISTS users;
