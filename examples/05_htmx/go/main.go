//go:build !js

// Example 05: Interactive HTMX
//
// Legend toggling via HTMX. Clicking a series name sends an HTMX
// request with ?hidden= params; the server re-renders the chart
// with those series hidden using WithHiddenSeries.
// Serves on http://localhost:1344.
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"codeberg.org/hum3/gogal"
)

var seriesNames = []string{"Temperature", "Humidity", "Wind Speed"}

type reading struct {
	hour, min int
	val       float64
}

var (
	tempR = []reading{
		{0, 0, 11.2}, {1, 20, 10.4}, {2, 50, 9.7}, {4, 10, 8.9},
		{5, 30, 9.1}, {6, 45, 10.8}, {7, 15, 12.5}, {8, 40, 14.8},
		{9, 50, 17.2}, {10, 30, 19.0}, {11, 10, 20.6}, {12, 25, 22.1},
		{13, 40, 23.0}, {14, 15, 22.5}, {15, 30, 21.3}, {16, 50, 19.7},
		{18, 0, 17.4}, {19, 20, 15.6}, {20, 40, 14.0}, {22, 0, 12.3},
		{23, 30, 11.5},
	}
	humR = []reading{
		{0, 0, 82.0}, {2, 30, 85.0}, {4, 45, 88.0}, {6, 30, 84.0},
		{8, 15, 72.0}, {10, 0, 58.0}, {11, 30, 50.0}, {13, 0, 45.0},
		{15, 0, 48.0}, {17, 0, 55.0}, {19, 0, 65.0}, {21, 30, 76.0},
		{23, 0, 80.0},
	}
	windR = []reading{
		{0, 30, 6.2}, {3, 0, 4.8}, {5, 15, 5.5}, {7, 0, 8.3},
		{9, 30, 12.1}, {11, 0, 15.4}, {12, 45, 18.2}, {14, 30, 16.7},
		{16, 0, 13.9}, {18, 30, 10.5}, {20, 0, 7.8}, {22, 30, 5.9},
	}
)

func toPoints(raw []reading, fmtStr string) []gogal.DataPoint {
	pts := make([]gogal.DataPoint, len(raw))
	for i, r := range raw {
		t := time.Date(2024, 6, 15, r.hour, r.min, 0, 0, time.UTC)
		pts[i] = gogal.DataPoint{Time: t, Y: r.val, Label: fmt.Sprintf(fmtStr, r.val)}
	}
	return pts
}

func sensorData() map[string][]gogal.DataPoint {
	return map[string][]gogal.DataPoint{
		"Temperature": toPoints(tempR, "%.1f\u00b0C"),
		"Humidity":    toPoints(humR, "%.0f%%"),
		"Wind Speed":  toPoints(windR, "%.1f km/h"),
	}
}

func renderChart(hidden []string) *gogal.Chart {
	data := sensorData()
	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithTitle("Weather Station"),
		gogal.WithGrid(true),
		gogal.WithTooltips(true),
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
		gogal.WithHiddenSeries(hidden...),
	)
	for _, name := range seriesNames {
		chart.Add(name, data[name])
	}
	return chart
}

func parseHidden(r *http.Request) []string {
	h := r.URL.Query().Get("hidden")
	if h == "" {
		return nil
	}
	return strings.Split(h, ",")
}

func toggleHidden(current []string, name string) []string {
	for i, h := range current {
		if h == name {
			return append(current[:i], current[i+1:]...)
		}
	}
	return append(current, name)
}

func hiddenParam(hidden []string) string {
	if len(hidden) == 0 {
		return ""
	}
	return "?hidden=" + strings.Join(hidden, ",")
}

func main() {
	fmt.Println("Serving HTMX interactive chart at http://localhost:1344")

	http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
		hidden := parseHidden(r)
		w.Header().Set("Content-Type", "text/html")
		renderChart(hidden).Render(w)

		// Legend toggle buttons
		fmt.Fprint(w, `<div style="margin-top: 1em;">`)
		for _, name := range seriesNames {
			newHidden := toggleHidden(hidden, name)
			isHidden := false
			for _, h := range hidden {
				if h == name {
					isHidden = true
					break
				}
			}
			style := "margin-right: 0.5em; cursor: pointer; padding: 0.25em 0.75em; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5;"
			if isHidden {
				style += " opacity: 0.4; text-decoration: line-through;"
			}
			fmt.Fprintf(w, `<button style="%s" hx-get="/chart%s" hx-target="#chart-container" hx-swap="innerHTML">%s</button>`,
				style, hiddenParam(newHidden), name)
		}
		fmt.Fprint(w, `</div>`)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head>
<title>gogal - HTMX Interactive</title>
<script src="https://unpkg.com/htmx.org@2.0.4"></script>
</head>
<body style="font-family: system-ui; padding: 2em; max-width: 900px; margin: 0 auto;">
<h1>Interactive Chart (HTMX)</h1>
<p>Click series names below the chart to toggle visibility.</p>
<div id="chart-container" hx-get="/chart" hx-trigger="load" hx-swap="innerHTML">
  <p>Loading chart...</p>
</div>
</body></html>`)
	})

	http.ListenAndServe(":1344", nil)
}
