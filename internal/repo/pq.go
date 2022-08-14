package repo

import (
	"context"
	"user-service/models"

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
	ctx context.Context, user *models.User) (*models.User, error) {

	sql := `INSERT INTO users (id, email, created_at) VALUES ($1, $2, $3)`

	_, err := db.pg.Exec(ctx, sql, user.UID, user.Email, user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete .
func (db *PostgresDB) Delete(ctx context.Context, uid string) error {

	sql := `DELETE FROM users WHERE id=$1`

	_, err := db.pg.Exec(ctx, sql, uid)
	return err
}

// List .
func (db *PostgresDB) List(ctx context.Context) ([]*models.User, error) {

	sql := `SELECT id, email, created_at 
	FROM users ORDER BY created_at`

	rows, err := db.pg.Query(ctx, sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)
	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.UID, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
