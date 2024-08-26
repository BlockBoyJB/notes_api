package repo

import (
	"context"
	"notes_api/internal/model/dbmodel"
	"notes_api/internal/repo/pgdb"
	"notes_api/pkg/postgres"
)

type User interface {
	Create(ctx context.Context, u dbmodel.User) error
	FindByUsername(ctx context.Context, username string) (dbmodel.User, error)
}

type Note interface {
	Create(ctx context.Context, n dbmodel.Note) (int, error)
	Find(ctx context.Context, username, sort string, offset, limit int) ([]dbmodel.Note, error)
}

type Repositories struct {
	User
	Note
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(pg),
		Note: pgdb.NewNoteRepo(pg),
	}
}
