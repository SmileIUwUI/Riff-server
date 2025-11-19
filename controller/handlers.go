package controller

import "time"

func (c *Controller) handlerAddModule(cmd Command) {
	result := CommandResult{
		Command: cmd,
	}

	defer func() {
		result.Timestamp = time.Now()
		cmd.Chanback <- result
		close(cmd.Chanback)
	}()

	configAny, existsConfig := cmd.Data["config"]
	constructorAny, existsConstructor := cmd.Data["constructor"]

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

	config, okConfigFormat := configAny.(map[string]any)
	if !okConfigFormat {
		result.Error = ErrConfigType
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
	result := CommandResult{
		Command: cmd,
	}

	defer func() {
		result.Timestamp = time.Now()
		cmd.Chanback <- result
		close(cmd.Chanback)
	}()

	nameAny, existsName := cmd.Data["name"]

	if !existsName {
		result.Error = ErrNameNotFound
		return
	}

	name, okNameFormat := nameAny.(string)
	if !okNameFormat {
		result.Error = ErrNameType
		return
	}

	_, exists := c.Modules[name]

	if !exists {
		result.Error = ErrModuleNotFound
		return
	}

	delete(c.Modules, name)

	result.Result = map[string]any{
		"name": name,
	}
}

func (c *Controller) handlerListModules(cmd Command) {
	result := CommandResult{
		Command: cmd,
	}

	defer func() {
		result.Timestamp = time.Now()
		cmd.Chanback <- result
		close(cmd.Chanback)
	}()

	resultMap := make(map[string]any)

	for name, module := range c.Modules {
		resultMap[name] = map[string]any{
			"name":    module.Name,
			"version": module.Version,
			"config":  module.Config,
		}
	}

	result.Result = map[string]any{
		"modules": resultMap,
		"count":   len(c.Modules),
	}
}
