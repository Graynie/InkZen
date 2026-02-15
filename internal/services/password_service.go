package services

import "golang.org/x/crypto/bcrypt"

// Hashea una contraseña antes de guardarla
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Verifica si la contraseña ingresada coincide con la guardada
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
