//go:build !(js && wasm)

// Example 01a: Axis Formats
//
// Visual test of axis tick formatting across different data ranges:
// numeric (small, large, decimal, negative), scientific, temporal
// (hours, days, months, years), and ordinal.
//
// Server mode: http://localhost:1339
package main

import "codeberg.org/hum3/lofigui"

func main() {
	app := lofigui.NewApp()
	app.Run(":1339", model)
}
