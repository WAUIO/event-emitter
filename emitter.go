package event_emitter

import (
	"errors"
	"reflect"
	"sync"
)

func New() EventEmitter {
	return EventEmitter{
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

func (e *EventEmitter) Fire(event string, payload ...interface{}) ([]interface{}, error) {
	fn, in, err := e.resolve(event, payload...)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, len(fn))

	for i, f := range fn {
		value := f.Call(in)
		results[i] = reflect.ValueOf(value[0].Interface())
	}

	return results, nil
}

func (e *EventEmitter) FireBackground(event string, payload ...interface{}) (chan interface{}, error) {
	fn, in, err := e.resolve(event, payload...)
	if err != nil {
		return nil, err
	}

	results := make(chan interface{}, len(fn))

	for _, f := range fn {
		go func(f reflect.Value) {
			value := f.Call(in)
			results <- reflect.ValueOf(value[0].Interface())
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

	fn := make([]reflect.Value, len(listeners))
	for i, listener := range listeners {
		f := reflect.ValueOf(listener)
		if len(payload) != f.Type().NumIn() {
			// f is not qualified to be a correct listener for the event
			continue
		}

		// add f into fn
		fn[i] = f
	}

	in := make([]reflect.Value, len(payload))
	for k, param := range payload {
		in[k] = reflect.ValueOf(param)
	}

	return fn, in, nil
}

func (e *EventEmitter) results(input [][]reflect.Value) (output []interface{}, err error) {
	output = make([]interface{}, len(input))
	for i, val := range input {
		output[i] = reflect.ValueOf(val[0].Interface())
	}

	return output, nil
}
