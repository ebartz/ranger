package users

import (
	"github.com/ranger/ranger/tests/framework/pkg/namegenerator"
)

const (
	defaultPasswordLength = 12
)

func GenerateUserPassword(password string) string {
	return namegenerator.RandStringLower(defaultPasswordLength)
}
