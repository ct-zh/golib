package tracktime

import (
	"fmt"
	"time"
)

// Do take a time track to count how long it has been executed.
// @usage in function start:  defer tracktime.Do("function_name")()
func Do(name string) func() {
	start := time.Now()

	return func() {
		// replace the dotting log
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
