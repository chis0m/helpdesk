-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tickets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  reporter_user_id BIGINT UNSIGNED NOT NULL,
  assigned_user_id BIGINT UNSIGNED NULL,
  title VARCHAR(180) NOT NULL,
  description TEXT NOT NULL,
  category VARCHAR(80) NOT NULL DEFAULT 'general',
  status VARCHAR(30) NOT NULL DEFAULT 'open',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME(3) NULL,

  CONSTRAINT fk_tickets_reporter_user_id
    FOREIGN KEY (reporter_user_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE,

  CONSTRAINT fk_tickets_assigned_user_id
    FOREIGN KEY (assigned_user_id) REFERENCES users(id)
    ON DELETE SET NULL ON UPDATE CASCADE,

  INDEX idx_tickets_reporter_user_id (reporter_user_id),
  INDEX idx_tickets_assigned_user_id (assigned_user_id),
  INDEX idx_tickets_status (status),
  INDEX idx_tickets_category (category),
  INDEX idx_tickets_deleted_at (deleted_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tickets;
-- +goose StatementEnd
