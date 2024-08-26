package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"notes_api/internal/service"
)

type noteRouter struct {
	note service.Note
}

func newNoteRouter(g *echo.Group, note service.Note) {
	r := &noteRouter{note: note}

	g.POST("/create", r.create)
	g.GET("/list", r.list)
	g.POST("/validate", r.validate)
}

type noteCreateInput struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"text" validate:"required"`
}

type noteCreateResponse struct {
	NoteId int `json:"note_id"`
}

//	@Summary		Create note
//	@Description	Create note
//	@Tags			note
//	@Accept			json
//	@Produce		json
//	@Param			input	body		noteCreateInput	true	"input"
//	@Success		200		{object}	noteCreateResponse
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/notes/create [post]
func (r *noteRouter) create(c echo.Context) error {
	var input noteCreateInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	username, ok := c.Get(usernameCtx).(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return nil
	}
	noteId, err := r.note.Create(c.Request().Context(), service.NoteInput{
		Username: username,
		Title:    input.Title,
		Text:     input.Text,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.JSON(http.StatusCreated, noteCreateResponse{NoteId: noteId})
}

type noteListInput struct {
	Sort   string `json:"sort"`
	Offset int    `json:"offset" validate:"value"`
	Limit  int    `json:"limit" validate:"limit"`
}

//	@Summary		Get list notes
//	@Description	Get user list notes
//	@Tags			note
//	@Accept			json
//	@Produce		json
//	@Param			input	body		noteListInput	true	"input"
//	@Success		200		{array}		service.NoteOutput
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/notes/list [get]
func (r *noteRouter) list(c echo.Context) error {
	var input noteListInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}

	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	username, ok := c.Get(usernameCtx).(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return nil
	}
	notes, err := r.note.Find(c.Request().Context(), service.NoteListInput{
		Username: username,
		Sort:     input.Sort,
		Offset:   input.Offset,
		Limit:    input.Limit,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.JSON(http.StatusOK, notes)
}

type noteValidateInput struct {
	Text string `json:"text" validate:"required"`
}

//	@Summary		Validate note
//	@Description	Validate spelling mistakes
//	@Tags			note
//	@Accept			json
//	@Produce		json
//	@Param			input	body		noteValidateInput	true	"input"
//	@Success		200		{array}		service.TextError
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/notes/validate [post]
func (r *noteRouter) validate(c echo.Context) error {
	var input noteValidateInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	textErrs, err := r.note.ValidateText(input.Text)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.JSON(http.StatusOK, textErrs)
}
