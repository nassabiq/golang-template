-- +goose Up
-- +goose StatementBegin
CREATE TABLE password_resets (
  id          VARCHAR(36) PRIMARY KEY,
  user_id     VARCHAR(36) NOT NULL,
  token_hash  TEXT NOT NULL,
  used        BOOLEAN NOT NULL DEFAULT false,
  expires_at  TIMESTAMP NOT NULL,
  created_at  TIMESTAMP NULL,
  updated_at  TIMESTAMP NULL,

  CONSTRAINT fk_password_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_password_resets_hash ON password_resets(token_hash);
CREATE INDEX idx_password_resets_user ON password_resets(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_password_resets_hash;
DROP INDEX idx_password_resets_user;
DROP TABLE password_resets;
-- +goose StatementEnd
