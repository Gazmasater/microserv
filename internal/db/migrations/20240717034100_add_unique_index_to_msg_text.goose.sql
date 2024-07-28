-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_msg_text ON msg (text);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_msg_text;
-- +goose StatementEnd
