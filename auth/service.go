package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(user_id int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
type serviceJwt struct {
}

func NewService() *serviceJwt {
	return &serviceJwt{}
}
func (s *serviceJwt) GenerateToken(id int) (string, error) {
	claim := jwt.MapClaims{
		"id": id,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte("bwa_secret_key_s4j4"))

	if err != nil {
		return token, err
	}
	return token, nil
}

func (s *serviceJwt) ValidateToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte("bwa_secret_key_s4j4"), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
