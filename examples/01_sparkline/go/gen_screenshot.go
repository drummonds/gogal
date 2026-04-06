//go:build ignore

// Generates screenshot SVGs for example 01 documentation.
// Run: go run gen_screenshot.go
// Outputs: ../../../docs/01_sparkline/01_regular.svg, 01_smooth.svg, 01_wide.svg
package main

import (
	"math"
	"os"

	"codeberg.org/hum3/gogal"
)

func main() {
	var points []gogal.DataPoint
	for i := 0; i < 50; i++ {
		x := float64(i)
		y := math.Sin(x*0.3)*10 + 20 + math.Sin(x*1.1)*3
		points = append(points, gogal.DataPoint{X: x, Y: y})
	}

	dir := "../../../docs/01_sparkline/"

	// Regular sparkline
	regular := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(200, 30),
	)
	regular.Add("temperature", points)
	f, _ := os.Create(dir + "01_regular.svg")
	regular.Render(f)
	f.Close()

	// Smooth sparkline
	smooth := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(200, 30),
		gogal.WithSmooth(true),
	)
	smooth.Add("temperature", points)
	f, _ = os.Create(dir + "01_smooth.svg")
	smooth.Render(f)
	f.Close()

	// Wide sparkline
	wide := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(400, 40),
	)
	wide.Add("temperature", points)
	f, _ = os.Create(dir + "01_wide.svg")
	wide.Render(f)
	f.Close()
}
