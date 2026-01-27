package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nassabiq/golang-template/internal/modules/{{MODULE}}/domain"
)

type {{MODULE}}Repository struct {
	db *sql.DB
}

func New{{MODULE}}Repository(db *sql.DB) *{{MODULE}}Repository {
	return &{{MODULE}}Repository{db: db}
}

func (r *{{MODULE}}Repository) FindByID(ctx context.Context, id string) (*domain.{{MODULE}}, error) {
	var {{MODULE|lower}} domain.{{MODULE}}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, created_at, updated_at FROM {{MODULE|lower}}s WHERE id = $1",
		id).Scan(&{{MODULE|lower}}.ID, &{{MODULE|lower}}.Name, &{{MODULE|lower}}.CreatedAt, &{{MODULE|lower}}.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &{{MODULE|lower}}, nil
}

func (r *{{MODULE}}Repository) Create(ctx context.Context, request *domain.{{MODULE}}Create) (*domain.{{MODULE}}, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO {{MODULE|lower}}s (id, name, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())`,
		request.ID,
		request.Name,
	)

	if err != nil {
		return nil, err
	}

	return &domain.{{MODULE}}{
		ID:        request.ID,
		Name:      request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *{{MODULE}}Repository) List(ctx context.Context, limit int, offset int) ([]domain.{{MODULE}}, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name, created_at, updated_at FROM {{MODULE|lower}}s ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var {{MODULE|lower}}s []domain.{{MODULE}}
	for rows.Next() {
		var {{MODULE|lower}} domain.{{MODULE}}
		err := rows.Scan(&{{MODULE|lower}}.ID, &{{MODULE|lower}}.Name, &{{MODULE|lower}}.CreatedAt, &{{MODULE|lower}}.UpdatedAt)
		if err != nil {
			return nil, err
		}
		{{MODULE|lower}}s = append({{MODULE|lower}}s, {{MODULE|lower}})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return {{MODULE|lower}}s, nil
}

func (r *{{MODULE}}Repository) Update(ctx context.Context, request *domain.{{MODULE}}Update) (*domain.{{MODULE}}, error) {
	query := "UPDATE {{MODULE|lower}}s SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if request.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *request.Name)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, created_at, updated_at", argIndex)
	args = append(args, request.ID)

	var {{MODULE|lower}} domain.{{MODULE}}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&{{MODULE|lower}}.ID, &{{MODULE|lower}}.Name, &{{MODULE|lower}}.CreatedAt, &{{MODULE|lower}}.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &{{MODULE|lower}}, nil
}

func (r *{{MODULE}}Repository) Delete(ctx context.Context, {{MODULE|lower}} *domain.{{MODULE}}) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM {{MODULE|lower}}s WHERE id = $1", {{MODULE|lower}}.ID)
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
