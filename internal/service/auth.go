package service

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	authServicePrefixLog = "/service/auth"
)

var defaultSignMethod = jwt.SigningMethodRS256

type TokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type authService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	tokenTTL   time.Duration
}

func newAuthService(privateKey, publicKey string, tokenTTL time.Duration) *authService {
	privateKeyData, err := os.ReadFile(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKeyData, err := os.ReadFile(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	private, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatal(err)
	}
	public, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		log.Fatal(err)
	}

	return &authService{
		privateKey: private,
		publicKey:  public,
		tokenTTL:   tokenTTL,
	}
}

func (s *authService) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(defaultSignMethod, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: username,
	})

	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		log.Errorf("%s/CreateToken error sign token: %s", authServicePrefixLog, err)
		return "", err
	}
	return signedToken, nil
}

func (s *authService) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrIncorrectSignMethod
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrCannotParseToken
	}
	return claims, nil
}
