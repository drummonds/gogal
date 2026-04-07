//go:build !(js && wasm)

// Example 01: Sparkline
//
// Interactive sparkline demo. Each click of Start generates a new
// random sparkline with 7 points. Uses the lofigui framework for
// browser interactivity via WASM.
//
// Server mode: http://localhost:1340
package main

import "codeberg.org/hum3/lofigui"

func main() {
	app := lofigui.NewApp()
	app.Run(":1340", model)
}
