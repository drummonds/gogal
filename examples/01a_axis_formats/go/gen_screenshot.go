//go:build ignore

// Generates screenshot SVG for example 01a documentation.
// Run: go run gen_screenshot.go
// Outputs: ../../../docs/01a_axis_formats/01a_axis_formats.svg
package main

import (
	"os"

	"codeberg.org/hum3/gogal"
)

func main() {
	dir := "../../../docs/01a_axis_formats/"
	os.MkdirAll(dir, 0o755)

	pts := []gogal.DataPoint{
		{X: 0, Y: 1},
		{X: 2, Y: 3},
		{X: 4, Y: 2},
		{X: 6, Y: 5},
		{X: 8, Y: 4},
		{X: 10, Y: 6},
	}

	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithSize(500, 250),
		gogal.WithGrid(true),
		gogal.WithTitle("Small integers (0–10)"),
	)
	chart.Add("data", pts)

	f, _ := os.Create(dir + "01a_axis_formats.svg")
	chart.Render(f)
	f.Close()
}
