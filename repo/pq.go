package repo

import (
	"context"
	"user-service/domain"

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
	ctx context.Context, id uuid.UUID, email string) (*domain.User, error) {

	sql := `INSERT INTO users (id, email) 
	VALUES ($1, $2) RETURNING id, email, created_at`

	var user *domain.User
	err := db.pg.QueryRow(
		ctx, sql, id, email).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete .
func (db *PostgresDB) Delete(ctx context.Context, id uuid.UUID) error {

	sql := `DELETE FROM users WHERE id=$1`

	_, err := db.pg.Exec(ctx, sql, id)
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
		var user *domain.User

		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
