-- +goose Up
-- +goose StatementBegin
ALTER TABLE ticket_comments
  ADD COLUMN uuid CHAR(36) NOT NULL,
  ADD UNIQUE INDEX ux_ticket_comments_uuid (uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ticket_comments
  DROP INDEX ux_ticket_comments_uuid,
  DROP COLUMN uuid;
-- +goose StatementEnd
