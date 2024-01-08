-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS files (
  id uuid PRIMARY KEY,
  name text NOT NULL,
  size integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
-- +goose StatementEnd
