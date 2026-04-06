//go:build ignore

package main

import (
	"fmt"
	"os"
	"time"

	"codeberg.org/hum3/gogal"
)

func main() {
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
	var temp, humidity []gogal.DataPoint
	for _, r := range tempR {
		t := time.Date(2024, 6, 15, r.hour, r.min, 0, 0, time.UTC)
		temp = append(temp, gogal.DataPoint{Time: t, Y: r.val, Label: fmt.Sprintf("%.1f\u00b0C", r.val)})
	}
	for _, r := range humR {
		t := time.Date(2024, 6, 15, r.hour, r.min, 0, 0, time.UTC)
		humidity = append(humidity, gogal.DataPoint{Time: t, Y: r.val, Label: fmt.Sprintf("%.0f%%", r.val)})
	}

	// Temporal
	tc := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithAxisMode(gogal.Temporal),
		gogal.WithTitle("Weather \u2014 Temporal Axis"),
		gogal.WithGrid(true),
		gogal.WithSmooth(true),
		gogal.WithTimeFormat("15:04"),
	)
	tc.Add("Temperature (\u00b0C)", temp)
	tc.Add("Humidity (%)", humidity)

	f, _ := os.Create("../../../docs/03_multi_series/03_temporal.svg")
	tc.Render(f)
	f.Close()

	// Ordinal
	oc := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithAxisMode(gogal.Ordinal),
		gogal.WithTitle("Weather \u2014 Ordinal Axis"),
		gogal.WithGrid(true),
		gogal.WithSmooth(true),
	)
	oc.Add("Temperature (\u00b0C)", temp)
	oc.Add("Humidity (%)", humidity)

	f, _ = os.Create("../../../docs/03_multi_series/03_ordinal.svg")
	oc.Render(f)
	f.Close()
}
