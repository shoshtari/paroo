package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type UserRepoImp struct {
	pool *pgxpool.Pool
}

func (m UserRepoImp) migrate(ctx context.Context) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS users(
			telegram_id int UNIQUE NOT NULL,
			telegram_username TEXT,
			wallex_token TEXT DEFAULT NULL,
			ramzinex_token TEXT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP DEFAULT NULL
		)
		`
	_, err := m.pool.Exec(ctx, stmt)
	return err

}

func (m UserRepoImp) GetOrCreate(ctx context.Context, user pkg.User) (pkg.User, error) {
	stmt := `
		INSERT INTO
			users(telegram_id, telegram_username, wallex_token, ramzinex_token)
			VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id)
			DO UPDATE SET telegram_id = $1, telegram_username = $2, wallex_token = $3, ramzinex_token = $4
		RETURNING
			telegram_id, COALESCE(telegram_username, ''), COALESCE(wallex_token, ''), COALESCE(ramzinex_token, ''), created_at, deleted_at
	`
	err := m.pool.QueryRow(ctx, stmt, user.TelegramID, user.TelegramUsername, user.WallexToken, user.RamzinexToken).Scan(
		&user.TelegramID,
		&user.TelegramUsername,
		&user.WallexToken,
		&user.RamzinexToken,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	return user, err
}

func (m UserRepoImp) GetAll(ctx context.Context) ([]pkg.User, error) {
	stmt := `
		SELECT telegram_id, COALESCE(telegram_username, ''), COALESCE(wallex_token, ''), COALESCE(ramzinex_token, ''), created_at
		FROM users
		WHERE deleted_at IS NULL
	`
	rows, err := m.pool.Query(ctx, stmt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()
	var users []pkg.User
	for rows.Next() {
		var user pkg.User
		err = rows.Scan(
			&user.TelegramID,
			&user.TelegramUsername,
			&user.WallexToken,
			&user.RamzinexToken,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		users = append(users, user)
	}

	return users, err
}

func NewUserRepo(ctx context.Context, pool *pgxpool.Pool) (repositories.UserRepo, error) {
	ans := UserRepoImp{
		pool: pool,
	}

	return ans, ans.migrate(ctx)
}
