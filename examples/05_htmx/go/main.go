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
)

func main() {
	fmt.Println("Serving HTMX interactive chart at http://localhost:1344")

	http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
		hidden := parseHiddenList(r.URL.Query().Get("hidden"))
		w.Header().Set("Content-Type", "text/html")
		renderChart(hidden).Render(w)

		// Legend toggle buttons
		fmt.Fprint(w, `<div style="margin-top: 1em;">`)
		for _, name := range seriesNames {
			newHidden := toggleHidden(hidden, name)
			fmt.Fprintf(w, `<button style="%s" hx-get="/chart%s" hx-target="#chart-container" hx-swap="innerHTML">%s</button>`,
				toggleButtonStyle(isHidden(hidden, name)), hiddenParam(newHidden), name)
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
