//go:build !js

// Example 03: Multi-Series
//
// Multiple data series with color cycling. Shows both Temporal and
// Ordinal axis modes side by side. Serves on http://localhost:1342.
package main

import (
	"fmt"
	"net/http"
	"time"

	"codeberg.org/hum3/gogal"
)

// sensorData returns irregularly-sampled readings from two sensors.
// Temperature readings are roughly every 1–2 hours; humidity roughly every 2–3 hours.
// This makes temporal vs ordinal comparison meaningful — temporal shows the gaps.
func sensorData() (temp, humidity []gogal.DataPoint) {
	type reading struct {
		hour, min int
		val       float64
	}
	tempR := []reading{
		{0, 0, 11.2}, {1, 20, 10.4}, {2, 50, 9.7}, {4, 10, 8.9},
		{5, 30, 9.1}, {6, 45, 10.8}, {7, 15, 12.5}, {8, 40, 14.8},
		{9, 50, 17.2}, {10, 30, 19.0}, {11, 10, 20.6}, {12, 25, 22.1},
		{13, 40, 23.0}, {14, 15, 22.5}, {15, 30, 21.3}, {16, 50, 19.7},
		{18, 0, 17.4}, {19, 20, 15.6}, {20, 40, 14.0}, {22, 0, 12.3},
		{23, 30, 11.5},
	}
	humR := []reading{
		{0, 0, 82.0}, {2, 30, 85.0}, {4, 45, 88.0}, {6, 30, 84.0},
		{8, 15, 72.0}, {10, 0, 58.0}, {11, 30, 50.0}, {13, 0, 45.0},
		{15, 0, 48.0}, {17, 0, 55.0}, {19, 0, 65.0}, {21, 30, 76.0},
		{23, 0, 80.0},
	}
	for _, r := range tempR {
		t := time.Date(2024, 6, 15, r.hour, r.min, 0, 0, time.UTC)
		temp = append(temp, gogal.DataPoint{Time: t, Y: r.val, Label: fmt.Sprintf("%.1f\u00b0C", r.val)})
	}
	for _, r := range humR {
		t := time.Date(2024, 6, 15, r.hour, r.min, 0, 0, time.UTC)
		humidity = append(humidity, gogal.DataPoint{Time: t, Y: r.val, Label: fmt.Sprintf("%.0f%%", r.val)})
	}
	return
}

func temporalChart(temp, humidity []gogal.DataPoint) *gogal.Chart {
	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithAxisMode(gogal.Temporal),
		gogal.WithTitle("Weather \u2014 Temporal Axis"),
		gogal.WithGrid(true),
		gogal.WithTooltips(true),
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
	)
	chart.Add("Temperature (\u00b0C)", temp)
	chart.Add("Humidity (%)", humidity)
	return chart
}

func ordinalChart(temp, humidity []gogal.DataPoint) *gogal.Chart {
	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithAxisMode(gogal.Ordinal),
		gogal.WithTitle("Weather \u2014 Ordinal Axis"),
		gogal.WithGrid(true),
		gogal.WithTooltips(true),
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
	)
	chart.Add("Temperature (\u00b0C)", temp)
	chart.Add("Humidity (%)", humidity)
	return chart
}

func writeDataTables(w http.ResponseWriter, temp, humidity []gogal.DataPoint) {
	const cs = `border: 1px solid #ccc; padding: 4px 10px;`
	const hs = cs + ` background: #f5f5f5;`
	fmt.Fprint(w, `<h2>Data</h2><div style="display: flex; gap: 2em; flex-wrap: wrap;">`)

	// Temperature table
	fmt.Fprint(w, `<div><h3>Temperature</h3>
<table style="border-collapse: collapse; font-size: 0.82em;">
<thead><tr><th style="`+hs+`">Time</th><th style="`+hs+`">°C</th></tr></thead><tbody>
`)
	for _, p := range temp {
		fmt.Fprintf(w, `<tr><td style="`+cs+`">%s</td><td style="`+cs+` text-align:right;">%.1f</td></tr>`+"\n",
			p.Time.Format("15:04"), p.Y)
	}
	fmt.Fprint(w, `</tbody></table></div>`)

	// Humidity table
	fmt.Fprint(w, `<div><h3>Humidity</h3>
<table style="border-collapse: collapse; font-size: 0.82em;">
<thead><tr><th style="`+hs+`">Time</th><th style="`+hs+`">%</th></tr></thead><tbody>
`)
	for _, p := range humidity {
		fmt.Fprintf(w, `<tr><td style="`+cs+`">%s</td><td style="`+cs+` text-align:right;">%.0f</td></tr>`+"\n",
			p.Time.Format("15:04"), p.Y)
	}
	fmt.Fprint(w, `</tbody></table></div></div>`)
}

func main() {
	temp, humidity := sensorData()

	tc := temporalChart(temp, humidity)
	svg, _ := tc.RenderString()
	fmt.Println(svg)

	fmt.Println("Serving multi-series charts at http://localhost:1342")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Multi-Series</title></head>
<body style="font-family: system-ui; padding: 2em; max-width: 900px; margin: 0 auto;">
<h1>Multi-Series Charts</h1>
<h2>Temporal Axis</h2>
<p>X positions are proportional to wall-clock time &mdash; gaps in data are visible.</p>
`)
		temporalChart(temp, humidity).Render(w)
		fmt.Fprint(w, `
<h2>Ordinal Axis</h2>
<p>X positions are equally spaced by event index &mdash; regardless of time gaps.</p>
`)
		ordinalChart(temp, humidity).Render(w)
		writeDataTables(w, temp, humidity)
		fmt.Fprint(w, "\n</body></html>")
	})
	http.ListenAndServe(":1342", nil)
}
