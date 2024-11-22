package vabastegi

import (
	"errors"
	"fmt"
	"reflect"
)

// Hub is a contract that the T must satisfy.
type Hub interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

var (
	ErrHubIsNotStruct       = errors.New("passed hub is not struct")
	ErrHubFieldNotFound     = errors.New("hub field not found")
	ErrTheFieldCanNotChange = errors.New("the field can not change")
)

func GetHubField(hub interface{}, name string) (interface{}, error) {
	hubValue := reflect.ValueOf(hub)
	if hubValue.Kind() != reflect.Struct {
		return nil, ErrHubIsNotStruct
	}

	theValue := hubValue.FieldByName(name)
	if !theValue.IsValid() {
		return nil, fmt.Errorf("%s field on hub: %w", name, ErrHubFieldNotFound)
	}

	return theValue.Interface(), nil
}

func SetHubField[T any](app *App[T], name string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("setting %s field: %v", name, e)
		}
	}()

	hubValue := reflect.ValueOf(&app.Hub).Elem()
	if hubValue.Kind() != reflect.Struct {
		return fmt.Errorf("hub is %s: %w", hubValue.Kind(), ErrHubIsNotStruct)
	}

	theValue := hubValue.FieldByName(name)
	if !theValue.IsValid() {
		return fmt.Errorf("%s field on hub: %w", name, ErrHubFieldNotFound)
	}

	if !theValue.CanAddr() {
		return ErrTheFieldCanNotChange
	}

	theValue.Set(reflect.ValueOf(value))

	return nil
}

func areTheSame(valueA, ValueB reflect.Value) bool {
	return false
}
