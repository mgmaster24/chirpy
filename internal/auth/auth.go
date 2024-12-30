package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HsshPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(pw, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	registeredClaims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&registeredClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	subject, err := registeredClaims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func GetBearerToken(header http.Header) (string, error) {
	return getValue(header, "Bearer")
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	val := hex.EncodeToString(b)
	return val, nil
}

func Authenticate(secret string, header http.Header) (uuid.UUID, error) {
	token, err := GetBearerToken(header)
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := ValidateJWT(token, secret)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func GetAPIKeys(header http.Header) (string, error) {
	return getValue(header, "ApiKey")
}

func getValue(header http.Header, authType string) (string, error) {
	authHeader := header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("Authorization header is not present")
	}
	val := strings.TrimPrefix(authHeader, authType)
	val = strings.TrimSpace(val)
	return val, nil
}
