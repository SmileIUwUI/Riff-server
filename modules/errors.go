package modules

import "errors"

var (
	ErrDBNotFound  = errors.New("db argument was not found")
	ErrDBType      = errors.New("db is required and must be of *sqlx.DB type")
	ErrConfigIsNil = errors.New("config is nil")
)
