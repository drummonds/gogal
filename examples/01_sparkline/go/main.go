// Example 01: Sparkline
//
// Demonstrates the simplest gogal usage: a minimal inline chart
// with no axes, labels, or legend. Outputs SVG to stdout and
// also serves it on http://localhost:1340.
package main

import (
	"fmt"
	"math"
	"net/http"

	"codeberg.org/hum3/gogal"
)

func main() {
	// Generate some sample data: a sine wave with noise
	var points []gogal.DataPoint
	for i := 0; i < 50; i++ {
		x := float64(i)
		y := math.Sin(x*0.3)*10 + 20 + math.Sin(x*1.1)*3
		points = append(points, gogal.DataPoint{X: x, Y: y})
	}

	// Create a sparkline
	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(200, 30),
	)
	chart.Add("temperature", points)

	// Write to stdout
	svg, _ := chart.RenderString()
	fmt.Println(svg)

	// Also serve on HTTP
	fmt.Println("Serving sparkline at http://localhost:1340")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Sparkline</title></head>
<body style="font-family: system-ui; padding: 2em;">
<h1>gogal Sparkline Example</h1>
<p>Current temperature: `)
		chart.Render(w)
		fmt.Fprint(w, `</p>
<p>Smooth variant: `)

		smoothChart := gogal.NewLineChart(
			gogal.WithVariant(gogal.Sparkline),
			gogal.WithSize(200, 30),
			gogal.WithSmooth(true),
		)
		smoothChart.Add("temperature", points)
		smooth, _ := smoothChart.RenderString()
		fmt.Fprintf(w, "%s</p>\n", smooth)

		fmt.Fprintf(w, `<p>Wider: `)
		wideChart := gogal.NewLineChart(
			gogal.WithVariant(gogal.Sparkline),
			gogal.WithSize(400, 40),
		)
		wideChart.Add("temperature", points)
		wideChart.Render(w)
		fmt.Fprintf(w, `</p>
</body></html>`)
	})
	http.ListenAndServe(":1340", nil)
}
