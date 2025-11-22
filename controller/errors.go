package controller

import (
	"errors"
)

var (
	ErrUnknownCommand      = errors.New("unknown command")
	ErrConstructorType     = errors.New("constructor is required and must be of ModuleConstructor type")
	ErrConfigType          = errors.New("config is required and must be of type map[string]any")
	ErrNameType            = errors.New("name is required and must be of type string")
	ErrConfigNotFound      = errors.New("config argument was not found")
	ErrConstructorNotFound = errors.New("constructor argument was not found")
	ErrNameNotFound        = errors.New("name argument was not found")
	ErrModuleAlreadyExists = errors.New("module with this name already exists")
	ErrModuleNotFound      = errors.New("a module with this was not found")
)

var (
	ErrConstructorNotFunction        = errors.New("constructor is not a function")
	ErrConstructorInvalidParams      = errors.New("constructor must have exactly 2 parameters")
	ErrConstructorInvalidReturns     = errors.New("constructor must return exactly 4 values")
	ErrConstructorInvalidReturnTypes = errors.New("constructor must return (any, string, string, error)")
	ErrConstructorLastReturnNotError = errors.New("last return value must be error")
	ErrConfigNotMap                  = errors.New("config must be map[string]any")
)
