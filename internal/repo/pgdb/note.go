package pgdb

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"notes_api/internal/model/dbmodel"
	"notes_api/internal/repo/pgerrs"
	"notes_api/pkg/postgres"
)

type NoteRepo struct {
	*postgres.Postgres
}

func NewNoteRepo(pg *postgres.Postgres) *NoteRepo {
	return &NoteRepo{pg}
}

func (r *NoteRepo) Create(ctx context.Context, n dbmodel.Note) (int, error) {
	sql, args, _ := r.Builder.
		Insert("note").
		Columns("username", "title", "text").
		Values(n.Username, n.Title, n.Text).
		Suffix("returning id").
		ToSql()

	var noteId int
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&noteId); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23503" {
				return 0, pgerrs.ErrForeignKey
			}
		}
		return 0, err
	}
	return noteId, nil
}

func (r *NoteRepo) Find(ctx context.Context, username, sort string, offset, limit int) ([]dbmodel.Note, error) {
	sql, args, _ := r.Builder.
		Select("id, title, text, created_at").
		From("note").
		Where("username = ?", username).
		OrderBy(sort).Offset(uint64(offset)).
		Limit(uint64(limit)).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dbmodel.Note
	for rows.Next() {
		var note dbmodel.Note

		err = rows.Scan(
			&note.Id,
			&note.Title,
			&note.Text,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, note)
	}
	return result, nil
}
