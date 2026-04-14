-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  actor_user_uuid CHAR(36) NULL,
  action VARCHAR(160) NOT NULL,
  path VARCHAR(512) NOT NULL DEFAULT '',
  http_method VARCHAR(16) NOT NULL DEFAULT '',
  resource_type VARCHAR(64) NULL,
  resource_id BIGINT UNSIGNED NULL,
  session_id CHAR(36) NULL,
  token_jti VARCHAR(64) NULL,
  success BOOLEAN NOT NULL DEFAULT TRUE,
  error_code VARCHAR(96) NULL,
  ip VARCHAR(64) NULL,
  user_agent VARCHAR(512) NULL,
  metadata JSON NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  INDEX idx_audit_logs_action (action),
  INDEX idx_audit_logs_actor_user_uuid (actor_user_uuid),
  INDEX idx_audit_logs_created_at (created_at),
  INDEX idx_audit_logs_resource (resource_type, resource_id),

  CONSTRAINT fk_audit_logs_actor_user_uuid
    FOREIGN KEY (actor_user_uuid) REFERENCES users(uuid)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_logs;
-- +goose StatementEnd
