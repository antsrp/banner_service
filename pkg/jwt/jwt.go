package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	method  *jwt.SigningMethodHMAC
	signKey []byte
}

type JWTOption func(s *Service)

func NewJwtService(signKey []byte, opts ...JWTOption) Service {
	s := Service{signKey: signKey}

	for _, opt := range opts {
		opt(&s)
	}

	if s.method == nil {
		WithMethodHS256(&s)
	}

	return s
}

func (js Service) Parse(token string) (map[string]any, error) {
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(jwtToken *jwt.Token) (interface{}, error) {

		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}
		return js.signKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't parse jwt token: %w", err)
	}

	//fmt.Printf("jwt: %#v", jwtToken)

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	m := make(map[string]any, len(claims))

	for k, v := range claims {
		m[k] = v
	}
	return m, nil
}

func (js Service) NewToken(fields map[string]any) (string, error) {
	fmap := make(jwt.MapClaims, len(fields))
	for k, v := range fields {
		fmap[k] = v
	}
	token := jwt.NewWithClaims(js.method, fmap)

	s, err := token.SignedString(js.signKey)
	if err != nil {
		return "", fmt.Errorf("can't sign token with key %s: %w", js.signKey, err)
	}
	return s, nil
}
