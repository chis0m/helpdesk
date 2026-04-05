-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ticket_comments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  ticket_id BIGINT UNSIGNED NOT NULL,
  author_user_id BIGINT UNSIGNED NOT NULL,
  body TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME(3) NULL,

  CONSTRAINT fk_ticket_comments_ticket_id
    FOREIGN KEY (ticket_id) REFERENCES tickets(id)
    ON DELETE CASCADE ON UPDATE CASCADE,

  CONSTRAINT fk_ticket_comments_author_user_id
    FOREIGN KEY (author_user_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE,

  INDEX idx_ticket_comments_ticket_id (ticket_id),
  INDEX idx_ticket_comments_author_user_id (author_user_id),
  INDEX idx_ticket_comments_deleted_at (deleted_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ticket_comments;
-- +goose StatementEnd
