package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	TokenMalformed = errors.New("token is malformed")
	TokenExpired   = errors.New("token is expired")
	TokenInvalid   = errors.New("token is invalid")
)

func NewToken(userId string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		Subject:   userId,
	}).SignedString([]byte("secret"))
}

func ValidateToken(t string) error {
	_, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte("secret"), nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return TokenMalformed
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return TokenExpired
			} else {
				return TokenInvalid
			}
		}
	}

	return nil
}
