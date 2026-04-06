//go:build ignore

package main

import (
	"fmt"
	"os"
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

func main() {
	temp := toPoints([]reading{
		{0, 0, 11.2}, {1, 20, 10.4}, {2, 50, 9.7}, {4, 10, 8.9},
		{5, 30, 9.1}, {6, 45, 10.8}, {7, 15, 12.5}, {8, 40, 14.8},
		{9, 50, 17.2}, {10, 30, 19.0}, {11, 10, 20.6}, {12, 25, 22.1},
		{13, 40, 23.0}, {14, 15, 22.5}, {15, 30, 21.3}, {16, 50, 19.7},
		{18, 0, 17.4}, {19, 20, 15.6}, {20, 40, 14.0}, {22, 0, 12.3},
		{23, 30, 11.5},
	}, "%.1f\u00b0C")
	humidity := toPoints([]reading{
		{0, 0, 82.0}, {2, 30, 85.0}, {4, 45, 88.0}, {6, 30, 84.0},
		{8, 15, 72.0}, {10, 0, 58.0}, {11, 30, 50.0}, {13, 0, 45.0},
		{15, 0, 48.0}, {17, 0, 55.0}, {19, 0, 65.0}, {21, 30, 76.0},
		{23, 0, 80.0},
	}, "%.0f%%")
	wind := toPoints([]reading{
		{0, 30, 6.2}, {3, 0, 4.8}, {5, 15, 5.5}, {7, 0, 8.3},
		{9, 30, 12.1}, {11, 0, 15.4}, {12, 45, 18.2}, {14, 30, 16.7},
		{16, 0, 13.9}, {18, 30, 10.5}, {20, 0, 7.8}, {22, 30, 5.9},
	}, "%.1f km/h")

	// All series visible
	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithTitle("Weather Station"),
		gogal.WithGrid(true),
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
	)
	chart.Add("Temperature", temp)
	chart.Add("Humidity", humidity)
	chart.Add("Wind Speed", wind)

	f, _ := os.Create("../../../docs/05_htmx/05_all.svg")
	chart.Render(f)
	f.Close()

	// Humidity hidden
	chart2 := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithTitle("Weather Station"),
		gogal.WithGrid(true),
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
		gogal.WithHiddenSeries("Humidity"),
	)
	chart2.Add("Temperature", temp)
	chart2.Add("Humidity", humidity)
	chart2.Add("Wind Speed", wind)

	f, _ = os.Create("../../../docs/05_htmx/05_toggled.svg")
	chart2.Render(f)
	f.Close()
}
