package gonuts

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// StateID is a unique identifier for a state.
type StateID string

// EventID is a unique identifier for an event.
type EventID string

// SMAction represents a function to be executed during state transitions.
type SMAction func(context map[string]interface{})

// SMCondition represents a function that returns a boolean based on the context.
type SMCondition func(context map[string]interface{}) bool

// State represents a state in the state machine.
type State struct {
	ID           StateID
	Name         string
	EntryActions []SMAction
	ExitActions  []SMAction
}

// Transition represents a transition between states.
type Transition struct {
	From      StateID
	To        StateID
	Event     EventID
	Condition SMCondition
	Actions   []SMAction
}

// TimedTransition represents a transition that occurs after a specified duration.
type TimedTransition struct {
	Transition
	Duration time.Duration
	timer    *time.Timer
}

// EventData encapsulates event information passed to the state machine.
type EventData struct {
	EventID EventID
	Data    map[string]interface{}
}

// StatesMan is a flexible, concurrent-safe state machine manager.
type StatesMan struct {
	Name             string
	mu               sync.RWMutex
	States           map[StateID]*State
	Transitions      []Transition
	TimedTransitions []TimedTransition
	CurrentState     StateID
	EventChannel     chan EventData
	Context          map[string]interface{}
	PreHooks         []SMAction
	PostHooks        []SMAction
}

// AnyState represents a wildcard state that matches any current state.
const AnyState StateID = "*"

// NewStatesMan creates a new StatesMan instance.
func NewStatesMan(name string) *StatesMan {
	return &StatesMan{
		Name:         name,
		States:       make(map[StateID]*State),
		Transitions:  []Transition{},
		EventChannel: make(chan EventData, 10),
		Context:      make(map[string]interface{}),
	}
}

// AddState adds a new state to the state machine.
func (sm *StatesMan) AddState(id StateID, name string, entryActions, exitActions []SMAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.States[id] = &State{
		ID:           id,
		Name:         name,
		EntryActions: entryActions,
		ExitActions:  exitActions,
	}
}

// SetInitialState sets the initial state of the state machine.
func (sm *StatesMan) SetInitialState(id StateID) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, exists := sm.States[id]; !exists {
		return fmt.Errorf("state %s does not exist", id)
	}
	sm.CurrentState = id
	return nil
}

// AddTransition adds a new transition to the state machine.
func (sm *StatesMan) AddTransition(from, to StateID, event EventID, condition SMCondition, actions ...SMAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.Transitions = append(sm.Transitions, Transition{
		From:      from,
		To:        to,
		Event:     event,
		Condition: condition,
		Actions:   actions,
	})
}

// AddTimedTransition adds a new timed transition to the state machine.
func (sm *StatesMan) AddTimedTransition(from, to StateID, duration time.Duration, actions ...SMAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.TimedTransitions = append(sm.TimedTransitions, TimedTransition{
		Transition: Transition{
			From:    from,
			To:      to,
			Actions: actions,
		},
		Duration: duration,
	})
}

// TriggerEvent triggers an event in the state machine without knowing the next state.
func (sm *StatesMan) TriggerEvent(event EventID, data map[string]interface{}) {
	sm.EventChannel <- EventData{EventID: event, Data: data}
}

// AddPreHook adds a pre-transition hook.
func (sm *StatesMan) AddPreHook(hook SMAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.PreHooks = append(sm.PreHooks, hook)
}

// AddPostHook adds a post-transition hook.
func (sm *StatesMan) AddPostHook(hook SMAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.PostHooks = append(sm.PostHooks, hook)
}

// Run starts the state machine event loop.
func (sm *StatesMan) Run() {
	for eventData := range sm.EventChannel {
		sm.mu.Lock()
		sm.handleEvent(eventData)
		sm.mu.Unlock()
	}
}

// handleEvent processes an incoming event and executes the appropriate transition.
func (sm *StatesMan) handleEvent(eventData EventData) {
	currentState := sm.States[sm.CurrentState]
	event := eventData.EventID
	context := sm.Context

	// Merge event data into context
	for k, v := range eventData.Data {
		context[k] = v
	}

	// Process transitions - if not found, no transition is executed - we might want to handle this case in a newer version
	// foundTransition := false
	for _, t := range sm.Transitions {
		if (t.From == sm.CurrentState || t.From == AnyState) && t.Event == event {
			if t.Condition == nil || t.Condition(context) {
				sm.executeTransition(currentState, sm.States[t.To], t.Actions)
				// foundTransition = true
				break
			}
		}
	}
	// if !foundTransition {
	// 	// Optionally handle no matching transition

	// }
	sm.checkTimedTransitions()
}

// executeTransition performs the transition between states.
func (sm *StatesMan) executeTransition(from, to *State, actions []SMAction) {
	context := sm.Context

	// Execute pre-hooks
	for _, hook := range sm.PreHooks {
		hook(context)
	}

	// Execute exit actions of the current state
	for _, action := range from.ExitActions {
		action(context)
	}

	// Execute transition actions
	for _, action := range actions {
		action(context)
	}

	// Update current state
	sm.CurrentState = to.ID

	// Execute entry actions of the new state
	for _, action := range to.EntryActions {
		action(context)
	}

	// Execute post-hooks
	for _, hook := range sm.PostHooks {
		hook(context)
	}

	// Reset and start timed transitions for the new state
	sm.resetTimedTransitions(to.ID)
}

// checkTimedTransitions initializes timers for timed transitions from the current state.
func (sm *StatesMan) checkTimedTransitions() {
	for i := range sm.TimedTransitions {
		tt := &sm.TimedTransitions[i]
		if tt.From == sm.CurrentState && tt.timer == nil {
			tt.timer = time.AfterFunc(tt.Duration, func() {
				sm.mu.Lock()
				defer sm.mu.Unlock()
				sm.executeTransition(sm.States[tt.From], sm.States[tt.To], tt.Actions)
			})
		}
	}
}

// resetTimedTransitions stops existing timers and starts new ones for the given state.
func (sm *StatesMan) resetTimedTransitions(stateID StateID) {
	for i := range sm.TimedTransitions {
		tt := &sm.TimedTransitions[i]
		if tt.timer != nil {
			tt.timer.Stop()
			tt.timer = nil
		}
		if tt.From == stateID {
			tt.timer = time.AfterFunc(tt.Duration, func() {
				sm.mu.Lock()
				defer sm.mu.Unlock()
				sm.executeTransition(sm.States[tt.From], sm.States[tt.To], tt.Actions)
			})
		}
	}
}

// GetCurrentState returns the current state of the state machine.
func (sm *StatesMan) GetCurrentState() StateID {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.CurrentState
}

// Export exports the state machine configuration to JSON.
func (sm *StatesMan) Export() (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	export := struct {
		Name        string
		States      map[StateID]*State
		Transitions []Transition
	}{
		Name:        sm.Name,
		States:      sm.States,
		Transitions: sm.Transitions,
	}
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Import imports the state machine configuration from JSON.
func (sm *StatesMan) Import(jsonStr string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	var imp struct {
		Name        string
		States      map[StateID]*State
		Transitions []Transition
	}
	err := json.Unmarshal([]byte(jsonStr), &imp)
	if err != nil {
		return err
	}
	sm.Name = imp.Name
	sm.States = imp.States
	sm.Transitions = imp.Transitions
	return nil
}

// GenerateDOT generates a DOT representation of the state machine for visualization.
func (sm *StatesMan) GenerateDOT() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	dot := "digraph " + sm.Name + " {\n"
	for _, state := range sm.States {
		dot += fmt.Sprintf("  %s [label=\"%s\"];\n", state.ID, state.Name)
	}
	for _, t := range sm.Transitions {
		dot += fmt.Sprintf("  %s -> %s [label=\"%s\"];\n", t.From, t.To, t.Event)
	}
	dot += "}\n"
	return dot
}

/*
// Example usage of the updated StatesMan.
func exampleUsage() {
	// Create a new state machine
	sm := NewStatesMan("JobProcessor")

	// Define states
	sm.AddState("Idle", "Idle State", nil, nil)
	sm.AddState("Processing", "Processing State", nil, nil)
	sm.AddState("Completed", "Completed State", nil, nil)
	sm.AddState("Failed", "Failed State", nil, nil)

	// Set initial state
	sm.SetInitialState("Idle")

	// Define transitions
	sm.AddTransition("Idle", "Processing", "StartProcessing", nil)
	sm.AddTransition("Processing", "Completed", "Success", nil)
	sm.AddTransition("Processing", "Failed", "Failure", nil)
	// Allow any state to transition back to Idle on Reset event
	sm.AddTransition(AnyState, "Idle", "Reset", nil)

	// Start the state machine in a separate goroutine
	go sm.Run()

	// Function that starts processing
	startProcessing := func() {
		sm.TriggerEvent("StartProcessing", nil)
	}

	// Function that simulates job processing and triggers Success or Failure
	processJob := func() {
		// Simulate job processing...
		// On success
		sm.TriggerEvent("Success", map[string]interface{}{"result": "Job completed successfully"})
		// On failure
		// sm.TriggerEvent("Failure", map[string]interface{}{"error": errors.New("Job failed")})
	}

	// Start processing
	startProcessing()
	processJob()

	// Check current state
	currentState := sm.GetCurrentState()
	fmt.Println("Current State:", currentState)
}
*/
