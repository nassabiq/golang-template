-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
  id          VARCHAR(36) PRIMARY KEY,
  user_id     VARCHAR(36) NOT NULL,
  token_hash  TEXT NOT NULL,
  revoked     BOOLEAN NOT NULL DEFAULT false,
  expires_at  TIMESTAMP NOT NULL,
  created_at  TIMESTAMP NULL,
  updated_at  TIMESTAMP NULL,

  CONSTRAINT fk_refresh_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_refresh_tokens_hash;
DROP INDEX idx_refresh_tokens_user;
DROP TABLE refresh_tokens;
-- +goose StatementEnd
