-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id            VARCHAR(36) PRIMARY KEY,
  name          VARCHAR(255) NOT NULL,
  email         VARCHAR(255) NOT NULL UNIQUE,
  password      TEXT NOT NULL,
  role_id       VARCHAR(36) NOT NULL,
  created_at    TIMESTAMP NULL,
  updated_at    TIMESTAMP NULL,

  CONSTRAINT fk_users_role
    FOREIGN KEY (role_id)
    REFERENCES roles(id)
    ON DELETE RESTRICT
);

CREATE INDEX idx_users_role_id ON users(role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_role_id;
DROP TABLE users;
-- +goose StatementEnd
