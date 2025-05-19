package store

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRoleStatus int

type UserRole struct {
	ID       DefaultRoleID   `json:"id"`
	Name     DefaultRoleName `json:"name"`
	IsActive bool            `json:"isActive"`
}
type UserFilterQuery struct {
	IsVerified bool   `json:"isVerified"`
	Search     string `json:"search"`
}
type User struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Password   password   `json:"-"`
	IsVerified bool       `json:"isVerified"`
	VerifiedAt *string    `json:"verifiedAt"`
	Roles      []UserRole `json:"roles"`
	Common
}

type password struct {
	Hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Hash = hash
	return nil

}
func (p *password) Compare(text string) bool {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(text))
	return err == nil

}
func (p *password) GetHash() []byte {
	return p.Hash

}

var (
	ErrDuplicateEmail    = errors.New("email has been taken")
	ErrDuplicateUserRole = errors.New("user already has this role")
	ErrVerifyUser        = errors.New("issue verifying user")
)
