package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type AuthService interface {
	GenerateJWT(username string) (string, error)
	ValidateJWT(token string) (bool, error)
	SetupRSAKeys()
}

type Claims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

type Auth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (a *Auth) SetupRSAKeys() {
	var err error
	a.privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key")
	}
	a.publicKey = &a.privateKey.PublicKey
}

func (a *Auth) GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24) // Set expiration time (e.g., 24 hours)

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *Auth) ValidateJWT(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return a.publicKey, nil
	})
	if err != nil {
		return false, err
	}
	if token.Valid {
		return true, nil
	}
	return false, nil
}
