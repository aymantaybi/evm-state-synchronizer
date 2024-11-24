package main

import (
	"fmt"

	"github.com/aymantaybi/ronin/common/fdlimit"
)

func Greet() {

}

// MakeDatabaseHandles raises out the number of allowed file handles per process
// for Geth and returns half of the allowance to assign to the database.
func MakeDatabaseHandles(max int) int {
	limit, err := fdlimit.Maximum()
	if err != nil {
		fmt.Printf("Failed to retrieve file descriptor allowance: %v", err)
	}
	switch {
	case max == 0:
		// User didn't specify a meaningful value, use system limits
	case max < 128:
		// User specified something unhealthy, just use system defaults
		fmt.Println("File descriptor limit invalid (<128)", "had", max, "updated", limit)
	case max > limit:
		// User requested more than the OS allows, notify that we can't allocate it
		fmt.Println("Requested file descriptors denied by OS", "req", max, "limit", limit)
	default:
		// User limit is meaningful and within allowed range, use that
		limit = max
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		fmt.Println("Failed to raise file descriptor allowance: %v", err)
	}
	return int(raised / 2) // Leave half for networking and other stuff
}
