package services

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/Graynie/InkZen/internal/models"
)

type UsuarioService struct{}

func NewUsuarioService() *UsuarioService {
	return &UsuarioService{}
}

func (s *UsuarioService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *UsuarioService) PrepareUser(user models.Usuario) (models.Usuario, error) {
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return user, err
	}

	user.Password = hashedPassword
	return user, nil
}

func (s *UsuarioService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
