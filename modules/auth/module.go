package auth

import (
	"Riff/controller"
	"Riff/modules"

	"github.com/jmoiron/sqlx"
)

const (
	NAME    = "auth"
	VERSION = "0.1.0"
)

type AuthModule struct {
	db       *sqlx.DB
	receiver chan controller.Command
}

func NewAuthModule(config map[string]any, receiver chan controller.Command) (any, string, string, error) {
	db, err := modules.GetArg[*sqlx.DB](config, "db")

	if err != nil {
		return nil, NAME, VERSION, err
	}

	module := &AuthModule{
		db:       db,
		receiver: receiver,
	}

	return module, NAME, VERSION, nil
}
