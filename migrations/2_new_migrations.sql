-- +goose Up
CREATE TABLE IF NOT EXISTS companies
(
    id         UUID               DEFAULT uuid_generate_v4() PRIMARY KEY,
    name       VARCHAR   NOT NULL,
    owner_id   UUID      NOT NULL
        CONSTRAINT companies_owner_id_fk REFERENCES users (id) ON DELETE CASCADE,
    manager_id UUID REFERENCES users (id),
    is_active  BOOLEAN   NOT NULL DEFAULT true,
    logo_url   VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS companies;
