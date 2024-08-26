package service

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"notes_api/internal/model/dbmodel"
	"notes_api/internal/repo"
	"notes_api/internal/repo/pgerrs"
	"notes_api/pkg/hasher"
)

const (
	userServicePrefixLog = "/service/user"
)

type userService struct {
	user   repo.User
	hasher hasher.Hasher
}

func newUserService(user repo.User, hasher hasher.Hasher) *userService {
	return &userService{
		user:   user,
		hasher: hasher,
	}
}

func (s *userService) Create(ctx context.Context, input UserInput) error {
	err := s.user.Create(ctx, dbmodel.User{
		Username: input.Username,
		Password: s.hasher.Hash(input.Password),
	})

	if err != nil {
		if errors.Is(err, pgerrs.ErrAlreadyExists) {
			return ErrUserAlreadyExists
		}
		log.Errorf("%s/Create error create user: %s", userServicePrefixLog, err)
		return err
	}
	return nil
}

func (s *userService) VerifyPassword(ctx context.Context, input UserInput) (bool, error) {
	u, err := s.user.FindByUsername(ctx, input.Username)

	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return false, ErrUserNotFound
		}
		log.Errorf("%s/VerifyPassword error find user: %s", userServicePrefixLog, err)
		return false, err
	}
	return s.hasher.Verify(input.Password, u.Password), nil
}
