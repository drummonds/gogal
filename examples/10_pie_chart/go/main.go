//go:build !js

// Example 10: Pie Chart
//
// Proportional slices with a legend of labels and percentages, in theme
// colours and with an explicit palette. Serves on http://localhost:1349.
package main

import (
	"fmt"
	"net/http"

	"git.bytestone.uk/hum3/gogal"
)

func languageChart(opts ...gogal.Option) *gogal.Chart {
	chart := gogal.NewPieChart(append([]gogal.Option{
		gogal.WithTitle("Language Popularity"),
		gogal.WithSize(450, 300),
	}, opts...)...)
	chart.AddSlice("Go", 35)
	chart.AddSlice("Python", 28)
	chart.AddSlice("Rust", 18)
	chart.AddSlice("JS", 12)
	chart.AddSlice("Other", 7)
	return chart
}

func main() {
	svg, _ := languageChart().RenderString()
	fmt.Println(svg)

	fmt.Println("Serving pie charts at http://localhost:1349")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Pie Chart</title></head>
<body style="font-family: system-ui; padding: 2em; max-width: 900px; margin: 0 auto;">
<h1>Pie Chart</h1>
<h2>Theme palette</h2>
`)
		languageChart().Render(w)
		fmt.Fprint(w, "<h2>Explicit colours</h2>\n")
		languageChart(gogal.WithColors("#3298dc", "#48c774", "#f14668", "#ffdd57", "#b86bff")).Render(w)
		fmt.Fprint(w, "\n</body></html>")
	})
	http.ListenAndServe(":1349", nil)
}
