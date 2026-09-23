//go:build ignore

package main

import (
	"os"

	"git.bytestone.uk/hum3/gogal"
)

func main() {
	build := func(opts ...gogal.Option) *gogal.Chart {
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
	write("../../../docs/10_pie_chart/10_theme.svg", build())
	write("../../../docs/10_pie_chart/10_colours.svg",
		build(gogal.WithColors("#3298dc", "#48c774", "#f14668", "#ffdd57", "#b86bff")))
}

func write(path string, chart *gogal.Chart) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	chart.Render(f)
}
