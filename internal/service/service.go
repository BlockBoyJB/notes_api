package service

import (
	"context"
	"notes_api/internal/repo"
	"notes_api/pkg/hasher"
	"time"
)

type (
	UserInput struct {
		Username string
		Password string
	}
	NoteInput struct {
		Username string
		Title    string
		Text     string
	}
	NoteListInput struct {
		Username string
		Sort     string
		Offset   int
		Limit    int
	}
	NoteOutput struct {
		Id        int       `json:"id"`
		Title     string    `json:"title"`
		Text      string    `json:"text"`
		CreatedAt time.Time `json:"created_at"`
	}
	TextError struct {
		Type         string   `json:"type"`
		Position     int      `json:"position"`
		Row          int      `json:"row"`
		Column       int      `json:"column"`
		Length       int      `json:"length"`
		Word         string   `json:"word"`
		Replacements []string `json:"replacements"`
	}
)

type Auth interface {
	CreateToken(username string) (string, error)
	ParseToken(tokenString string) (*TokenClaims, error)
}

type User interface {
	Create(ctx context.Context, input UserInput) error
	VerifyPassword(ctx context.Context, input UserInput) (bool, error)
}

type Note interface {
	Create(ctx context.Context, input NoteInput) (int, error)
	Find(ctx context.Context, input NoteListInput) ([]NoteOutput, error)
	ValidateText(text string) ([]TextError, error)
}

type (
	Services struct {
		Auth Auth
		User User
		Note Note
	}
	ServicesDependencies struct {
		Repos      *repo.Repositories
		Hasher     hasher.Hasher
		TokenTTL   time.Duration
		PrivateKey string
		PublicKey  string
	}
)

func NewServices(d *ServicesDependencies) *Services {
	return &Services{
		Auth: newAuthService(d.PrivateKey, d.PublicKey, d.TokenTTL),
		User: newUserService(d.Repos.User, d.Hasher),
		Note: newNoteService(d.Repos.Note),
	}
}
