package controller

import (
	"time"

	"github.com/google/uuid"
)

type Controller struct {
	modules          map[string]Module
	modulesReceivers []chan Command
	queueCommands    []Command
	externalReceiver chan Command
	commandHandlers  map[CommandType]CommandHandler
}

func NewController() *Controller {
	controller := &Controller{
		modules:          make(map[string]Module),
		modulesReceivers: make([]chan Command, 0),
		queueCommands:    make([]Command, 0),
		externalReceiver: make(chan Command, 256),
		commandHandlers:  make(map[CommandType]CommandHandler),
	}

	controller.registerCommandHandlers()

	go controller.commandLoop()

	return controller
}

func (c *Controller) SubmitCommand(commandType CommandType, data map[string]any) <-chan CommandResult {
	cmd := Command{
		ID:       uuid.New(),
		Type:     commandType,
		Source:   CommandSourceExternal,
		Sender:   "external",
		Data:     data,
		Chanback: make(chan CommandResult, 1),
	}

	c.externalReceiver <- cmd
	return cmd.Chanback
}

func (c *Controller) registerCommandHandlers() {
	c.commandHandlers[CommandAddModule] = c.handlerAddModule
	c.commandHandlers[CommandRemoveModule] = c.handlerRemoveModule
	c.commandHandlers[CommandListModules] = c.handlerListModules
}

func (c *Controller) commandLoop() {
	for {
		c.processExternalCommands()

		c.processInternalCommands()

		c.processExecuteCommands()

		time.Sleep(5 * time.Millisecond)
	}
}

func (c *Controller) processExternalCommands() {
	for {
		select {
		case command, ok := <-c.externalReceiver:
			if !ok {
				return
			}
			c.queueCommands = append(c.queueCommands, command)
		default:
			return
		}
	}
}

func (c *Controller) processInternalCommands() {
	for _, receiver := range c.modulesReceivers {
	ReceiverLoop:
		for {
			select {
			case command, ok := <-receiver:
				if !ok {
					break ReceiverLoop
				}
				c.queueCommands = append(c.queueCommands, command)
			default:
				break ReceiverLoop
			}
		}
	}
}

func (c *Controller) processExecuteCommands() {
	if len(c.queueCommands) == 0 {
		return
	}

	for _, cmd := range c.queueCommands {
		c.executeCommand(cmd)
	}

	c.queueCommands = make([]Command, 0)
}

func (c *Controller) executeCommand(cmd Command) {
	handler, exists := c.commandHandlers[cmd.Type]

	if !exists {
		result := CommandResult{
			Command:   cmd,
			Error:     ErrUnknownCommand,
			Timestamp: time.Now(),
		}

		cmd.Chanback <- result
		close(cmd.Chanback)
		return
	}

	handler(cmd)
}
