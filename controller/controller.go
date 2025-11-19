package controller

import "time"

type Controller struct {
	Modules          map[string]Module
	ModulesReceivers []chan Command
	QueueCommands    []Command
	ExternalReceiver chan Command
	commandHandlers  map[CommandType]CommandHandler
}

func NewController() *Controller {
	controller := &Controller{
		Modules:          make(map[string]Module),
		ModulesReceivers: make([]chan Command, 0),
		QueueCommands:    make([]Command, 0),
		ExternalReceiver: make(chan Command),
	}

	controller.registerCommandHandlers()

	go controller.commandLoop()

	return controller
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
		case command, ok := <-c.ExternalReceiver:
			if !ok {
				return
			}
			c.QueueCommands = append(c.QueueCommands, command)
		default:
			return
		}
	}
}

func (c *Controller) processInternalCommands() {
	for _, receiver := range c.ModulesReceivers {
	ReceiverLoop:
		for {
			select {
			case command, ok := <-receiver:
				if !ok {
					break ReceiverLoop
				}
				c.QueueCommands = append(c.QueueCommands, command)
			default:
				break ReceiverLoop
			}
		}
	}
}

func (c *Controller) processExecuteCommands() {
	if len(c.QueueCommands) == 0 {
		return
	}

	for _, cmd := range c.QueueCommands {
		c.executeCommand(cmd)
	}

	c.QueueCommands = make([]Command, 0)
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
