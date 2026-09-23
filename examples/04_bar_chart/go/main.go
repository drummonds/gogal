//go:build !js

// Example 04: Bar Chart
//
// Vertical bar chart with category labels, value labels and a grouped
// two-series variant. Serves on http://localhost:1343.
package main

import (
	"fmt"
	"net/http"

	"git.bytestone.uk/hum3/gogal"
)

var (
	fib       = []float64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	fibLabels = []string{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10"}
	quarters  = []string{"Q1", "Q2", "Q3", "Q4"}
	sales2024 = []float64{120, 95, 140, 160}
	sales2025 = []float64{135, 110, 150, 185}
)

func fibChart(n int) *gogal.Chart {
	chart := gogal.NewBarChart(
		gogal.WithTitle(fmt.Sprintf("Fibonacci (n=%d)", n)),
		gogal.WithSize(500, 220),
		gogal.WithYFormat("%.0f"),
	)
	chart.AddCategories("Fibonacci", fibLabels[:n], fib[:n])
	return chart
}

func groupedChart() *gogal.Chart {
	chart := gogal.NewBarChart(
		gogal.WithTitle("Quarterly sales"),
		gogal.WithSize(500, 260),
		gogal.WithYTitle("k£"),
		gogal.WithYFormat("%.0f"),
	)
	chart.AddCategories("2024", quarters, sales2024)
	chart.AddCategories("2025", quarters, sales2025)
	return chart
}

func main() {
	svg, _ := fibChart(10).RenderString()
	fmt.Println(svg)

	fmt.Println("Serving bar charts at http://localhost:1343")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Bar Chart</title></head>
<body style="font-family: system-ui; padding: 2em; max-width: 900px; margin: 0 auto;">
<h1>Bar Chart</h1>
<h2>Single series</h2>
`)
		fibChart(10).Render(w)
		fmt.Fprint(w, "<h2>Grouped series</h2>\n")
		groupedChart().Render(w)
		fmt.Fprint(w, "\n</body></html>")
	})
	http.ListenAndServe(":1343", nil)
}
