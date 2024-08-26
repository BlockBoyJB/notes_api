package pgdb

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"notes_api/internal/model/dbmodel"
	"notes_api/internal/repo/pgerrs"
	"notes_api/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Create(ctx context.Context, u dbmodel.User) error {
	sql, args, _ := r.Builder.
		Insert("\"user\"").
		Columns("username", "password").
		Values(u.Username, u.Password).
		ToSql()

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return pgerrs.ErrAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (dbmodel.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("\"user\"").
		Where("username = ?", username).
		ToSql()

	var user dbmodel.User

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbmodel.User{}, pgerrs.ErrNotFound
		}
		return dbmodel.User{}, err
	}
	return user, nil
}
