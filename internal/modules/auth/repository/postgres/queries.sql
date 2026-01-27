-- name: CreateUser
INSERT INTO users (id, name, email, password, role_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FindUserByEmail
SELECT id, name, email, password, role_id, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: FindUserByID
SELECT id, name, email, password, role_id, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;


-- name: UpdateUserPassword
UPDATE users SET password = $1 WHERE id = $2;

-- name: StoreRefreshToken
INSERT INTO refresh_tokens (id, user_id, token_hash, revoked, expires_at, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FindValidRefreshToken
SELECT id, user_id, token_hash, revoked, expires_at, created_at, updated_at FROM refresh_tokens 
WHERE token_hash = $1 AND revoked = false AND expires_at > NOW() LIMIT 1;

-- name: RevokeRefreshToken
UPDATE refresh_tokens SET revoked = true
WHERE token_hash = $1;

-- name: RevokeAllRefreshTokens
UPDATE refresh_tokens SET revoked = true 
WHERE user_id = $1;

-- name: StorePasswordReset
INSERT INTO password_resets (id, user_id, token_hash, expires_at, used, created_at, updated_at) 
VALUES ($1, $2, $3, $4, false, $5, $6);

-- name: FindValidPasswordReset
SELECT id, user_id, token_hash, expires_at, used, created_at, updated_at 
FROM password_resets WHERE token_hash = $1 AND used = false AND expires_at > NOW() LIMIT 1;

-- name: MarkPasswordResetUsed
UPDATE password_resets
SET used = true WHERE id = $1;

-- name: FindRoleByName
SELECT id, name FROM roles WHERE name = $1 LIMIT 1;

-- name: FindRoleByID
SELECT id, name FROM roles WHERE id = $1 LIMIT 1;

-- users
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
ON users(email);

-- refresh_tokens
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash
ON refresh_tokens(token_hash);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user
ON refresh_tokens(user_id);

-- password_resets
CREATE INDEX IF NOT EXISTS idx_password_resets_hash
ON password_resets(token_hash);










