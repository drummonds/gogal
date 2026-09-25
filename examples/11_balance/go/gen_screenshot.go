//go:build ignore

package main

import (
	"os"
	"os/exec"
)

// The chart lives in chart.go; run the package's own generator so the
// screenshot and the served chart never drift apart.
func main() {
	cmd := exec.Command("go", "run", "-tags", "screenshot", ".")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
