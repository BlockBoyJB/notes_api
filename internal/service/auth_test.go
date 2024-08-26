package service

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuthService_CreateToken(t *testing.T) {
	auth := newAuthService(
		"../../private.key",
		"../../public.key",
		time.Hour,
	)

	testCases := []struct {
		testName string
		username string
	}{
		{
			testName: "Correct test",
			username: "vasya",
		},
		{
			testName: "Token without username",
			username: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			token, err := auth.CreateToken(tc.username)

			assert.Equal(t, nil, err)
			assert.NotEqual(t, "", token)
		})
	}
}

func TestAuthService_ParseToken(t *testing.T) {
	auth := newAuthService(
		"../../private.key",
		"../../public.key",
		time.Hour*72,
	)
	token, err := auth.CreateToken("vasya")
	if err != nil {
		panic(err)
	}

	et := jwt.NewWithClaims(defaultSignMethod, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() - 60,
			IssuedAt:  time.Now().Unix() - 2000,
		},
		Username: "vasya",
	})
	expiredToken, err := et.SignedString(auth.privateKey)
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		testName  string
		username  string
		token     string
		expectErr error
	}{
		{
			testName:  "Correct test",
			username:  "vasya",
			token:     token,
			expectErr: nil,
		},
		{
			testName:  "Not a token",
			token:     "fooooobar",
			expectErr: jwt.NewValidationError("token contains an invalid number of segments", jwt.ValidationErrorMalformed),
		},
		{
			testName: "Incorrect token sign",
			token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMjQ1OTAxMTUsImlhdCI6MTcyNDU5MDExNSwidXNlcm5hbWUiOiJ2YXN5YSJ9.lLqA8eUk3ad7qnsExUay_ukE7Du_VxAGUtzKpEIyXdgBosQiH6dAeW0CFr1cjJaqU7hO3vpS-ZTOdgO1_F1R2NoTgWFIplYzf7GFbQ4P1AWKp7QcNU2pJvpfEyoaghrbGGKw1rKrWzDZX2TwY79LspNG31D1ioqDXahsxQBHBsj6gXB8Q0bvSrI-CbrC1NBlYJndLBkMfDPFS7AhedlLNz8DqmPadpCbTs-KUqfBJvSDev8RLK2X8bbxJ3V0gwZCfnBTq0PkO2wQdFLI047MpSg2E0401jd7FXPKj6KiuRk4vk9-k3ceoSsdWqqAtoEHldvCa2GWHxLSo1SqC3yzjg",
			expectErr: &jwt.ValidationError{
				Inner:  rsa.ErrVerification,
				Errors: jwt.ValidationErrorSignatureInvalid,
			},
		},
		{
			testName: "Incorrect token sign method",
			token:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1NDEwMjMsImlhdCI6MTcyNDUzMzgyMywidXNlcm5hbWUiOiJ2YXN5YSJ9.ZZ-rwDr3SOWvmmNl9jUS1rmdNj5PARa2DTf3Yr8_kyVvEewMCxXeB42KNl1ulYOkZiCA7Xllp49GqRT_wh4b4w",
			expectErr: &jwt.ValidationError{
				Inner:  ErrIncorrectSignMethod,
				Errors: jwt.ValidationErrorUnverifiable,
			},
		},
		{
			testName: "Expired token",
			token:    expiredToken,
			expectErr: &jwt.ValidationError{
				Inner:  errors.New("token is expired by 1m0s"),
				Errors: jwt.ValidationErrorExpired,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			claims, err := auth.ParseToken(tc.token)

			assert.Equal(t, tc.expectErr, err)

			if tc.expectErr == nil {
				assert.Equal(t, tc.username, claims.Username)
			}
		})
	}
}
