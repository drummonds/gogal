package gogal

import (
	"strings"
	"testing"
	"time"
)

// dailySeries returns one point per day from start for the given number of
// days, counting up from 20.
func dailySeries(start time.Time, days int) ([]time.Time, []float64) {
	var times []time.Time
	var values []float64
	for i := 0; i <= days; i++ {
		times = append(times, start.AddDate(0, 0, i))
		values = append(values, float64(20+i))
	}
	return times, values
}

func tickLabels(axis *AxisLayout) []string {
	var labels []string
	for _, tick := range axis.Ticks {
		labels = append(labels, tick.Label)
	}
	return labels
}

func TestRender_WithoutPoints(t *testing.T) {
	times, values := dailySeries(time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), 30)
	svg, err := NewLineChart(WithPoints(false), WithTooltips(false)).
		AddTimeSeries("customers", times, values).RenderString()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(svg, "<circle") || strings.Contains(svg, "data-point") {
		t.Error("WithPoints(false) should draw no point markers")
	}
	if !strings.Contains(svg, `class="line"`) {
		t.Error("the line itself should still be drawn")
	}
}

// Counts need whole-number ticks: with a %.0f format, half steps would be
// rounded into repeated labels placed at the wrong value.
func TestYAxis_IntegerTicks(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   string
	}{
		{"flat", []float64{20, 20}, "19 20 21"},
		{"narrow", []float64{20, 23}, "20 21 22 23"},
		{"wide", []float64{20, 1121}, "0 200 400 600 800 1000 1200"},
	}
	for _, tt := range tests {
		chart := NewLineChart(WithIntegerY(true), WithYFormat("%.0f"))
		chart.AddXY("n", []float64{0, 1}, tt.values)
		if got := strings.Join(tickLabels(chart.Layout().YAxis), " "); got != tt.want {
			t.Errorf("%s: y labels = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// The Y domain is widened to whole ticks, so the axis is labelled at both
// the bottom and the top of the plot.
func TestYAxis_NiceDomainLabelsBothEnds(t *testing.T) {
	chart := NewLineChart(WithYFormat("%.0f"))
	chart.AddXY("n", []float64{0, 1, 2}, []float64{20, 500, 1121})
	layout := chart.Layout()

	ticks := layout.YAxis.Ticks
	bottom := layout.PlotArea.Y + layout.PlotArea.Height
	if first := ticks[0]; first.Label != "0" || first.Position != bottom {
		t.Errorf("first y tick = %q at %v, want \"0\" at the plot bottom %v", first.Label, first.Position, bottom)
	}
	if last := ticks[len(ticks)-1]; last.Label != "1200" || last.Position != layout.PlotArea.Y {
		t.Errorf("last y tick = %q at %v, want \"1200\" at the plot top %v", last.Label, last.Position, layout.PlotArea.Y)
	}
}

// A date-only format cannot show sub-day ticks, so the step never drops
// below a day; likewise a month format never steps in weeks.
func TestTemporalScale_StepNoFinerThanFormat(t *testing.T) {
	tests := []struct {
		format     string
		start, end time.Time
		want       string
	}{
		{"2006-01-02", time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC),
			"2026-01-15 2026-01-16"},
		{"Jan 2006", time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
			"Feb 2026 Mar 2026"},
	}
	for _, tt := range tests {
		ls := NewLinearScale(float64(tt.start.Unix()), float64(tt.end.Unix()))
		ls.SetRange(0, 560)
		var labels []string
		for _, tick := range NewTemporalScale(ls, tt.format, time.UTC).Ticks() {
			labels = append(labels, tick.Label)
		}
		if got := strings.Join(labels, " "); got != tt.want {
			t.Errorf("%s: labels = %q, want %q", tt.format, got, tt.want)
		}
	}
}

// WithEndLabels pins the first and last dates to the axis ends and drops
// calendar ticks whose labels would collide with them.
func TestXAxis_EndLabels(t *testing.T) {
	start := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	for _, days := range []int{1, 9, 30, 120, 400, 1100} {
		times, values := dailySeries(start, days)
		chart := NewLineChart(WithSize(660, 180), WithMargins(10, 40, 24, 60), WithEndLabels(true))
		chart.AddTimeSeries("customers", times, values)
		layout := chart.Layout()
		ticks := layout.XAxis.Ticks

		first, last := ticks[0], ticks[len(ticks)-1]
		if first.Label != times[0].Format("2006-01-02") || first.Position != layout.PlotArea.X {
			t.Errorf("%dd: first x tick = %q at %v, want the start date at the axis start", days, first.Label, first.Position)
		}
		end := layout.PlotArea.X + layout.PlotArea.Width
		if last.Label != times[len(times)-1].Format("2006-01-02") || last.Position != end {
			t.Errorf("%dd: last x tick = %q at %v, want the end date at the axis end", days, last.Label, last.Position)
		}
		for i := 1; i < len(ticks); i++ {
			if gap := ticks[i].Position - ticks[i-1].Position; gap < 70 {
				t.Errorf("%dd: labels %q and %q are %.0fpx apart, they overlap", days, ticks[i-1].Label, ticks[i].Label, gap)
			}
		}
	}
}

// Plain XY data carries no times, so its X axis is numeric even though
// the default axis mode is temporal.
func TestXAxis_NumericWhenPointsCarryNoTime(t *testing.T) {
	chart := NewLineChart(WithYFormat("%.0f"))
	chart.AddXY("n", []float64{0, 2, 4, 6, 8, 10}, []float64{1, 3, 2, 5, 4, 6})
	if got := strings.Join(tickLabels(chart.Layout().XAxis), " "); got != "0 2 4 6 8 10" {
		t.Errorf("x labels = %q, want \"0 2 4 6 8 10\"", got)
	}
}

// WithMaxYTicks caps the number of labelled Y ticks, picking the finest
// nice step that fits once the domain is widened to whole steps. The cap
// is at least 3: a domain crossing zero always needs -step, 0 and +step.
func TestYAxis_MaxTicks(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		max    int
		want   string
	}{
		{"step 10", []float64{20, 52}, 5, "20 30 40 50 60"},
		{"step 500", []float64{20, 1121}, 5, "0 500 1000 1500"},
		{"flat", []float64{20, 20}, 5, "19 20 21"},
		{"three", []float64{20, 142}, 3, "0 100 200"},
		{"two floors at three", []float64{20, 142}, 2, "0 100 200"},
	}
	for _, tt := range tests {
		chart := NewLineChart(WithIntegerY(true), WithYFormat("%.0f"), WithMaxYTicks(tt.max))
		chart.AddXY("n", []float64{0, 1}, tt.values)
		if got := strings.Join(tickLabels(chart.Layout().YAxis), " "); got != tt.want {
			t.Errorf("%s: y labels = %q, want %q", tt.name, got, tt.want)
		}
	}

	for hi := 1.0; hi <= 5000; hi *= 1.37 {
		for _, limit := range []int{2, 3, 4, 5, 6} {
			chart := NewLineChart(WithMaxYTicks(limit))
			chart.AddXY("n", []float64{0, 1}, []float64{-hi / 3, hi})
			if n := len(chart.Layout().YAxis.Ticks); n > max(limit, 3) || n < 2 {
				t.Errorf("data -%.3g..%.3g, max %d: %d y ticks", hi/3, hi, limit, n)
			}
		}
	}
}
