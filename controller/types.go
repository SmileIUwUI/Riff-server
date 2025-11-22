package controller

import (
	"time"

	"github.com/google/uuid"
)

type CommandType string

type CommandSource string

type CommandHandler func(cmd Command)

type CommandResult struct {
	Command   Command
	Result    any
	Error     error
	Timestamp time.Time
}

type Command struct {
	ID       uuid.UUID
	Type     CommandType
	Source   CommandSource
	Sender   string
	Executor string // moduleName.method
	Data     map[string]any
	Chanback chan CommandResult
}

type Module struct {
	Name     string
	Instance any
	Version  string
	Config   map[string]any
}
