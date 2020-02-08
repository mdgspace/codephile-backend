package errors

import (
	"errors"
)

var UserAlreadyExistError = errors.New("could not create user: Username already exists")

var UserNotFoundError = errors.New("user not found")

var HandleNotFoundError = errors.New("handle not available")
