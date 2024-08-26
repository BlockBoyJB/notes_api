package v1

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"notes_api/internal/model/dbmodel"
	"notes_api/internal/service"
)

type noteTestsData struct {
	username string
	password string
	token    string
}

func setupNoteTestsData(s *APITestSuite) *noteTestsData {
	username, password := "vasya", "abc"
	err := s.services.User.Create(s.ctx, service.UserInput{
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	token, err := s.services.Auth.CreateToken(username)
	if err != nil {
		panic(err)
	}
	return &noteTestsData{
		username: username,
		password: password,
		token:    token,
	}
}

func (s *APITestSuite) TestNoteRouter_create() {
	setup := setupNoteTestsData(s) // default user is vasya with abc password

	fakeUser, err := s.services.Auth.CreateToken("petya")
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		testName   string
		inputBody  string
		token      string
		expectCode int
	}{
		{
			testName:   "Correct test",
			inputBody:  `{"title": "Foobar", "text": "some text"}`,
			token:      setup.token,
			expectCode: 201,
		},
		{
			testName:   "User not exists",
			inputBody:  `{"title": "Foobar", "text": "some text"}`,
			token:      fakeUser,
			expectCode: 400,
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/notes/create", bytes.NewBufferString(tc.inputBody))

		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		r.Header.Set(echo.HeaderAuthorization, "Bearer "+tc.token)

		s.router.ServeHTTP(w, r)

		s.Assert().Equal(tc.expectCode, w.Code, tc.testName)
	}
}

func (s *APITestSuite) TestNoteRouter_list() {
	setup := setupNoteTestsData(s)

	note := dbmodel.Note{
		Username: setup.username,
		Title:    "Foobar",
		Text:     "SomeText",
	}

	var err error
	note.Id, err = s.repos.Note.Create(s.ctx, note)
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		testName   string
		inputBody  string
		expectCode int
	}{
		{
			testName:   "Correct test",
			inputBody:  `{"sort": "id", "limit": 20}`,
			expectCode: 200,
		},
		{
			testName:   "Incorrect offset",
			inputBody:  `{"sort": "id", "offset": -1, "limit": 20}`,
			expectCode: 400,
		},
		{
			testName:   "Incorrect limit (too high)",
			inputBody:  `{"sort": "id", "limit": 100000000}`,
			expectCode: 400,
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/notes/list", bytes.NewBufferString(tc.inputBody))

		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		r.Header.Set(echo.HeaderAuthorization, "Bearer "+setup.token)

		s.router.ServeHTTP(w, r)

		s.Assert().Equal(tc.expectCode, w.Code, tc.testName)

		if tc.expectCode == 200 {
			var response []service.NoteOutput
			if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			s.Assert().Equal(note.Id, response[0].Id)
			s.Assert().Equal(note.Title, response[0].Title)
			s.Assert().Equal(note.Text, response[0].Text)
		}
	}
}
