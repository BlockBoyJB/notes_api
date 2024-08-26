package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"notes_api/internal/service"
)

type authRouter struct {
	auth service.Auth
	user service.User
}

func newAuthRouter(g *echo.Group, auth service.Auth, user service.User) {
	r := &authRouter{
		auth: auth,
		user: user,
	}
	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
}

type signUpInput struct {
	Username string `json:"username" validate:"username,required"`
	Password string `json:"password" validate:"required"`
}

//	@Summary		Sign up
//	@Description	Create account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body	signUpInput	true	"input"
//	@Success		201
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/auth/sign-up [post]
func (r *authRouter) signUp(c echo.Context) error {
	var input signUpInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	err := r.user.Create(c.Request().Context(), service.UserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

type signInInput struct {
	Username string `json:"username" validate:"username,required"`
	Password string `json:"password" validate:"required"`
}

type signInResponse struct {
	Token string `json:"token"`
}

//	@Summary		Sign in
//	@Description	Sign in to account for getting token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		signInInput	true	"input"
//	@Success		200		{object}	signInResponse
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/auth/sign-in [post]
func (r *authRouter) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	ok, err := r.user.VerifyPassword(c.Request().Context(), service.UserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	if !ok {
		errorResponse(c, http.StatusForbidden, echo.ErrForbidden)
		return nil
	}

	token, err := r.auth.CreateToken(input.Username)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.JSON(http.StatusOK, signInResponse{Token: token})
}
