package jwt

import "github.com/golang-jwt/jwt/v5"

func WithMethodHS256(s *Service) {
	s.method = jwt.SigningMethodHS256
}

func WithMethodHS512(s *Service) {
	s.method = jwt.SigningMethodHS512
}
