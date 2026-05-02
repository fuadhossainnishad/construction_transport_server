package utils

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordHasher struct{}

func (BcryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (BcryptPasswordHasher) Compare(hash_password, raw_password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash_password), []byte(raw_password)) == nil
}
