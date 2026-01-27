package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nassabiq/golang-template/internal/modules/user/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, role_id, created_at, updated_at FROM users WHERE id = $1",
		id).Scan(&user.ID, &user.Name, &user.Email, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, request *domain.UserCreate) (*domain.User, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, name, email, password_hash, role_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		request.ID,
		request.Name,
		request.Email,
		request.Password,
		request.RoleID,
	)

	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        request.ID,
		Name:      request.Name,
		Email:     request.Email,
		RoleID:    request.RoleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *UserRepository) List(ctx context.Context, limit int, offset int) ([]domain.User, int64, error) {
	// Count total
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name, email, role_id, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, request *domain.UserUpdate) (*domain.User, error) {
	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if request.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *request.Name)
		argIndex++
	}

	if request.Email != nil {
		query += fmt.Sprintf(", email = $%d", argIndex)
		args = append(args, *request.Email)
		argIndex++
	}

	if request.RoleID != nil {
		query += fmt.Sprintf(", role_id = $%d", argIndex)
		args = append(args, *request.RoleID)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, email, role_id, created_at, updated_at", argIndex)
	args = append(args, request.ID)

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.Name, &user.Email, &user.RoleID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, user *domain.User) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
