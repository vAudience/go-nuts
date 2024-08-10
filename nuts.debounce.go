package gonuts

import (
	"reflect"
	"sync"
	"time"
)

// Debounce that executes first call immediately and last call after delays in calls
func Debounce(fn any, duration time.Duration, callback func(int)) func(...any) {
	var timer *time.Timer
	var args []reflect.Value
	var callCount int
	var firstCall bool = true
	var mutex sync.Mutex
	fnVal := reflect.ValueOf(fn)

	return func(callArgs ...any) {
		mutex.Lock()
		defer mutex.Unlock()
		callCount++

		// Convert call arguments to reflect.Value
		args = make([]reflect.Value, len(callArgs))
		for i, arg := range callArgs {
			args[i] = reflect.ValueOf(arg)
		}

		if firstCall {
			// Execute the function immediately on the first call
			fnVal.Call(args)
			firstCall = false
		}

		// Reset the timer if it's already set
		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(duration, func() {
			mutex.Lock()
			defer mutex.Unlock()

			// Call the function with the latest arguments after the duration elapses
			fnVal.Call(args)

			// If a callback is provided, call it with the number of accumulated calls
			if callback != nil {
				callback(callCount)
			}

			// Reset call count and first call state after executing
			callCount = 0
			firstCall = true
		})
	}
}

/*

func main() {
    // Example function to debounce
    printMessage := func(message string) {
        fmt.Println("Message:", message)
    }

    // Example callback to execute after debouncing
    countCalls := func(count int) {
        fmt.Println("Function was called", count, "times")
    }

    // Creating a debounced version of printMessage that executes after 2 seconds of inactivity
    debouncedPrint := Debounce(printMessage, 2*time.Second, countCalls)

    // Simulating rapid calls to the debounced function
    debouncedPrint("Hello")
    time.Sleep(1 * time.Second)
    debouncedPrint("World")
    time.Sleep(1 * time.Second)
    debouncedPrint("Again")
    time.Sleep(3 * time.Second) // Wait enough time to ensure the last call and the callback execute
}

*/
