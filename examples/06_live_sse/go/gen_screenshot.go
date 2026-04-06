//go:build ignore

package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"codeberg.org/hum3/gogal"
)

func main() {
	var points []gogal.DataPoint
	base := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	for i := range 30 {
		t := base.Add(time.Duration(i*2) * time.Second)
		y := 50 + 20*math.Sin(float64(i)*0.3) + rand.Float64()*10 - 5
		points = append(points, gogal.DataPoint{
			Time:  t,
			Y:     y,
			Label: fmt.Sprintf("%.1f", y),
		})
	}

	chart := gogal.NewLineChart(
		gogal.WithVariant(gogal.Static),
		gogal.WithTitle("Live Sensor Data"),
		gogal.WithGrid(true),
		gogal.WithSmooth(true),
		gogal.WithTooltips(true),
		gogal.WithTimeFormat("15:04:05"),
		gogal.WithYFormat("%.0f"),
	)
	chart.Add("Sensor", points)

	f, _ := os.Create("../../../docs/06_live_sse/06_live.svg")
	chart.Render(f)
	f.Close()
}
