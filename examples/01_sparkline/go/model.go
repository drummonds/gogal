package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"codeberg.org/hum3/gogal"
	"codeberg.org/hum3/lofigui"
)

// model generates sparkline points one at a time, updating the chart
// after each new point so the user sees it grow.
func model(app *lofigui.App) {
	const nPoints = 7

	var points []gogal.DataPoint
	y := 20.0

	for i := range nPoints {
		y += rand.Float64()*10 - 5 // random walk: ±5
		points = append(points, gogal.DataPoint{X: float64(i), Y: y})

		lofigui.Reset()
		lofigui.HTML(fmt.Sprintf(`<p class="has-text-grey is-size-7">Point %d of %d</p>`, i+1, nPoints))
		lofigui.HTML(`<p>Sparkline: ` + renderSparkline(points, false) + `</p>`)
		lofigui.HTML(`<p>Smooth: ` + renderSparkline(points, true) + `</p>`)

		if i < nPoints-1 {
			app.Sleep(500 * time.Millisecond)
		}
	}
}

func renderSparkline(points []gogal.DataPoint, smooth bool) string {
	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(300, 40),
		gogal.WithSmooth(smooth),
	)
	chart.Add("data", points)
	svg, _ := chart.RenderString()
	return svg
}
