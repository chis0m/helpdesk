-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auth_sessions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  session_id CHAR(36) NOT NULL UNIQUE,
  user_uuid CHAR(36) NOT NULL,
  refresh_jti CHAR(36) NOT NULL,
  refresh_expires_at DATETIME(3) NOT NULL,
  revoked_at DATETIME(3) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT fk_auth_sessions_user_uuid
    FOREIGN KEY (user_uuid) REFERENCES users(uuid)
    ON DELETE CASCADE ON UPDATE CASCADE,

  INDEX idx_auth_sessions_session_id (session_id),
  INDEX idx_auth_sessions_user_uuid (user_uuid),
  INDEX idx_auth_sessions_refresh_jti (refresh_jti),
  INDEX idx_auth_sessions_revoked_at (revoked_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth_sessions;
-- +goose StatementEnd
