package security

import (
	"errors"
	"net/http"
	"pet-market/internal/logger"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Authorization struct {
	Log *logger.Logger
}

const SecretKey = "secret-key"
const TokenExp = time.Hour * 1

type Claims struct {
	jwt.RegisteredClaims
	Login  string
	UserID int
}

func New(logger *logger.Logger) Authorization {
	return Authorization{
		Log: logger,
	}
}

func (auth *Authorization) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token, err := GetToken(token)
		if err != nil || !auth.verifyToken(token) {
			w.WriteHeader(http.StatusUnauthorized)
			return

		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func (auth *Authorization) Authorize(_ http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	token, err := GetToken(token)
	if err != nil || !auth.verifyToken(token) {
		return false

	} else {
		return true
	}
}

func (auth *Authorization) BuildJWTString(login string, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Login:  login,
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetLogin(tokenString string) string {
	var claims Claims
	jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	return claims.Login
}
func GetUserID(tokenString string) int {
	var claims Claims
	jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	return claims.UserID
}

func (auth *Authorization) verifyToken(tokenString string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		auth.Log.Log.Error("Error during verifying token", zap.String("error", err.Error()))
		return false
	}
	return token.Valid
}
func GetToken(str string) (string, error) {
	if str == "" {
		return "", errors.New("token is invalid")
	}
	tokens := strings.Split(str, " ")
	if len(tokens) != 2 {
		return "", errors.New("token is invalid")
	}
	return tokens[1], nil
}

func (auth *Authorization) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (auth *Authorization) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
