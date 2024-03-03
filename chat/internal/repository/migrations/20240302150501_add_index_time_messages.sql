-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE INDEX time_idx ON messages (time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX time_idx;
-- +goose StatementEnd