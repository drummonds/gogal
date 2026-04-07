package main

import (
	"math/rand/v2"

	"codeberg.org/hum3/gogal"
	"codeberg.org/hum3/lofigui"
)

// model contains the business logic — shared by server and WASM builds.
// Each invocation generates a new random sparkline with 7 points.
func model(app *lofigui.App) {
	points := randomPoints(7)

	regular := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(300, 40),
	)
	regular.Add("data", points)
	regSVG, _ := regular.RenderString()

	smooth := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(300, 40),
		gogal.WithSmooth(true),
	)
	smooth.Add("data", points)
	smoothSVG, _ := smooth.RenderString()

	lofigui.HTML(`<p>Sparkline: ` + regSVG + `</p>`)
	lofigui.HTML(`<p>Smooth: ` + smoothSVG + `</p>`)
}

func randomPoints(n int) []gogal.DataPoint {
	pts := make([]gogal.DataPoint, n)
	y := 20.0
	for i := range pts {
		y += rand.Float64()*10 - 5 // random walk: ±5
		pts[i] = gogal.DataPoint{X: float64(i), Y: y}
	}
	return pts
}
