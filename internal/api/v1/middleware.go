package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"notes_api/internal/service"
	"os"
	"strings"
)

const (
	bearerPrefix = "Bearer "
	usernameCtx  = "username"
)

type authMiddleware struct {
	auth service.Auth
}

func (r *authMiddleware) authHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := parseToken(c.Request())
		if !ok {
			errorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader)
			return nil
		}
		claims, err := r.auth.ParseToken(token)
		if err != nil {
			if errors.Is(err, service.ErrCannotParseToken) {
				errorResponse(c, http.StatusUnauthorized, err)
				return nil
			} else {
				errorResponse(c, http.StatusForbidden, err)
				return nil
			}
		}
		c.Set(usernameCtx, claims.Username)
		return next(c)
	}
}

func parseToken(r *http.Request) (string, bool) {
	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}
	token := strings.Split(header, bearerPrefix)
	if len(token) != 2 {
		return "", false
	}
	return token[1], true
}

func LoggingMiddleware(h *echo.Echo, output string) {
	cfg := middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}", "method":"${method}","uri":"${uri}", "status":${status}, "error":"${error}"}` + "\n",
	}
	if output == "stdout" {
		cfg.Output = os.Stdout
	} else {
		file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Output = file
	}
	h.Use(middleware.LoggerWithConfig(cfg))
}
