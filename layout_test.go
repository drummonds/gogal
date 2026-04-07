package gogal

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSparklineLayout(t *testing.T) {
	chart := NewLineChart(WithVariant(Sparkline))
	chart.Add("temp", []DataPoint{
		{X: 0, Y: 10},
		{X: 1, Y: 20},
		{X: 2, Y: 15},
		{X: 3, Y: 25},
	})

	layout := chart.Layout()

	// Sparkline: plot area fills entire SVG
	if layout.PlotArea.X != 0 || layout.PlotArea.Y != 0 {
		t.Errorf("sparkline plot area should start at origin, got (%v,%v)",
			layout.PlotArea.X, layout.PlotArea.Y)
	}

	// No axes, no legend
	if layout.XAxis != nil {
		t.Error("sparkline should not have X axis")
	}
	if layout.YAxis != nil {
		t.Error("sparkline should not have Y axis")
	}
	if layout.Legend != nil {
		t.Error("sparkline should not have legend")
	}

	// Should have one series with a path
	if len(layout.Series) != 1 {
		t.Fatalf("expected 1 series, got %d", len(layout.Series))
	}
	if layout.Series[0].Path == "" {
		t.Error("series should have a path")
	}
	if !strings.HasPrefix(layout.Series[0].Path, "M") {
		t.Errorf("path should start with M, got %q", layout.Series[0].Path[:10])
	}
}

func TestSparklineRender(t *testing.T) {
	chart := NewLineChart(WithVariant(Sparkline))
	chart.Add("data", []DataPoint{
		{X: 0, Y: 5},
		{X: 1, Y: 10},
		{X: 2, Y: 3},
	})

	var buf bytes.Buffer
	err := chart.Render(&buf)
	if err != nil {
		t.Fatal(err)
	}

	svg := buf.String()

	// Basic SVG structure
	if !strings.Contains(svg, "<svg") {
		t.Error("output should contain <svg")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("output should contain </svg>")
	}
	if !strings.Contains(svg, "viewBox") {
		t.Error("output should contain viewBox")
	}
	if !strings.Contains(svg, "<path") {
		t.Error("output should contain <path")
	}
	// Sparkline: no tooltip, no axis, no title
	if strings.Contains(svg, "tooltip") {
		t.Error("sparkline should not contain tooltip")
	}
	if strings.Contains(svg, "x-axis") {
		t.Error("sparkline should not contain x-axis")
	}
}

func TestStaticLayout(t *testing.T) {
	chart := NewLineChart(
		WithTitle("Test Chart"),
		WithSize(800, 400),
	)
	chart.Add("A", []DataPoint{
		{X: 0, Y: 10},
		{X: 1, Y: 20},
		{X: 2, Y: 30},
	})
	chart.Add("B", []DataPoint{
		{X: 0, Y: 5},
		{X: 1, Y: 15},
		{X: 2, Y: 25},
	})

	layout := chart.Layout()

	// Should have title
	if layout.Title == nil {
		t.Fatal("expected title")
	}
	if layout.Title.Text != "Test Chart" {
		t.Errorf("title = %q, want %q", layout.Title.Text, "Test Chart")
	}

	// Should have axes
	if layout.XAxis == nil {
		t.Fatal("expected X axis")
	}
	if layout.YAxis == nil {
		t.Fatal("expected Y axis")
	}

	// Should have legend (2 series)
	if layout.Legend == nil {
		t.Fatal("expected legend")
	}
	if len(layout.Legend.Entries) != 2 {
		t.Errorf("expected 2 legend entries, got %d", len(layout.Legend.Entries))
	}

	// Series should have paths and points
	if len(layout.Series) != 2 {
		t.Fatalf("expected 2 series, got %d", len(layout.Series))
	}
	for i, s := range layout.Series {
		if s.Path == "" {
			t.Errorf("series[%d] should have a path", i)
		}
		if len(s.Points) != 3 {
			t.Errorf("series[%d] should have 3 points, got %d", i, len(s.Points))
		}
	}

	// Points should be within plot area
	for i, s := range layout.Series {
		for j, p := range s.Points {
			if p.X < layout.PlotArea.X || p.X > layout.PlotArea.X+layout.PlotArea.Width {
				t.Errorf("series[%d].point[%d].X = %v, outside plot area", i, j, p.X)
			}
			if p.Y < layout.PlotArea.Y || p.Y > layout.PlotArea.Y+layout.PlotArea.Height {
				t.Errorf("series[%d].point[%d].Y = %v, outside plot area", i, j, p.Y)
			}
		}
	}
}

func TestHiddenSeries(t *testing.T) {
	chart := NewLineChart(
		WithHiddenSeries("B"),
	)
	chart.Add("A", []DataPoint{{X: 0, Y: 10}, {X: 1, Y: 20}})
	chart.Add("B", []DataPoint{{X: 0, Y: 5}, {X: 1, Y: 15}})

	layout := chart.Layout()

	if len(layout.Series) != 1 {
		t.Fatalf("expected 1 visible series, got %d", len(layout.Series))
	}
	if layout.Series[0].Name != "A" {
		t.Errorf("visible series should be A, got %q", layout.Series[0].Name)
	}
}

func TestStaticRender(t *testing.T) {
	chart := NewLineChart(
		WithTitle("Temperature"),
		WithSize(800, 400),
		WithTooltips(true),
	)
	chart.Add("Sensor", []DataPoint{
		{X: 0, Y: 20},
		{X: 1, Y: 22},
		{X: 2, Y: 19},
	})

	var buf bytes.Buffer
	err := chart.Render(&buf)
	if err != nil {
		t.Fatal(err)
	}

	svg := buf.String()

	// Should have all expected elements
	expects := []string{
		"<svg",
		"viewBox",
		"<title",
		"Temperature",
		"<style",
		"tooltip",
		"x-axis",
		"y-axis",
		"<path",
		"<circle",
	}
	for _, e := range expects {
		if !strings.Contains(svg, e) {
			t.Errorf("static chart should contain %q", e)
		}
	}
}

func TestSmoothPath(t *testing.T) {
	chart := NewLineChart(WithVariant(Sparkline), WithSmooth(true))
	chart.Add("data", []DataPoint{
		{X: 0, Y: 0},
		{X: 1, Y: 10},
		{X: 2, Y: 5},
		{X: 3, Y: 15},
	})

	layout := chart.Layout()
	path := layout.Series[0].Path

	// Smooth path should use C (cubic bezier) commands
	if !strings.Contains(path, "C") {
		t.Errorf("smooth path should contain C commands, got %q", path)
	}
}

func TestEmptyChart(t *testing.T) {
	chart := NewLineChart()
	layout := chart.Layout()

	// Should not panic, should have no series
	if len(layout.Series) != 0 {
		t.Errorf("empty chart should have 0 series, got %d", len(layout.Series))
	}

	// Render should not error
	var buf bytes.Buffer
	err := chart.Render(&buf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRenderString(t *testing.T) {
	chart := NewLineChart(WithVariant(Sparkline))
	chart.Add("data", []DataPoint{{X: 0, Y: 1}, {X: 1, Y: 2}})

	svg, err := chart.RenderString()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(svg, "<svg") {
		t.Error("RenderString should return SVG")
	}
}

func makeTimePoints(hours []float64, ys []float64) []DataPoint {
	pts := make([]DataPoint, len(hours))
	for i, h := range hours {
		hr := int(h)
		min := int((h - float64(hr)) * 60)
		pts[i] = DataPoint{
			Time:  time.Date(2024, 6, 15, hr, min, 0, 0, time.UTC),
			Y:     ys[i],
			Label: "tooltip",
		}
	}
	return pts
}

func TestOrdinalLayout_PointsWithinPlotArea(t *testing.T) {
	temp := makeTimePoints(
		[]float64{0, 1, 2, 3, 4},
		[]float64{10, 12, 14, 11, 13},
	)
	hum := makeTimePoints(
		[]float64{0, 2, 4},
		[]float64{80, 70, 75},
	)

	chart := NewLineChart(
		WithVariant(Static),
		WithAxisMode(Ordinal),
		WithSize(800, 400),
		WithTimeFormat("15:04"),
	)
	chart.Add("Temp", temp)
	chart.Add("Humidity", hum)

	layout := chart.Layout()

	for i, s := range layout.Series {
		for j, p := range s.Points {
			if p.X < layout.PlotArea.X-1 || p.X > layout.PlotArea.X+layout.PlotArea.Width+1 {
				t.Errorf("series[%d].point[%d].X = %.1f, outside plot area [%.1f, %.1f]",
					i, j, p.X, layout.PlotArea.X, layout.PlotArea.X+layout.PlotArea.Width)
			}
			if p.Y < layout.PlotArea.Y-1 || p.Y > layout.PlotArea.Y+layout.PlotArea.Height+1 {
				t.Errorf("series[%d].point[%d].Y = %.1f, outside plot area [%.1f, %.1f]",
					i, j, p.Y, layout.PlotArea.Y, layout.PlotArea.Y+layout.PlotArea.Height)
			}
		}
	}
}

func TestOrdinalLayout_XAxisLabels(t *testing.T) {
	temp := makeTimePoints(
		[]float64{0, 1, 2},
		[]float64{10, 12, 14},
	)

	chart := NewLineChart(
		WithVariant(Static),
		WithAxisMode(Ordinal),
		WithTimeFormat("15:04"),
	)
	chart.Add("Temp", temp)

	layout := chart.Layout()

	if layout.XAxis == nil {
		t.Fatal("expected X axis")
	}
	if len(layout.XAxis.Ticks) != 3 {
		t.Fatalf("expected 3 X ticks, got %d", len(layout.XAxis.Ticks))
	}

	wantLabels := []string{"00:00", "01:00", "02:00"}
	for i, tick := range layout.XAxis.Ticks {
		if tick.Label != wantLabels[i] {
			t.Errorf("tick[%d].Label = %q, want %q", i, tick.Label, wantLabels[i])
		}
	}
}

func TestOrdinalLayout_LabelCount(t *testing.T) {
	long := makeTimePoints(
		[]float64{0, 1, 2, 3, 4},
		[]float64{10, 12, 14, 11, 13},
	)
	short := makeTimePoints(
		[]float64{0, 2},
		[]float64{80, 70},
	)

	chart := NewLineChart(
		WithVariant(Static),
		WithAxisMode(Ordinal),
		WithTimeFormat("15:04"),
	)
	chart.Add("Long", long)
	chart.Add("Short", short)

	layout := chart.Layout()

	if layout.XAxis == nil {
		t.Fatal("expected X axis")
	}
	// Tick count should equal the longest series length.
	if len(layout.XAxis.Ticks) != 5 {
		t.Errorf("expected 5 X ticks (longest series), got %d", len(layout.XAxis.Ticks))
	}
}

func TestYFormatApplied(t *testing.T) {
	chart := NewLineChart(
		WithVariant(Static),
		WithSize(800, 400),
		WithYFormat("%.0f"),
	)
	chart.Add("A", []DataPoint{
		{X: 0, Y: 10},
		{X: 1, Y: 20},
		{X: 2, Y: 30},
	})

	layout := chart.Layout()

	if layout.YAxis == nil {
		t.Fatal("expected Y axis")
	}
	for _, tick := range layout.YAxis.Ticks {
		if strings.Contains(tick.Label, ".") {
			t.Errorf("Y tick label %q should have no decimal with format %%.0f", tick.Label)
		}
	}
}

func TestTemporalLayout_XAxisLabels(t *testing.T) {
	base := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	pts := make([]DataPoint, 7)
	for i := range pts {
		pts[i] = DataPoint{
			Time: base.Add(time.Duration(i*4) * time.Hour),
			Y:    float64(10 + i),
		}
	}

	chart := NewLineChart(
		WithVariant(Static),
		WithSize(800, 400),
		WithTimeFormat("15:04"),
	)
	chart.Add("data", pts)

	layout := chart.Layout()

	if layout.XAxis == nil {
		t.Fatal("expected X axis")
	}

	// Tick labels should be time-formatted, not scientific notation.
	for _, tick := range layout.XAxis.Ticks {
		if strings.Contains(tick.Label, "e+") {
			t.Errorf("temporal tick label %q should not be scientific notation", tick.Label)
		}
		if !strings.Contains(tick.Label, ":") {
			t.Errorf("temporal tick label %q should contain ':', expected HH:MM", tick.Label)
		}
	}
}
