package main

import (
	"time"

	"codeberg.org/hum3/gogal"
	"codeberg.org/hum3/lofigui"
)

type axisExample struct {
	title  string
	points []gogal.DataPoint
	opts   []gogal.Option
}

// model renders all axis format examples as SVG charts.
func model(app *lofigui.App) {
	lofigui.HTML(`<h2>Axis Format Gallery</h2>
<p>Visual test of axis tick formatting across different data ranges.</p>
<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1.5em;">`)

	for _, ex := range examples() {
		opts := []gogal.Option{
			gogal.WithVariant(gogal.Static),
			gogal.WithSize(500, 250),
			gogal.WithGrid(true),
			gogal.WithTitle(ex.title),
		}
		opts = append(opts, ex.opts...)

		chart := gogal.NewLineChart(opts...)
		chart.Add("data", ex.points)
		svg, _ := chart.RenderString()
		lofigui.HTML(`<div>` + svg + `</div>`)
	}
	lofigui.HTML(`</div>`)
}

func examples() []axisExample {
	return []axisExample{
		smallIntegers(),
		decimals(),
		hundreds(),
		largeRange(),
		scientificSmall(),
		scientificLarge(),
		negativeRange(),
		narrowRange(),
		temporalHours(),
		temporalDays(),
		temporalMonths(),
		temporalYears(),
		ordinalDays(),
		ordinalCategories(),
	}
}

func smallIntegers() axisExample {
	pts := xyPoints([]float64{0, 2, 4, 6, 8, 10}, []float64{1, 3, 2, 5, 4, 6})
	return axisExample{"Small integers (0-10)", pts, nil}
}

func decimals() axisExample {
	pts := xyPoints(
		[]float64{0, 0.2, 0.4, 0.6, 0.8, 1.0},
		[]float64{0.1, 0.3, 0.2, 0.5, 0.4, 0.6},
	)
	return axisExample{"Decimals (0-1)", pts, nil}
}

func hundreds() axisExample {
	pts := xyPoints(
		[]float64{0, 100, 200, 300, 400, 500},
		[]float64{50, 150, 120, 280, 350, 420},
	)
	return axisExample{"Hundreds (0-500)", pts, nil}
}

func largeRange() axisExample {
	pts := xyPoints(
		[]float64{0, 20000, 40000, 60000, 80000, 100000},
		[]float64{10000, 30000, 25000, 55000, 70000, 90000},
	)
	return axisExample{"Large (0-100k)", pts, nil}
}

func scientificSmall() axisExample {
	pts := xyPoints(
		[]float64{1e-6, 2e-6, 5e-6, 1e-5, 5e-5, 1e-4, 5e-4, 1e-3},
		[]float64{1e-6, 3e-6, 2e-6, 8e-6, 4e-5, 7e-5, 3e-4, 8e-4},
	)
	return axisExample{"Scientific small (1e-6 to 1e-3)", pts, nil}
}

func scientificLarge() axisExample {
	pts := xyPoints(
		[]float64{1e6, 1e7, 1e8, 5e8, 1e9},
		[]float64{5e6, 2e7, 8e7, 3e8, 7e8},
	)
	return axisExample{"Scientific large (1e6 to 1e9)", pts, nil}
}

func negativeRange() axisExample {
	pts := xyPoints(
		[]float64{-50, -30, -10, 10, 30, 50},
		[]float64{-40, -15, 5, 20, -5, 35},
	)
	return axisExample{"Negative range (-50 to 50)", pts, nil}
}

func narrowRange() axisExample {
	pts := xyPoints(
		[]float64{10.0, 10.1, 10.2, 10.3, 10.4, 10.5},
		[]float64{10.05, 10.15, 10.1, 10.25, 10.35, 10.45},
	)
	return axisExample{"Narrow (10.0-10.5)", pts, nil}
}

func temporalHours() axisExample {
	base := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	pts := make([]gogal.DataPoint, 7)
	for i := range pts {
		t := base.Add(time.Duration(i*4) * time.Hour)
		pts[i] = gogal.DataPoint{Time: t, Y: float64(10 + i*3)}
	}
	return axisExample{
		"Temporal: hours (24h)",
		pts,
		[]gogal.Option{gogal.WithTimeFormat("15:04")},
	}
}

func temporalDays() axisExample {
	base := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	pts := make([]gogal.DataPoint, 7)
	for i := range pts {
		t := base.Add(time.Duration(i*5) * 24 * time.Hour)
		pts[i] = gogal.DataPoint{Time: t, Y: float64(20 + i*2)}
	}
	return axisExample{
		"Temporal: days (30d)",
		pts,
		[]gogal.Option{gogal.WithTimeFormat("Jan 02")},
	}
}

func temporalMonths() axisExample {
	pts := make([]gogal.DataPoint, 12)
	for i := range pts {
		t := time.Date(2024, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
		pts[i] = gogal.DataPoint{Time: t, Y: float64(5 + i)}
	}
	return axisExample{
		"Temporal: months (12 months)",
		pts,
		[]gogal.Option{gogal.WithTimeFormat("Jan")},
	}
}

func temporalYears() axisExample {
	pts := make([]gogal.DataPoint, 6)
	for i := range pts {
		t := time.Date(2020+i*2, 1, 1, 0, 0, 0, 0, time.UTC)
		pts[i] = gogal.DataPoint{Time: t, Y: float64(100 + i*50)}
	}
	return axisExample{
		"Temporal: years (2020-2030)",
		pts,
		[]gogal.Option{gogal.WithTimeFormat("2006")},
	}
}

func ordinalDays() axisExample {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	pts := make([]gogal.DataPoint, len(days))
	for i, d := range days {
		pts[i] = gogal.DataPoint{X: float64(i), Y: float64(10 + i*3), Label: d}
	}
	return axisExample{
		"Ordinal: weekdays",
		pts,
		[]gogal.Option{gogal.WithAxisMode(gogal.Ordinal)},
	}
}

func ordinalCategories() axisExample {
	cats := []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"}
	pts := make([]gogal.DataPoint, len(cats))
	for i, c := range cats {
		pts[i] = gogal.DataPoint{X: float64(i), Y: float64(5 + i*2), Label: c}
	}
	return axisExample{
		"Ordinal: categories",
		pts,
		[]gogal.Option{gogal.WithAxisMode(gogal.Ordinal)},
	}
}

func xyPoints(xs, ys []float64) []gogal.DataPoint {
	n := min(len(xs), len(ys))
	pts := make([]gogal.DataPoint, n)
	for i := range n {
		pts[i] = gogal.DataPoint{X: xs[i], Y: ys[i]}
	}
	return pts
}
