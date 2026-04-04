package gogal

import "time"

// DataPoint represents a single data point in a series.
type DataPoint struct {
	X     float64   // quantitative X value (or Unix timestamp for temporal)
	Y     float64   // quantitative Y value
	Label string    // optional per-point label (for tooltips)
	Color string    // optional per-point color override
	Time  time.Time // temporal X value (used when AxisMode is Temporal)
}

// Series represents a named collection of data points.
type Series struct {
	Name   string
	Points []DataPoint
	Color  string // series color (from theme if empty)
	Hidden bool   // whether series is currently hidden (for HTMX toggle)
}
