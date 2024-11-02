package storage

import "github.com/pkg/errors"

var (
	ErrProfileDoesNotExist = errors.New("profile does not exist")
	ErrProfileAlreadyExist = errors.New("profile already exist")
)
