package errors

import (
	"errors"
)

var UserAlreadyExistError = errors.New("could not create user: Username already exists")

var UserNotFoundError = errors.New("user not found")

var PasswordIncorrectError = errors.New("password is incorrect")

var HandleNotFoundError = errors.New("handle not available")

var ErrJobQueueFull = errors.New("Job queue completely full")

var FieldEmptyError = errors.New("empty field forbidden")

var UserUnverifiedError = errors.New("E-mail not verified")