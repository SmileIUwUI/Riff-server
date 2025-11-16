package controller

type Controller struct {
	Modules           map[string]Module
	ModulesReceivers  []chan CommandBatch
	QueueCommands     []Command
	ExternalReceivers chan CommandBatch
}

func NewController() *Controller {
	return &Controller{
		Modules:           make(map[string]Module),
		ModulesReceivers:  make([]chan CommandBatch, 0),
		QueueCommands:     make([]Command, 0),
		ExternalReceivers: make(chan CommandBatch),
	}
}
