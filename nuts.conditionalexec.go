package gonuts

// Condition represents a function that returns a boolean
type Condition func() bool

// Action represents a function that performs some action
type Action func()

// ConditionalExecution executes actions based on conditions
type ConditionalExecution struct {
	conditions []Condition
	actions    []Action
	elseAction Action
}

// NewConditionalExecution creates a new ConditionalExecution
//
// Returns:
//   - *ConditionalExecution: a new instance of ConditionalExecution
//
// Example usage:
//
//	ce := gonuts.NewConditionalExecution().
//	    If(func() bool { return x > 10 }).
//	    Then(func() { fmt.Println("x is greater than 10") }).
//	    ElseIf(func() bool { return x > 5 }).
//	    Then(func() { fmt.Println("x is greater than 5 but not greater than 10") }).
//	    Else(func() { fmt.Println("x is 5 or less") }).
//	    Execute()
func NewConditionalExecution() *ConditionalExecution {
	return &ConditionalExecution{}
}

// If adds a condition to the execution chain
//
// Parameters:
//   - condition: a function that returns a boolean
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) If(condition Condition) *ConditionalExecution {
	ce.conditions = append(ce.conditions, condition)
	return ce
}

// Then adds an action to be executed if the previous condition is true
//
// Parameters:
//   - action: a function to be executed
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) Then(action Action) *ConditionalExecution {
	ce.actions = append(ce.actions, action)
	return ce
}

// ElseIf is an alias for If to improve readability
//
// Parameters:
//   - condition: a function that returns a boolean
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) ElseIf(condition Condition) *ConditionalExecution {
	return ce.If(condition)
}

// Else adds an action to be executed if all conditions are false
//
// Parameters:
//   - action: a function to be executed
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) Else(action Action) *ConditionalExecution {
	ce.elseAction = action
	return ce
}

// Execute runs the conditional execution chain
//
// This method evaluates each condition in order and executes the corresponding
// action for the first true condition. If no conditions are true and an Else
// action is defined, it executes the Else action.
func (ce *ConditionalExecution) Execute() {
	for i, condition := range ce.conditions {
		if condition() {
			if i < len(ce.actions) {
				ce.actions[i]()
			}
			return
		}
	}
	if ce.elseAction != nil {
		ce.elseAction()
	}
}

// ExecuteWithFallthrough runs the conditional execution chain with fallthrough behavior
//
// This method is similar to Execute, but it continues to evaluate conditions and
// execute actions even after a true condition is found, until it encounters a condition
// that returns false or reaches the end of the chain.
func (ce *ConditionalExecution) ExecuteWithFallthrough() {
	for i, condition := range ce.conditions {
		if condition() {
			if i < len(ce.actions) {
				ce.actions[i]()
			}
		} else {
			return
		}
	}
	if ce.elseAction != nil {
		ce.elseAction()
	}
}

// IfThen is a convenience function for simple if-then execution
//
// Parameters:
//   - condition: a function that returns a boolean
//   - action: a function to be executed if the condition is true
//
// Example usage:
//
//	gonuts.IfThen(
//	    func() bool { return x > 10 },
//	    func() { fmt.Println("x is greater than 10") },
//	)
func IfThen(condition Condition, action Action) {
	if condition() {
		action()
	}
}

// IfThenElse is a convenience function for simple if-then-else execution
//
// Parameters:
//   - condition: a function that returns a boolean
//   - thenAction: a function to be executed if the condition is true
//   - elseAction: a function to be executed if the condition is false
//
// Example usage:
//
//	gonuts.IfThenElse(
//	    func() bool { return x > 10 },
//	    func() { fmt.Println("x is greater than 10") },
//	    func() { fmt.Println("x is 10 or less") },
//	)
func IfThenElse(condition Condition, thenAction, elseAction Action) {
	if condition() {
		thenAction()
	} else {
		elseAction()
	}
}
