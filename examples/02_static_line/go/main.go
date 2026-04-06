//go:build !js

// Example 02: Static Line
//
// Full line chart with axes, legend, grid, CSS tooltips, and animation.
// Serves on http://localhost:1341.
package main

import (
	"fmt"
	"net/http"
	"time"

	"codeberg.org/hum3/gogal"
)

// measurements records irregular weather-station readings throughout the day.
// Minutes are deliberately non-uniform to show how gogal handles uneven temporal spacing.
var measurements = []struct {
	hour, min int
	temp      float64
}{
	{0, 0, 11.2}, {1, 15, 10.5}, {2, 45, 9.8}, {3, 10, 9.1},
	{4, 50, 8.7}, {6, 5, 9.4}, {7, 0, 11.0}, {7, 40, 12.3},
	{8, 30, 14.1}, {9, 55, 16.8}, {10, 20, 18.5}, {11, 0, 20.1},
	{12, 35, 22.4}, {13, 10, 23.0}, {14, 0, 22.7}, {14, 45, 21.9},
	{15, 50, 20.6}, {16, 30, 19.2}, {17, 15, 17.8}, {18, 40, 15.9},
	{19, 0, 14.3}, {20, 25, 13.1}, {21, 50, 12.0}, {23, 10, 11.4},
}

func buildChart() (*gogal.Chart, []gogal.DataPoint) {
	var points []gogal.DataPoint
	for _, m := range measurements {
		t := time.Date(2024, 6, 15, m.hour, m.min, 0, 0, time.UTC)
		points = append(points, gogal.DataPoint{
			Time:  t,
			Y:     m.temp,
			Label: fmt.Sprintf("%.1f\u00b0C", m.temp),
		})
	}

	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithTitle("Temperature \u2014 15 June 2024"),
		gogal.WithXTitle("Hour"),
		gogal.WithYTitle("\u00b0C"),
		gogal.WithGrid(true),
		gogal.WithTooltips(true),
		gogal.WithAnimate(true),
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
		gogal.WithYFormat("%.0f"),
	)
	chart.Add("Temperature", points)
	return chart, points
}

const cellStyle = `border: 1px solid #ccc; padding: 4px 10px;`
const hdrStyle = cellStyle + ` background: #f5f5f5;`

func writeDebugTable(w http.ResponseWriter, chart *gogal.Chart, points []gogal.DataPoint) {
	layout := chart.Layout()

	// Summary stats
	fmt.Fprint(w, `<h2>Debug: layout pipeline</h2>`)
	fmt.Fprintf(w, `<p><strong>PlotArea:</strong> x=%.1f y=%.1f w=%.1f h=%.1f</p>`,
		layout.PlotArea.X, layout.PlotArea.Y, layout.PlotArea.Width, layout.PlotArea.Height)

	if len(points) > 0 {
		xMin, xMax := points[0].X, points[0].X
		yMin, yMax := points[0].Y, points[0].Y
		for _, p := range points[1:] {
			if p.X < xMin {
				xMin = p.X
			}
			if p.X > xMax {
				xMax = p.X
			}
			if p.Y < yMin {
				yMin = p.Y
			}
			if p.Y > yMax {
				yMax = p.Y
			}
		}
		fmt.Fprintf(w, `<p><strong>Data X range:</strong> %.0f — %.0f (span %.0f)</p>`, xMin, xMax, xMax-xMin)
		fmt.Fprintf(w, `<p><strong>Data Y range:</strong> %.1f — %.1f (span %.1f)</p>`, yMin, yMax, yMax-yMin)
	}

	// Per-point table: input data → DataPoint.X → plot coords
	fmt.Fprintf(w, `<h3>Points (%d)</h3>`, len(points))
	fmt.Fprint(w, `<table style="border-collapse: collapse; font-size: 0.82em; margin-top: 0.5em;">
<thead><tr>
<th style="`+hdrStyle+`">#</th>
<th style="`+hdrStyle+`">Time</th>
<th style="`+hdrStyle+`">DataPoint.X</th>
<th style="`+hdrStyle+`">Y</th>
<th style="`+hdrStyle+`">Plot X (svg)</th>
<th style="`+hdrStyle+`">Plot Y (svg)</th>
</tr></thead><tbody>
`)

	var plotPoints []gogal.PointLayout
	if len(layout.Series) > 0 {
		plotPoints = layout.Series[0].Points
	}
	for i, p := range points {
		plotX, plotY := "-", "-"
		if i < len(plotPoints) {
			plotX = fmt.Sprintf("%.1f", plotPoints[i].X)
			plotY = fmt.Sprintf("%.1f", plotPoints[i].Y)
		}
		fmt.Fprintf(w, `<tr>
<td style="`+cellStyle+`">%d</td>
<td style="`+cellStyle+`">%s</td>
<td style="`+cellStyle+` text-align:right;">%.0f</td>
<td style="`+cellStyle+` text-align:right;">%.1f</td>
<td style="`+cellStyle+` text-align:right;">%s</td>
<td style="`+cellStyle+` text-align:right;">%s</td>
</tr>
`, i, p.Time.Format("15:04"), p.X, p.Y, plotX, plotY)
	}
	fmt.Fprint(w, "</tbody></table>\n")
}

func main() {
	chart, points := buildChart()
	svg, _ := chart.RenderString()
	fmt.Println(svg)

	fmt.Println("Serving static line chart at http://localhost:1341")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Static Line</title></head>
<body style="font-family: system-ui; padding: 2em; max-width: 900px; margin: 0 auto;">
<h1>Static Line Chart</h1>
`)
		chart.Render(w)
		writeDebugTable(w, chart, points)
		fmt.Fprint(w, "\n</body></html>")
	})
	http.ListenAndServe(":1341", nil)
}
