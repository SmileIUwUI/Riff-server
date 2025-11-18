package controller

import "time"

type Controller struct {
	Modules          map[string]Module
	ModulesReceivers []chan Command
	QueueCommands    []Command
	ExternalReceiver chan Command
}

func NewController() *Controller {
	controller := &Controller{
		Modules:          make(map[string]Module),
		ModulesReceivers: make([]chan Command, 0),
		QueueCommands:    make([]Command, 0),
		ExternalReceiver: make(chan Command),
	}

	go controller.commandLoop()

	return controller
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
		for {
			select {
			case command, ok := <-receiver:
				if !ok {
					break
				}

				c.QueueCommands = append(c.QueueCommands, command)
			default:
				break
			}
			break
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

}
