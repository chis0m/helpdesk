-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS invites (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(120) NOT NULL,
  token_hash CHAR(64) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  middle_name VARCHAR(100) NULL,
  invited_by_user_id BIGINT UNSIGNED NOT NULL,
  target_role VARCHAR(20) NOT NULL DEFAULT 'staff',
  expires_at DATETIME(3) NOT NULL,
  used_at DATETIME(3) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE INDEX idx_invites_token_hash (token_hash),
  INDEX idx_invites_email (email),
  INDEX idx_invites_expires_at (expires_at),
  INDEX idx_invites_used_at (used_at),

  CONSTRAINT fk_invites_invited_by_user_id
    FOREIGN KEY (invited_by_user_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invites;
-- +goose StatementEnd
