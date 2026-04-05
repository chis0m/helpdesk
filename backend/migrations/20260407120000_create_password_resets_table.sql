-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_resets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_uuid CHAR(36) NOT NULL,
  token_hash CHAR(64) NOT NULL,
  expires_at DATETIME(3) NOT NULL,
  used_at DATETIME(3) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE INDEX idx_password_resets_token_hash (token_hash),
  INDEX idx_password_resets_user_uuid (user_uuid),
  INDEX idx_password_resets_expires_at (expires_at),
  INDEX idx_password_resets_used_at (used_at),

  CONSTRAINT fk_password_resets_user_uuid
    FOREIGN KEY (user_uuid) REFERENCES users(uuid)
    ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS password_resets;
-- +goose StatementEnd
