-- +goose Up
-- +goose StatementBegin
ALTER TABLE tickets
  ADD COLUMN uuid CHAR(36) NOT NULL,
  ADD UNIQUE INDEX ux_tickets_uuid (uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tickets
  DROP INDEX ux_tickets_uuid,
  DROP COLUMN uuid;
-- +goose StatementEnd
