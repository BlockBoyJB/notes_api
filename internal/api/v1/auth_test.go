package v1

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"notes_api/internal/service"
)

func (s *APITestSuite) TestAuthRouter_signUp() {
	testCases := []struct {
		testName   string
		inputBody  string
		username   string
		expectCode int
	}{
		{
			testName:   "Correct test",
			username:   "vasya",
			inputBody:  `{"username": "vasya", "password": "abc"}`,
			expectCode: 201,
		},
		{
			testName:   "Repeated username test",
			inputBody:  `{"username": "vasya", "password": "abc"}`,
			expectCode: 400,
		},
		{
			testName:   "Incorrect username: too short",
			inputBody:  `{"username": "o", "password": "abc"}`,
			expectCode: 400,
		},
		{
			testName:   "Incorrect username: too long",
			inputBody:  `{"username": "tooloooooooooooooooooooooooooooooong", "password": "abc"}`,
			expectCode: 400,
		},
		{
			testName:   "Incorrect username: invalid characters",
			inputBody:  `{"username": "va$$$ya", "password": "abc"}`,
			expectCode: 400,
		},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.inputBody))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		s.router.ServeHTTP(w, r)

		s.Assert().Equal(tc.expectCode, w.Code)
		if tc.expectCode == 201 {
			u, err := s.repos.User.FindByUsername(s.ctx, tc.username)
			if err != nil {
				panic(err)
			}
			s.Assert().Equal(tc.username, u.Username)
		}
	}
}

func (s *APITestSuite) TestAuthRouter_signIn() {
	// Предварительно создаем пользователя
	if err := s.services.User.Create(s.ctx, service.UserInput{
		Username: "vasya",
		Password: "abc",
	}); err != nil {
		panic(err)
	}

	testCases := []struct {
		testName   string
		inputBody  string
		expectCode int
	}{
		{
			testName:   "Correct test",
			inputBody:  `{"username": "vasya", "password": "abc"}`,
			expectCode: 200,
		},
		{
			testName:   "Wrong password",
			inputBody:  `{"username": "vasya", "password": "foobar"}`,
			expectCode: 403,
		},
		{
			testName:   "User not exists",
			inputBody:  `{"username": "petya", "password": "abc"}`,
			expectCode: 400,
		},
		{
			testName:   "Incorrect username input",
			inputBody:  `{"username": "va$$$ya", "password": "abc"}`,
			expectCode: 400,
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(tc.inputBody))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		s.router.ServeHTTP(w, r)

		s.Assert().Equal(tc.expectCode, w.Code, tc.testName)
		if tc.expectCode == 200 {
			var response signInResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			s.Assert().NotEqual("", response.Token)
		}
	}
}
