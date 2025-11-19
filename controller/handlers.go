package controller

import "time"

func (c *Controller) handlerAddModule(cmd Command) {
	configAny, existsConfig := cmd.Data["config"]
	constructorAny, existsConstructor := cmd.Data["constructor"]

	result := CommandResult{
		Command:   cmd,
		Timestamp: time.Now(),
	}

	defer func() {
		cmd.Chanback <- result
		close(cmd.Chanback)
	}()

	if !existsConfig {
		result.Error = ErrConfigNotFound
		return
	}

	if !existsConstructor {
		result.Error = ErrConstructorNotFound
		return
	}

	constructor, okConstructorFormat := constructorAny.(ModuleConstructor)
	if !okConstructorFormat {
		result.Error = ErrConstructorType
		return
	}

	config, okConfig := configAny.(map[string]any)
	if !okConfig {
		result.Error = ErrConfigNotFound
		return
	}

	receiver := make(chan Command)

	instance, name, version, err := constructor(config, receiver)

	if err != nil {
		result.Error = err
		return
	}

	if _, exists := c.Modules[name]; exists {
		result.Error = ErrModuleAlreadyExists
		return
	}

	c.Modules[name] = Module{
		Name:     name,
		Instance: instance,
		Version:  version,
		Config:   config,
	}

	result.Result = map[string]any{
		"name":          name,
		"version":       version,
		"modules_count": len(c.Modules),
	}
}

func (c *Controller) handlerRemoveModule(cmd Command) {

}

func (c *Controller) handlerListModules(cmd Command) {

}
