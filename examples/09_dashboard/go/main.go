//go:build !js

// Example 09: Dashboard
//
// Multiple charts on one page: different sizes, themes, and configurations.
// Shows how to compose a dashboard from gogal charts.
// Serves on http://localhost:1348.
package main

import (
	"fmt"
	"net/http"
	"time"

	"codeberg.org/hum3/gogal"
)

type reading struct {
	hour, min int
	val       float64
}

func toPoints(raw []reading, fmtStr string) []gogal.DataPoint {
	pts := make([]gogal.DataPoint, len(raw))
	for i, r := range raw {
		t := time.Date(2024, 6, 15, r.hour, r.min, 0, 0, time.UTC)
		pts[i] = gogal.DataPoint{Time: t, Y: r.val, Label: fmt.Sprintf(fmtStr, r.val)}
	}
	return pts
}

func sensorData() (temp, humidity, pressure []gogal.DataPoint) {
	temp = toPoints([]reading{
		{0, 0, 11.2}, {1, 20, 10.4}, {2, 50, 9.7}, {4, 10, 8.9},
		{5, 30, 9.1}, {6, 45, 10.8}, {7, 15, 12.5}, {8, 40, 14.8},
		{9, 50, 17.2}, {10, 30, 19.0}, {11, 10, 20.6}, {12, 25, 22.1},
		{13, 40, 23.0}, {14, 15, 22.5}, {15, 30, 21.3}, {16, 50, 19.7},
		{18, 0, 17.4}, {19, 20, 15.6}, {20, 40, 14.0}, {22, 0, 12.3},
		{23, 30, 11.5},
	}, "%.1f\u00b0C")
	humidity = toPoints([]reading{
		{0, 0, 82.0}, {2, 30, 85.0}, {4, 45, 88.0}, {6, 30, 84.0},
		{8, 15, 72.0}, {10, 0, 58.0}, {11, 30, 50.0}, {13, 0, 45.0},
		{15, 0, 48.0}, {17, 0, 55.0}, {19, 0, 65.0}, {21, 30, 76.0},
		{23, 0, 80.0},
	}, "%.0f%%")
	pressure = toPoints([]reading{
		{0, 0, 1015.0}, {3, 0, 1014.0}, {6, 0, 1013.0}, {9, 0, 1016.0},
		{12, 0, 1018.0}, {15, 0, 1017.0}, {18, 0, 1014.0}, {21, 0, 1013.0},
		{23, 45, 1014.0},
	}, "%.0f hPa")
	return
}

func main() {
	temp, humidity, pressure := sensorData()

	fmt.Println("Serving dashboard at http://localhost:1348")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Dashboard</title></head>
<body style="font-family: system-ui; padding: 2em; max-width: 1100px; margin: 0 auto;">
<h1>Weather Dashboard</h1>

<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1em;">
<div>
<h3>Temperature</h3>
`)
		// Temperature — light theme, smooth
		tc := gogal.NewLineChart(
			gogal.WithVariant(gogal.Static),
			gogal.WithSize(500, 300),
			gogal.WithGrid(true),
			gogal.WithTooltips(true),
			gogal.WithSmooth(true),
			gogal.WithTimeFormat("15:04"),
			gogal.WithYFormat("%.0f\u00b0C"),
		)
		tc.Add("Temperature", temp)
		tc.Render(w)

		fmt.Fprint(w, `</div>
<div>
<h3>Humidity</h3>
`)
		// Humidity — dark theme
		hc := gogal.NewLineChart(
			gogal.WithVariant(gogal.Static),
			gogal.WithSize(500, 300),
			gogal.WithTheme(gogal.ThemeDark),
			gogal.WithGrid(true),
			gogal.WithTooltips(true),
			gogal.WithSmooth(true),
			gogal.WithTimeFormat("15:04"),
			gogal.WithYFormat("%.0f%%"),
		)
		hc.Add("Humidity", humidity)
		hc.Render(w)

		fmt.Fprint(w, `</div>
</div>

<h3>All Metrics (Ordinal)</h3>
`)
		// All three — ordinal axis
		all := gogal.NewLineChart(
			gogal.WithVariant(gogal.Static),
			gogal.WithSize(1060, 350),
			gogal.WithAxisMode(gogal.Ordinal),
			gogal.WithGrid(true),
			gogal.WithTooltips(true),
			gogal.WithSmooth(true),
			gogal.WithLegend(true),
		)
		all.Add("Temperature (\u00b0C)", temp)
		all.Add("Humidity (%)", humidity)
		all.Add("Pressure (hPa)", pressure)
		all.Render(w)

		fmt.Fprint(w, `

<div style="display: flex; gap: 2em; margin-top: 1em;">
`)
		// Sparklines
		for _, s := range []struct {
			name   string
			points []gogal.DataPoint
		}{
			{"Temp", temp},
			{"Humidity", humidity},
			{"Pressure", pressure},
		} {
			sp := gogal.NewLineChart(
				gogal.WithVariant(gogal.Sparkline),
				gogal.WithSize(150, 25),
				gogal.WithSmooth(true),
			)
			sp.Add(s.name, s.points)
			fmt.Fprintf(w, `<span style="font-size: 14px;">%s: `, s.name)
			sp.Render(w)
			fmt.Fprint(w, `</span>`)
		}

		fmt.Fprint(w, `
</div>
</body></html>`)
	})
	http.ListenAndServe(":1348", nil)
}
