package gogal

import (
	"math"
	"strings"
	"testing"
)

var fib = []float64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
var fibLabels = []string{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10"}

func TestBarLayout_SingleSeries(t *testing.T) {
	chart := NewBarChart(WithTitle("Fibonacci"), WithSize(500, 220))
	chart.AddCategories("Fibonacci", fibLabels, fib)

	layout := chart.Layout()

	if len(layout.Series) != 1 {
		t.Fatalf("expected 1 series, got %d", len(layout.Series))
	}
	bars := layout.Series[0].Bars
	if len(bars) != len(fib) {
		t.Fatalf("expected %d bars, got %d", len(fib), len(bars))
	}
	if layout.Series[0].Path != "" {
		t.Error("bar series should not have a line path")
	}

	// Category labels on the X axis.
	var ticks []string
	for _, tick := range layout.XAxis.Ticks {
		ticks = append(ticks, tick.Label)
	}
	if got := strings.Join(ticks, " "); got != strings.Join(fibLabels, " ") {
		t.Errorf("x ticks = %q", got)
	}

	// Bars stand on the baseline (Y=0 at the bottom of the plot area).
	bottom := layout.PlotArea.Y + layout.PlotArea.Height
	for i, b := range bars {
		if math.Abs(b.Y+b.Height-bottom) > 0.01 {
			t.Errorf("bar %d bottom = %v, want baseline %v", i, b.Y+b.Height, bottom)
		}
		if b.X < layout.PlotArea.X || b.X+b.Width > layout.PlotArea.X+layout.PlotArea.Width+0.01 {
			t.Errorf("bar %d outside plot area", i)
		}
	}

	// Heights are proportional to value: F10 (55) is 55x F1 (1).
	if ratio := bars[9].Height / bars[0].Height; math.Abs(ratio-55) > 0.01 {
		t.Errorf("height ratio F10/F1 = %v, want 55", ratio)
	}
	// Tallest bar leaves headroom for its value label.
	if bars[9].Y <= layout.PlotArea.Y {
		t.Errorf("tallest bar top %v should sit below plot top %v", bars[9].Y, layout.PlotArea.Y)
	}

	// Equal widths, narrower than the band, left to right.
	band := layout.PlotArea.Width / float64(len(fib))
	for i, b := range bars {
		if math.Abs(b.Width-bars[0].Width) > 0.01 {
			t.Errorf("bar %d width %v differs from %v", i, b.Width, bars[0].Width)
		}
		if b.Width >= band {
			t.Errorf("bar %d width %v should be less than band %v", i, b.Width, band)
		}
		if i > 0 && b.X <= bars[i-1].X {
			t.Errorf("bar %d not to the right of bar %d", i, i-1)
		}
	}

	// Value labels and categories.
	if bars[9].Value != "55" || bars[9].Category != "F10" {
		t.Errorf("bar 9 = %+v, want Value 55, Category F10", bars[9])
	}
}

func TestBarLayout_GroupedSeries(t *testing.T) {
	chart := NewBarChart()
	chart.AddCategories("2024", []string{"Q1", "Q2", "Q3"}, []float64{10, 20, 30})
	chart.AddCategories("2025", []string{"Q1", "Q2", "Q3"}, []float64{15, 25, 35})

	layout := chart.Layout()
	if len(layout.Series) != 2 {
		t.Fatalf("expected 2 series, got %d", len(layout.Series))
	}
	a, b := layout.Series[0].Bars, layout.Series[1].Bars
	if len(a) != 3 || len(b) != 3 {
		t.Fatalf("expected 3 bars per series, got %d and %d", len(a), len(b))
	}
	for i := range a {
		// Side by side within the band, no overlap.
		if a[i].X+a[i].Width > b[i].X+0.01 {
			t.Errorf("category %d: series bars overlap (%v..%v vs %v)", i, a[i].X, a[i].X+a[i].Width, b[i].X)
		}
		if i < 2 && b[i].X+b[i].Width > a[i+1].X+0.01 {
			t.Errorf("category %d spills into category %d", i, i+1)
		}
	}
	if layout.Legend == nil || len(layout.Legend.Entries) != 2 {
		t.Error("grouped bar chart should have a two-entry legend")
	}
}

func TestBarLayout_NegativeValuesHangFromBaseline(t *testing.T) {
	chart := NewBarChart(WithValueLabels(false))
	chart.AddCategories("delta", []string{"a", "b"}, []float64{-4, 8})

	layout := chart.Layout()
	bars := layout.Series[0].Bars
	if bars[0].Value != "" {
		t.Errorf("value labels disabled, got %q", bars[0].Value)
	}
	// Negative bar starts at the baseline and extends downward; positive ends at it.
	baseline := bars[1].Y + bars[1].Height
	if math.Abs(bars[0].Y-baseline) > 0.01 {
		t.Errorf("negative bar top %v should equal baseline %v", bars[0].Y, baseline)
	}
	if ratio := bars[1].Height / bars[0].Height; math.Abs(ratio-2) > 0.01 {
		t.Errorf("height ratio = %v, want 2", ratio)
	}
}

func TestBarSparkline(t *testing.T) {
	chart := NewBarChart(WithVariant(Sparkline))
	chart.AddCategories("s", []string{"a", "b", "c", "d"}, []float64{1, 3, 2, 4})

	layout := chart.Layout()
	if layout.XAxis != nil || layout.YAxis != nil || layout.Legend != nil {
		t.Error("sparkline bars should have no axes or legend")
	}
	bars := layout.Series[0].Bars
	if len(bars) != 4 {
		t.Fatalf("expected 4 bars, got %d", len(bars))
	}
	last := bars[3]
	if last.X+last.Width > 100.01 || last.Y < -0.01 || last.Y+last.Height > 20.01 {
		t.Errorf("sparkline bars should fill the 100x20 box, got %+v", last)
	}
}

func TestBarRender(t *testing.T) {
	chart := NewBarChart(WithTitle("Fibonacci (n=10)"), WithSize(500, 220))
	chart.AddCategories("Fibonacci", fibLabels, fib)

	svg, err := chart.RenderString()
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(svg, `class="bar"`); n != 10 {
		t.Errorf("expected 10 bar rects, got %d", n)
	}
	if strings.Contains(svg, "<circle") {
		t.Error("bar chart should not draw data-point circles")
	}
	if !strings.Contains(svg, ">55</text>") {
		t.Error("expected value label 55 above the tallest bar")
	}
	if !strings.Contains(svg, ">F10</text>") {
		t.Error("expected category label F10 on the x axis")
	}
	if !strings.Contains(svg, "Fibonacci (n=10)") {
		t.Error("expected title")
	}
}
