package repo

import (
	"context"
	"user-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresDB .
type PostgresDB struct {
	pg *pgxpool.Pool
}

// NewPostgres .
func NewPostgres(pg *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{pg: pg}
}

// Add .
func (db *PostgresDB) Add(
	ctx context.Context, user *domain.User) (*domain.User, error) {

	sql := `INSERT INTO users (id, email, created_at) VALUES ($1, $2, $3)`

	_, err := db.pg.Exec(ctx, sql, user.ID.String(), user.Email, user.CreatedAt)
	// err := row.Scan(&user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete .
func (db *PostgresDB) Delete(ctx context.Context, id uuid.UUID) error {

	sql := `DELETE FROM users WHERE id=$1`

	_, err := db.pg.Exec(ctx, sql, id.String())
	return err
}

// List .
func (db *PostgresDB) List(ctx context.Context) ([]*domain.User, error) {

	sql := `SELECT id, email, created_at 
	FROM users ORDER BY created_at`

	rows, err := db.pg.Query(ctx, sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, 0)
	for rows.Next() {
		var (
			user domain.User
			id   string
		)

		if err := rows.Scan(&id, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		user.ID = uid
		users = append(users, &user)
	}

	return users, nil
}
