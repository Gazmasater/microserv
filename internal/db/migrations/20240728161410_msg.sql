-- +goose Up
-- +goose StatementBegin
CREATE TABLE msg (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL UNIQUE,
    status_1 VARCHAR(50) NOT NULL,
    status_2 VARCHAR(50),  
    created_at_1 TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_2 TIMESTAMPTZ DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS msg;
-- +goose StatementEnd