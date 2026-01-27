-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
  id,
  name,
  email,
  password,
  role_id,
  created_at,
  updated_at
) VALUES (
  '00000000-0000-0000-0000-000000000010',
  'Super Admin',
  'superadmin@app.com',
  '$2a$10$2c3aB8KfW7KA7DaLW8BbQOScx/uoR3BIlQ0UJryweiMHHTQh9FjcW',
  '00000000-0000-0000-0000-000000000003',
  NOW(),
  NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = '00000000-0000-0000-0000-000000000010';
-- +goose StatementEnd
