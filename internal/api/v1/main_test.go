package v1

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"notes_api/internal/repo"
	"notes_api/internal/service"
	"notes_api/pkg/hasher"
	"notes_api/pkg/postgres"
	"notes_api/pkg/validator"
	"testing"
	"time"
)

type APITestSuite struct {
	suite.Suite
	ctx      context.Context
	router   *echo.Echo
	pg       *postgres.Postgres
	repos    *repo.Repositories
	services *service.Services
	m        *migrate.Migrate
}

func (s *APITestSuite) SetupTest() {
	testPGUrl := "postgres://postgres:1234567890@localhost:6000/postgres"
	m, err := migrate.New("file://../../../migrations", testPGUrl+"?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	s.m = m

	s.ctx = context.Background()

	pg, err := postgres.NewPG(testPGUrl)
	if err != nil {
		panic(err)
	}
	s.pg = pg

	s.repos = repo.NewRepositories(pg)
	d := &service.ServicesDependencies{
		Repos:      s.repos,
		Hasher:     hasher.NewHasher("secret"),
		TokenTTL:   time.Hour,
		PrivateKey: "../../../private.key",
		PublicKey:  "../../../public.key",
	}
	s.services = service.NewServices(d)

	s.router = echo.New()

	s.router.Validator, err = validator.NewValidator()
	if err != nil {
		panic(err)
	}

	NewRouter(s.router, s.services)
}

func (s *APITestSuite) TearDownTest() {
	_ = s.m.Drop()
	s.pg.Close()
}

func TestAllRoutes(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
