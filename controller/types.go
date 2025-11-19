package controller

import (
	"time"

	"github.com/google/uuid"
)

type CommandType string

type CommandSource string

type CommandHandler func(cmd Command)

type ModuleConstructor func(config map[string]any, receiver chan Command) (any, string, string, error) // Instance Name Version Error

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
