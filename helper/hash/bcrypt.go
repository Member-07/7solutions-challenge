package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, password string) error
}

type bcryptHasher struct {
	cost int
}

func NewBcrypt() BcryptHasher {
	return &bcryptHasher{cost: bcrypt.DefaultCost}
}

func (b *bcryptHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (b *bcryptHasher) Compare(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
