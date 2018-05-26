package event_emitter

import (
	"errors"
	"reflect"
	"sync"
)

func New() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]interface{}),
	}
}

type Listener func(payload ... interface{}) interface{}

type EventEmitter struct {
	listeners map[string][]interface{}

	mu sync.Mutex
}

func (e *EventEmitter) On(event string, listener interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if reflect.ValueOf(listener).Type().Kind() != reflect.Func {
		return errors.New("listener should be a func")
	}

	e.listeners[event] = append(e.listeners[event], listener)

	return nil
}

func (e *EventEmitter) Fire(event string, payload ...interface{}) ([][]interface{}, error) {
	fn, in, err := e.resolve(event, payload...)
	if err != nil {
		return nil, err
	}

	results := make([][]interface{}, len(fn))

	for i, f := range fn {
		// results in values is the slice of returned values from our listener
		values := f.Call(in)

		// here we try to parse this slice to a slice of concrete values
		results[i], _ = getValues(values)
	}

	return results, nil
}

func (e *EventEmitter) FireBackground(event string, payload ...interface{}) (chan []interface{}, error) {
	fn, in, err := e.resolve(event, payload...)
	if err != nil {
		return nil, err
	}

	results := make(chan []interface{}, len(fn))

	for _, f := range fn {
		go func(f reflect.Value) {
			values := f.Call(in)

			if output, err := getValues(values); err == nil {
				results <- output
			} else {
				results <- nil
			}
		}(f)
	}

	return results, nil
}

func (e *EventEmitter) Clear(event string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.listeners[event]; !ok {
		return errors.New("event not defined")
	}

	delete(e.listeners, event)

	return nil
}

func (e *EventEmitter) ClearEvents() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners = make(map[string][]interface{})
}

func (e *EventEmitter) Has(event string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, ok := e.listeners[event]

	return ok
}

func (e *EventEmitter) resolve(event string, payload ...interface{}) ([]reflect.Value, []reflect.Value, error) {
	e.mu.Lock()
	listeners, ok := e.listeners[event]
	e.mu.Unlock()

	if !ok {
		return make([]reflect.Value, 0), nil, errors.New("No listener defined for event " + event)
	}

	fns := make([]reflect.Value, len(listeners))
	for i, listener := range listeners {
		fn := reflect.ValueOf(listener)
		if len(payload) != fn.Type().NumIn() {
			// fn is not qualified to be a correct listener for the event
			continue
		}

		// add f into fn
		fns[i] = fn
	}

	in := make([]reflect.Value, len(payload))
	for k, param := range payload {
		in[k] = reflect.ValueOf(param)
	}

	return fns, in, nil
}

func getValues(input []reflect.Value) (values []interface{}, err error) {
	values = make([]interface{}, len(input))
	for i, val := range input {
		values[i] = val.Interface()
	}

	return values, nil
}

func mapValues(inputs [][]reflect.Value) (values [][]interface{}, err error) {
	values = make([][]interface{}, len(inputs))

	for i, input := range inputs {
		values[i], _ = getValues(input)
	}

	return values, err
}
