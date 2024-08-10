package gonuts

import (
	"fmt"
	"reflect"
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// EventEmitter is a flexible publish-subscribe event system with named listeners
type EventEmitter struct {
	listeners map[string]map[string]reflect.Value
	mu        sync.RWMutex
}

// NewEventEmitter creates a new EventEmitter
//
// Returns:
//   - *EventEmitter: a new instance of EventEmitter
//
// Example usage:
//
//	emitter := gonuts.NewEventEmitter()
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string]map[string]reflect.Value),
	}
}

// On subscribes a named function to an event
//
// Parameters:
//   - event: the name of the event to subscribe to
//   - name: a unique name for this listener (if empty, a unique ID will be generated)
//   - fn: the function to be called when the event is emitted
//
// Returns:
//   - string: the name or generated ID of the listener
//   - error: any error that occurred during subscription
//
// Example usage:
//
//	id, err := emitter.On("userLoggedIn", "logLoginTime", func(username string) {
//	    fmt.Printf("User logged in: %s at %v\n", username, time.Now())
//	})
//	if err != nil {
//	    log.Printf("Error subscribing to event: %v", err)
//	}
func (ee *EventEmitter) On(event, name string, fn interface{}) (string, error) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return "", fmt.Errorf("third argument to On must be a function")
	}

	ee.mu.Lock()
	defer ee.mu.Unlock()

	if ee.listeners[event] == nil {
		ee.listeners[event] = make(map[string]reflect.Value)
	}

	if name == "" {
		var err error
		name, err = gonanoid.New()
		if err != nil {
			return "", fmt.Errorf("failed to generate unique ID: %w", err)
		}
	}

	ee.listeners[event][name] = reflect.ValueOf(fn)
	return name, nil
}

// Off unsubscribes a named function from an event
//
// Parameters:
//   - event: the name of the event to unsubscribe from
//   - name: the name or ID of the listener to unsubscribe
//
// Returns:
//   - error: any error that occurred during unsubscription
func (ee *EventEmitter) Off(event, name string) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if listeners, ok := ee.listeners[event]; ok {
		if _, exists := listeners[name]; exists {
			delete(listeners, name)
			return nil
		}
	}
	return fmt.Errorf("listener not found for event: %s, name: %s", event, name)
}

// Emit triggers an event with the given arguments
//
// Parameters:
//   - event: the name of the event to emit
//   - args: the arguments to pass to the event listeners
//
// Returns:
//   - error: any error that occurred during emission
func (ee *EventEmitter) Emit(event string, args ...interface{}) error {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	if listeners, ok := ee.listeners[event]; ok {
		for _, listener := range listeners {
			if err := ee.callListener(listener, args); err != nil {
				return err
			}
		}
	}
	return nil
}

// EmitConcurrent triggers an event with the given arguments, calling listeners concurrently
//
// Parameters:
//   - event: the name of the event to emit
//   - args: the arguments to pass to the event listeners
//
// Returns:
//   - error: any error that occurred during emission
func (ee *EventEmitter) EmitConcurrent(event string, args ...interface{}) error {
	ee.mu.RLock()
	listeners := ee.listeners[event]
	ee.mu.RUnlock()

	if len(listeners) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(listeners))

	for _, listener := range listeners {
		wg.Add(1)
		go func(l reflect.Value) {
			defer wg.Done()
			if err := ee.callListener(l, args); err != nil {
				errChan <- err
			}
		}(listener)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Once subscribes a one-time named function to an event
//
// Parameters:
//   - event: the name of the event to subscribe to
//   - name: a unique name for this listener (if empty, a unique ID will be generated)
//   - fn: the function to be called when the event is emitted
//
// Returns:
//   - string: the name or generated ID of the listener
//   - error: any error that occurred during subscription
//
// The function will be automatically unsubscribed after it is called once.
func (ee *EventEmitter) Once(event, name string, fn interface{}) (string, error) {
	wrapper := reflect.MakeFunc(reflect.TypeOf(fn), func(args []reflect.Value) []reflect.Value {
		reflect.ValueOf(fn).Call(args)
		ee.Off(event, name)
		return nil
	})
	return ee.On(event, name, wrapper.Interface())
}

// ListenerCount returns the number of listeners for a given event
//
// Parameters:
//   - event: the name of the event
//
// Returns:
//   - int: the number of listeners for the event
func (ee *EventEmitter) ListenerCount(event string) int {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	return len(ee.listeners[event])
}

// ListenerNames returns a list of all listener names for a given event
//
// Parameters:
//   - event: the name of the event
//
// Returns:
//   - []string: a slice containing all listener names for the event
func (ee *EventEmitter) ListenerNames(event string) []string {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	names := make([]string, 0, len(ee.listeners[event]))
	for name := range ee.listeners[event] {
		names = append(names, name)
	}
	return names
}

// Events returns a list of all events that have listeners
//
// Returns:
//   - []string: a slice containing all events with listeners
func (ee *EventEmitter) Events() []string {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	events := make([]string, 0, len(ee.listeners))
	for event := range ee.listeners {
		events = append(events, event)
	}
	return events
}

func (ee *EventEmitter) callListener(listener reflect.Value, args []interface{}) error {
	listenerType := listener.Type()
	if listenerType.NumIn() != len(args) {
		return fmt.Errorf("event handler expects %d arguments, but got %d", listenerType.NumIn(), len(args))
	}

	callArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		expectedType := listenerType.In(i)
		argValue := reflect.ValueOf(arg)
		if !argValue.Type().AssignableTo(expectedType) {
			return fmt.Errorf("argument %d has wrong type: got %v, want %v", i, argValue.Type(), expectedType)
		}
		callArgs[i] = argValue
	}

	listener.Call(callArgs)
	return nil
}
