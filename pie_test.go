package gogal

import (
	"math"
	"strings"
	"testing"
)

func languagePie(opts ...Option) *Chart {
	chart := NewPieChart(append([]Option{WithTitle("Language Popularity"), WithSize(450, 300)}, opts...)...)
	chart.AddSlice("Go", 35)
	chart.AddSlice("Python", 28)
	chart.AddSlice("Rust", 18)
	chart.AddSlice("JS", 12)
	chart.AddSlice("Other", 7)
	return chart
}

func TestPieLayout_Slices(t *testing.T) {
	layout := languagePie().Layout()

	if len(layout.Slices) != 5 {
		t.Fatalf("expected 5 slices, got %d", len(layout.Slices))
	}
	if layout.XAxis != nil || layout.YAxis != nil {
		t.Error("pie chart should have no axes")
	}

	total := 0.0
	for _, s := range layout.Slices {
		total += s.Percent
		if !strings.HasPrefix(s.Path, "M") || !strings.Contains(s.Path, "A") {
			t.Errorf("slice %s path %q should be an arc sector", s.Label, s.Path)
		}
		if s.Color == "" {
			t.Errorf("slice %s has no colour", s.Label)
		}
	}
	if math.Abs(total-100) > 0.01 {
		t.Errorf("percentages sum to %v, want 100", total)
	}
	if got := layout.Slices[0]; got.Label != "Go" || math.Abs(got.Percent-35) > 0.01 {
		t.Errorf("first slice = %+v, want Go at 35%%", got)
	}

	// Legend lists label and percentage, one entry per slice, stacked vertically.
	if layout.Legend == nil || len(layout.Legend.Entries) != 5 {
		t.Fatal("expected a five-entry legend")
	}
	if layout.Legend.Entries[0].Name != "Go 35.0%" {
		t.Errorf("legend entry = %q, want %q", layout.Legend.Entries[0].Name, "Go 35.0%")
	}
	e := layout.Legend.Entries
	if e[1].Y <= e[0].Y || e[1].X != e[0].X {
		t.Errorf("pie legend entries should stack vertically: %+v %+v", e[0], e[1])
	}
	// Legend sits to the right of the pie.
	if e[0].X <= layout.Pie.CX+layout.Pie.R {
		t.Errorf("legend x %v should be right of the pie (cx %v r %v)", e[0].X, layout.Pie.CX, layout.Pie.R)
	}
	if layout.Pie.R <= 0 {
		t.Error("pie radius should be positive")
	}
}

func TestPieLayout_ColoursPerSlice(t *testing.T) {
	colours := []string{"#3298dc", "#48c774", "#f14668", "#ffdd57", "#b86bff"}
	layout := languagePie(WithColors(colours...)).Layout()
	for i, s := range layout.Slices {
		if s.Color != colours[i] {
			t.Errorf("slice %d colour = %s, want %s", i, s.Color, colours[i])
		}
	}
	if layout.Legend.Entries[2].Color != colours[2] {
		t.Errorf("legend colour = %s, want %s", layout.Legend.Entries[2].Color, colours[2])
	}
}

func TestPieLayout_SingleSliceIsFullCircle(t *testing.T) {
	chart := NewPieChart()
	chart.AddSlice("all", 42)
	layout := chart.Layout()
	if len(layout.Slices) != 1 {
		t.Fatalf("expected 1 slice, got %d", len(layout.Slices))
	}
	p := layout.Slices[0].Path
	if strings.Count(p, "A") != 2 || strings.Contains(p, "L") {
		t.Errorf("single slice should be a full circle of two arcs, got %q", p)
	}
}

func TestPieLayout_SkipsNonPositiveValues(t *testing.T) {
	chart := NewPieChart()
	chart.AddSlice("a", 3)
	chart.AddSlice("none", 0)
	chart.AddSlice("neg", -2)
	chart.AddSlice("b", 1)
	layout := chart.Layout()
	if len(layout.Slices) != 2 {
		t.Fatalf("expected 2 slices, got %d", len(layout.Slices))
	}
	if math.Abs(layout.Slices[0].Percent-75) > 0.01 {
		t.Errorf("a = %v%%, want 75", layout.Slices[0].Percent)
	}
}

func TestPieSparkline(t *testing.T) {
	chart := NewPieChart(WithVariant(Sparkline))
	chart.AddSlice("a", 1)
	chart.AddSlice("b", 3)
	layout := chart.Layout()
	if layout.Legend != nil || layout.Title != nil {
		t.Error("sparkline pie should have no legend or title")
	}
	if len(layout.Slices) != 2 {
		t.Fatalf("expected 2 slices, got %d", len(layout.Slices))
	}
	if math.Abs(layout.Pie.R-10) > 0.01 {
		t.Errorf("sparkline pie radius = %v, want 10 in a 100x20 box", layout.Pie.R)
	}
}

func TestPieRender(t *testing.T) {
	svg, err := languagePie().RenderString()
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(svg, `class="slice"`); n != 5 {
		t.Errorf("expected 5 slice paths, got %d", n)
	}
	if !strings.Contains(svg, ">Go 35.0%</text>") {
		t.Error("expected legend text with percentage")
	}
	if !strings.Contains(svg, "<title>Go: 35 (35.0%)</title>") {
		t.Error("expected per-slice hover title")
	}
	if !strings.Contains(svg, "Language Popularity") {
		t.Error("expected chart title")
	}
}
