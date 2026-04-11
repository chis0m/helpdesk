-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  actor_user_id BIGINT UNSIGNED NULL,
  action VARCHAR(160) NOT NULL,
  path VARCHAR(512) NOT NULL DEFAULT '',
  http_method VARCHAR(16) NOT NULL DEFAULT '',
  session_id CHAR(36) NULL,
  token_jti VARCHAR(64) NULL,
  success BOOLEAN NOT NULL DEFAULT TRUE,
  error_code VARCHAR(96) NULL,
  resource_type VARCHAR(64) NULL,
  resource_id BIGINT UNSIGNED NULL,
  ip VARCHAR(64) NULL,
  user_agent VARCHAR(512) NULL,
  metadata JSON NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  INDEX idx_audit_logs_action (action),
  INDEX idx_audit_logs_actor_user_id (actor_user_id),
  INDEX idx_audit_logs_created_at (created_at),
  INDEX idx_audit_logs_resource (resource_type, resource_id),

  CONSTRAINT fk_audit_logs_actor_user_id
    FOREIGN KEY (actor_user_id) REFERENCES users(id)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_logs;
-- +goose StatementEnd
