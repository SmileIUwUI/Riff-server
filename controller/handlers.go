package controller

import (
	"fmt"
	"reflect"
	"time"
)

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

	instance, name, version, err := c.callConstructor(constructorAny, configAny)
	if err != nil {
		result.Error = err
		return
	}

	if _, exists := c.modules[name]; exists {
		result.Error = ErrModuleAlreadyExists
		return
	}

	c.modules[name] = Module{
		Name:     name,
		Instance: instance,
		Version:  version,
		Config:   configAny.(map[string]any),
	}

	result.Result = map[string]any{
		"name":          name,
		"version":       version,
		"modules_count": len(c.modules),
	}
}

func (c *Controller) callConstructor(constructorAny any, configAny any) (instance any, name string, version string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("constructor panic: %v", r)
		}
	}()

	constructorValue := reflect.ValueOf(constructorAny)

	if constructorValue.Kind() != reflect.Func {
		return nil, "", "", ErrConstructorNotFunction
	}

	constructorType := constructorValue.Type()

	if constructorType.NumIn() != 2 {
		return nil, "", "", ErrConstructorInvalidParams
	}

	if constructorType.NumOut() != 4 {
		return nil, "", "", ErrConstructorInvalidReturns
	}

	if constructorType.Out(1).Kind() != reflect.String ||
		constructorType.Out(2).Kind() != reflect.String {
		return nil, "", "", ErrConstructorInvalidReturnTypes
	}

	if !constructorType.Out(3).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, "", "", ErrConstructorLastReturnNotError
	}

	receiver := make(chan Command)

	config, ok := configAny.(map[string]any)
	if !ok {
		return nil, "", "", ErrConfigNotMap
	}

	args := []reflect.Value{
		reflect.ValueOf(config),
		reflect.ValueOf(receiver),
	}

	results := constructorValue.Call(args)

	instance = results[0].Interface()
	name = results[1].String()
	version = results[2].String()

	if errResult := results[3].Interface(); errResult != nil {
		err = errResult.(error)
	}

	return instance, name, version, err
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

	_, exists := c.modules[name]

	if !exists {
		result.Error = ErrModuleNotFound
		return
	}

	delete(c.modules, name)

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

	for name, module := range c.modules {
		resultMap[name] = map[string]any{
			"name":    module.Name,
			"version": module.Version,
			"config":  module.Config,
		}
	}

	result.Result = map[string]any{
		"modules": resultMap,
		"count":   len(c.modules),
	}
}
