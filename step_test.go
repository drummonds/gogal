package gogal

import (
	"fmt"
	"strings"
	"testing"
)

// pathVertices parses "Mx,y Lx,y ..." into its vertices.
func pathVertices(t *testing.T, path string) [][2]float64 {
	t.Helper()
	var out [][2]float64
	for _, seg := range strings.Fields(path) {
		var x, y float64
		if _, err := fmt.Sscanf(seg[1:], "%f,%f", &x, &y); err != nil {
			t.Fatalf("cannot parse segment %q of %q: %v", seg, path, err)
		}
		out = append(out, [2]float64{x, y})
	}
	return out
}

// A step chart holds each value until the next point, then jumps: every
// segment is either horizontal or vertical, alternating, and the path
// visits every data point.
func TestStepLayout_HoldsValueUntilNextPoint(t *testing.T) {
	chart := NewStepChart(WithSize(400, 200), WithSmooth(true)) // Smooth is ignored for steps
	chart.Add("balance", []DataPoint{{X: 0, Y: 100}, {X: 3, Y: 250}, {X: 4, Y: 180}, {X: 9, Y: 180}, {X: 10, Y: 300}})

	layout := chart.Layout()
	if len(layout.Series) != 1 {
		t.Fatalf("expected 1 series, got %d", len(layout.Series))
	}
	sl := layout.Series[0]
	if strings.Contains(sl.Path, "C") {
		t.Fatalf("step path must not be smoothed: %q", sl.Path)
	}
	v := pathVertices(t, sl.Path)
	pts := sl.Points

	// Vertices: p0, then for each next point a horizontal move to its x
	// at the old y and a vertical move to its y.
	wantLen := 1 + 2*(len(pts)-1)
	if len(v) != wantLen {
		t.Fatalf("expected %d vertices, got %d: %q", wantLen, len(v), sl.Path)
	}
	if v[0] != [2]float64{pts[0].X, pts[0].Y} {
		t.Errorf("path should start at the first point")
	}
	for i := 1; i < len(pts); i++ {
		h, vert := v[2*i-1], v[2*i]
		if h[0] != pts[i].X || h[1] != pts[i-1].Y {
			t.Errorf("point %d: horizontal step to (%v,%v), want (%v,%v)", i, h[0], h[1], pts[i].X, pts[i-1].Y)
		}
		if vert[0] != pts[i].X || vert[1] != pts[i].Y {
			t.Errorf("point %d: vertical step to (%v,%v), want (%v,%v)", i, vert[0], vert[1], pts[i].X, pts[i].Y)
		}
	}
	// Equal consecutive values give a flat run with no vertical jump in height.
	if v[6][1] != v[5][1] {
		t.Errorf("equal values should not change height: %v vs %v", v[5], v[6])
	}
}

func TestStepSparkline(t *testing.T) {
	chart := NewStepChart(WithVariant(Sparkline))
	chart.AddXY("s", []float64{0, 1, 2}, []float64{1, 3, 2})
	layout := chart.Layout()
	if layout.XAxis != nil || layout.Legend != nil {
		t.Error("sparkline step chart should have no axes or legend")
	}
	if n := len(pathVertices(t, layout.Series[0].Path)); n != 5 {
		t.Errorf("expected 5 vertices for 3 points, got %d", n)
	}
}

func TestStepRender(t *testing.T) {
	chart := NewStepChart(WithTitle("Balance"))
	chart.AddXY("balance", []float64{0, 1, 2}, []float64{10, 20, 15})
	svg, err := chart.RenderString()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `class="line"`) || !strings.Contains(svg, "<circle") {
		t.Error("step chart should draw a line path and data points")
	}
}
