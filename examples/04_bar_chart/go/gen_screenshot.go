//go:build ignore

package main

import (
	"os"

	"git.bytestone.uk/hum3/gogal"
)

func main() {
	fib := gogal.NewBarChart(
		gogal.WithTitle("Fibonacci (n=10)"),
		gogal.WithSize(500, 220),
		gogal.WithYFormat("%.0f"),
	)
	fib.AddCategories("Fibonacci",
		[]string{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10"},
		[]float64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55})
	write("../../../docs/04_bar_chart/04_fibonacci.svg", fib)

	grouped := gogal.NewBarChart(
		gogal.WithTitle("Quarterly sales"),
		gogal.WithSize(500, 260),
		gogal.WithYTitle("k£"),
		gogal.WithYFormat("%.0f"),
	)
	quarters := []string{"Q1", "Q2", "Q3", "Q4"}
	grouped.AddCategories("2024", quarters, []float64{120, 95, 140, 160})
	grouped.AddCategories("2025", quarters, []float64{135, 110, 150, 185})
	write("../../../docs/04_bar_chart/04_grouped.svg", grouped)
}

func write(path string, chart *gogal.Chart) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	chart.Render(f)
}
