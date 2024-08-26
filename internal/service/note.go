package service

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"notes_api/internal/model/dbmodel"
	"notes_api/internal/repo"
	"notes_api/internal/repo/pgerrs"
	"notes_api/pkg/validator"
)

const (
	noteServicePrefixLog = "/service/auth"
)

type noteService struct {
	note repo.Note
}

func newNoteService(note repo.Note) *noteService {
	return &noteService{
		note: note,
	}
}

func (s *noteService) Create(ctx context.Context, input NoteInput) (int, error) {
	noteId, err := s.note.Create(ctx, dbmodel.Note{
		Username: input.Username,
		Title:    input.Title,
		Text:     input.Text,
	})
	if err != nil {
		if errors.Is(err, pgerrs.ErrForeignKey) {
			return 0, ErrUserNotFound
		}
		log.Errorf("%s/Create error create note: %s", noteServicePrefixLog, err)
		return 0, err
	}

	return noteId, nil
}

func (s *noteService) Find(ctx context.Context, input NoteListInput) ([]NoteOutput, error) {
	// глупенькая сортировка
	switch input.Sort {
	case "id":
		input.Sort = "id ASC"
	case "date":
		input.Sort = "created_at DESC"
	default:
		input.Sort = "id DESC"
	}

	notes, err := s.note.Find(ctx, input.Username, input.Sort, input.Offset, input.Limit)
	if err != nil {
		log.Errorf("%s/Find error find user notes: %s", noteServicePrefixLog, err)
		return nil, err
	}

	result := make([]NoteOutput, 0)
	for _, note := range notes {
		result = append(result, NoteOutput{
			Id:        note.Id,
			Title:     note.Title,
			Text:      note.Text,
			CreatedAt: note.CreatedAt,
		})
	}
	return result, nil
}

func (s *noteService) ValidateText(text string) ([]TextError, error) {
	spellResult, err := validator.CheckText(text)
	if err != nil {
		log.Errorf("%s/Validate error validate text: %s", noteServicePrefixLog, err)
		return nil, err
	}

	result := make([]TextError, 0)

	for _, e := range spellResult {
		result = append(result, TextError{
			Type:         validator.Errors[e.Code],
			Position:     e.Pos,
			Row:          e.Row,
			Column:       e.Col,
			Length:       e.Len,
			Word:         e.Word,
			Replacements: e.S,
		})
	}
	return result, nil
}
