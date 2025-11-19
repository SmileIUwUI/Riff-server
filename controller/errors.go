package controller

import (
	"errors"
)

var (
	ErrUnknownCommand      = errors.New("unknown command")
	ErrConstructorType     = errors.New("constructor is required and must be ModuleConstructor type")
	ErrConfigNotFound      = errors.New("argumet config was not found")
	ErrConstructorNotFound = errors.New("constructor config was not found")
	ErrModuleAlreadyExists = errors.New("a module with that name already exists")
)
