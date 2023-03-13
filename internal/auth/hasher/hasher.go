package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	GetHashedPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

func NewHasher() Hasher {
	return &hasher{}
}

type hasher struct{}

func (h *hasher) GetHashedPassword(password string) (string, error) {
	pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pswd), err
}

func (h *hasher) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
