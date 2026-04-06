//go:build ignore

package main

import (
	"fmt"
	"os"
	"time"

	"codeberg.org/hum3/gogal"
)

func main() {
	// Same irregular measurements as main.go
	meas := []struct {
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
	var points []gogal.DataPoint
	for _, m := range meas {
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
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
		gogal.WithYFormat("%.0f"),
	)
	chart.Add("Temperature", points)

	f, _ := os.Create("../../../docs/02_static_line/02_static_line.svg")
	chart.Render(f)
	f.Close()
}
