-- +goose Up
-- +goose StatementBegin
CREATE TABLE msg (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE msg;
-- +goose StatementEnd
