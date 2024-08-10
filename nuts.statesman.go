package gonuts

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// StateID is a unique identifier for a state
type StateID string

// EventID is a unique identifier for an event
type EventID string

// Action represents a function to be executed
type SMAction func(context interface{})

// Condition represents a function that returns a boolean
type SMCondition func(context interface{}) bool

// State represents a state in the state machine
type State struct {
	ID           StateID
	Name         string
	EntryActions []SMAction
	ExitActions  []SMAction
	ParentID     *StateID
	Children     []StateID
}

// Transition represents a transition between states
type Transition struct {
	From      StateID
	To        StateID
	Event     EventID
	Condition SMCondition
	Actions   []SMAction
}

// TimedTransition represents a transition that occurs after a specified duration
type TimedTransition struct {
	Transition
	Duration time.Duration
	timer    *time.Timer
}

// Region represents an orthogonal region in the state machine
type Region struct {
	ID           string
	InitialState StateID
	CurrentState StateID
}

// StatesMan is a flexible state machine manager
type StatesMan struct {
	Name             string
	mu               sync.RWMutex
	States           map[StateID]*State
	Transitions      []Transition
	TimedTransitions []TimedTransition
	Regions          []Region
	EventChannel     chan EventID
	Context          interface{}
	PreHooks         []SMAction
	PostHooks        []SMAction
}

// NewStatesMan creates a new StatesMan instance
func NewStatesMan(name string) *StatesMan {
	return &StatesMan{
		Name:         name,
		States:       make(map[StateID]*State),
		Transitions:  []Transition{},
		Regions:      []Region{{ID: "main"}},
		EventChannel: make(chan EventID, 10),
	}
}

// AddState adds a new state to the state machine
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

// AddChildState adds a child state to a parent state
func (sm *StatesMan) AddChildState(parentID, childID StateID) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	parent, exists := sm.States[parentID]
	if !exists {
		return fmt.Errorf("parent state %s does not exist", parentID)
	}
	child, exists := sm.States[childID]
	if !exists {
		return fmt.Errorf("child state %s does not exist", childID)
	}
	parent.Children = append(parent.Children, childID)
	child.ParentID = &parentID
	return nil
}

// SetInitialState sets the initial state of the state machine
func (sm *StatesMan) SetInitialState(id StateID) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, exists := sm.States[id]; !exists {
		return fmt.Errorf("state %s does not exist", id)
	}
	sm.Regions[0].InitialState = id
	sm.Regions[0].CurrentState = id
	return nil
}

// AddTransition adds a new transition to the state machine
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

// AddTimedTransition adds a new timed transition to the state machine
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

// TriggerEvent triggers an event in the state machine
func (sm *StatesMan) TriggerEvent(event EventID) {
	sm.EventChannel <- event
}

// AddPreHook adds a pre-transition hook
func (sm *StatesMan) AddPreHook(hook SMAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.PreHooks = append(sm.PreHooks, hook)
}

// AddPostHook adds a post-transition hook
func (sm *StatesMan) AddPostHook(hook SMAction) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.PostHooks = append(sm.PostHooks, hook)
}

// Run starts the state machine
func (sm *StatesMan) Run() {
	for event := range sm.EventChannel {
		sm.mu.Lock()
		sm.handleEvent(event)
		sm.mu.Unlock()
	}
}

func (sm *StatesMan) handleEvent(event EventID) {
	for _, region := range sm.Regions {
		currentState := sm.States[region.CurrentState]
		for _, t := range sm.Transitions {
			if t.From == currentState.ID && t.Event == event {
				if t.Condition == nil || t.Condition(sm.Context) {
					sm.executeTransition(&region, currentState, sm.States[t.To], t.Actions)
					break
				}
			}
		}
	}
	sm.checkTimedTransitions()
}

func (sm *StatesMan) executeTransition(region *Region, from, to *State, actions []SMAction) {
	// Execute pre-hooks
	for _, hook := range sm.PreHooks {
		hook(sm.Context)
	}

	// Execute exit actions of the current state
	for _, action := range from.ExitActions {
		action(sm.Context)
	}

	// Execute transition actions
	for _, action := range actions {
		action(sm.Context)
	}

	// Update current state
	region.CurrentState = to.ID

	// Execute entry actions of the new state
	for _, action := range to.EntryActions {
		action(sm.Context)
	}

	// Execute post-hooks
	for _, hook := range sm.PostHooks {
		hook(sm.Context)
	}

	// Reset and start timed transitions for the new state
	sm.resetTimedTransitions(to.ID)
}

func (sm *StatesMan) checkTimedTransitions() {
	for i, tt := range sm.TimedTransitions {
		if tt.From == sm.Regions[0].CurrentState && tt.timer == nil {
			sm.TimedTransitions[i].timer = time.AfterFunc(tt.Duration, func() {
				sm.mu.Lock()
				defer sm.mu.Unlock()
				sm.executeTransition(&sm.Regions[0], sm.States[tt.From], sm.States[tt.To], tt.Actions)
			})
		}
	}
}

func (sm *StatesMan) resetTimedTransitions(stateID StateID) {
	for i := range sm.TimedTransitions {
		// Stop all existing timers
		if sm.TimedTransitions[i].timer != nil {
			sm.TimedTransitions[i].timer.Stop()
			sm.TimedTransitions[i].timer = nil
		}

		// Start new timers for transitions from the current state
		if sm.TimedTransitions[i].From == stateID {
			sm.TimedTransitions[i].timer = time.AfterFunc(sm.TimedTransitions[i].Duration, func() {
				sm.mu.Lock()
				defer sm.mu.Unlock()
				sm.executeTransition(&sm.Regions[0], sm.States[sm.TimedTransitions[i].From], sm.States[sm.TimedTransitions[i].To], sm.TimedTransitions[i].Actions)
			})
		}
	}
}

// GetCurrentState returns the current state of the state machine
func (sm *StatesMan) GetCurrentState() StateID {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.Regions[0].CurrentState
}

// Export exports the state machine configuration to JSON
func (sm *StatesMan) Export() (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	type exportStatesMan struct {
		Name        string
		States      map[StateID]*State
		Transitions []Transition
		Regions     []Region
	}
	export := exportStatesMan{
		Name:        sm.Name,
		States:      sm.States,
		Transitions: sm.Transitions,
		Regions:     sm.Regions,
	}
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Import imports the state machine configuration from JSON
func (sm *StatesMan) Import(jsonStr string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	type importStatesMan struct {
		Name        string
		States      map[StateID]*State
		Transitions []Transition
		Regions     []Region
	}
	var imp importStatesMan
	err := json.Unmarshal([]byte(jsonStr), &imp)
	if err != nil {
		return err
	}
	sm.Name = imp.Name
	sm.States = imp.States
	sm.Transitions = imp.Transitions
	sm.Regions = imp.Regions
	return nil
}

// GenerateDOT generates a DOT representation of the state machine for visualization
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

// Example usage:
//
// func main() {
//     sm := gonuts.NewStatesMan("TrafficLight")
//
//     sm.AddState("Red", "Red Light", []gonuts.Action{func(ctx interface{}) { fmt.Println("Red light on") }}, nil)
//     sm.AddState("Yellow", "Yellow Light", []gonuts.Action{func(ctx interface{}) { fmt.Println("Yellow light on") }}, nil)
//     sm.AddState("Green", "Green Light", []gonuts.Action{func(ctx interface{}) { fmt.Println("Green light on") }}, nil)
//
//     sm.SetInitialState("Red")
//
//     sm.AddTransition("Red", "Green", "Next", nil)
//     sm.AddTransition("Green", "Yellow", "Next", nil)
//     sm.AddTransition("Yellow", "Red", "Next", nil)
//
//     sm.AddTimedTransition("Green", "Yellow", 30*time.Second)
//     sm.AddTimedTransition("Yellow", "Red", 5*time.Second)
//
//     sm.AddPreHook(func(ctx interface{}) { fmt.Println("About to change state") })
//     sm.AddPostHook(func(ctx interface{}) { fmt.Println("State changed") })
//
//     go sm.Run()
//
//     for i := 0; i < 6; i++ {
//         time.Sleep(2 * time.Second)
//         sm.TriggerEvent("Next")
//         fmt.Printf("Current State: %s\n", sm.GetCurrentState())
//     }
//
//     jsonExport, _ := sm.Export()
//     fmt.Println("Exported JSON:", jsonExport)
//
//     dotOutput := sm.GenerateDOT()
//     fmt.Println("DOT representation:", dotOutput)
// }
