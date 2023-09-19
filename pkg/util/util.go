package util

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

func Workers(workers, jobs int) int {
	if workers > jobs {
		workers = jobs
	}

	return workers
}

func TimeTrack(start time.Time) string {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(2)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*/(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	return fmt.Sprintf("%s took %s", name, elapsed)
}
