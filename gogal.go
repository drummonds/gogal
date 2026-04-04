// Package gogal provides pure SVG chart generation for Go.
//
// Charts are rendered server-side as pure SVG with CSS-only interactivity
// (hover tooltips, transitions) and optional HTMX integration (legend
// toggling, axis switching, live updates via SSE).
//
// The architecture separates layout computation from SVG rendering,
// enabling comprehensive testing without SVG parsing.
package gogal

import (
	"bytes"
	"io"
	"time"
)

// Chart is the main type for building and rendering charts.
type Chart struct {
	config    ChartConfig
	series    []Series
	chartType string
}

// NewLineChart creates a new line chart with the given options.
func NewLineChart(opts ...Option) *Chart {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Chart{config: cfg, chartType: "line"}
}

// Add adds a data series to the chart.
func (c *Chart) Add(name string, points []DataPoint) *Chart {
	c.series = append(c.series, Series{
		Name:   name,
		Points: points,
	})
	return c
}

// AddXY adds a series from parallel X and Y slices.
func (c *Chart) AddXY(name string, xs, ys []float64) *Chart {
	n := min(len(xs), len(ys))
	points := make([]DataPoint, n)
	for i := 0; i < n; i++ {
		points[i] = DataPoint{X: xs[i], Y: ys[i]}
	}
	return c.Add(name, points)
}

// AddTimeSeries adds a series with time-based X values.
func (c *Chart) AddTimeSeries(name string, times []time.Time, values []float64) *Chart {
	n := min(len(times), len(values))
	points := make([]DataPoint, n)
	for i := 0; i < n; i++ {
		points[i] = DataPoint{
			X:    float64(times[i].Unix()),
			Y:    values[i],
			Time: times[i],
		}
	}
	return c.Add(name, points)
}

// Layout computes the chart layout without rendering SVG.
// This is useful for testing layout logic.
func (c *Chart) Layout() *LayoutResult {
	return computeLayout(&c.config, c.series, c.chartType)
}

// Render writes the chart as SVG to w.
func (c *Chart) Render(w io.Writer) error {
	layout := c.Layout()
	return render(w, layout)
}

// RenderString returns the chart SVG as a string.
func (c *Chart) RenderString() (string, error) {
	var buf bytes.Buffer
	err := c.Render(&buf)
	return buf.String(), err
}
