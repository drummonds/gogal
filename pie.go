package gogal

import (
	"fmt"
	"math"
)

// PieLayout positions the pie: centre and radius in SVG coordinates.
type PieLayout struct {
	CX, CY, R float64
}

// SliceLayout holds the computed geometry of one pie slice.
type SliceLayout struct {
	Label    string
	Value    string  // formatted value
	Percent  float64 // share of the total, 0–100
	Color    string
	Path     string // SVG path d attribute for the sector
	CSSClass string // e.g. "slice-0"
}

// pieLegendRowHeight is the vertical pitch of pie legend entries.
const pieLegendRowHeight = 18

// computePieLayout lays out a pie chart: title, sectors, and a vertical
// legend of "label percent%" entries to the right of the pie.
func computePieLayout(result *LayoutResult, series []Series, cfg *ChartConfig) {
	if cfg.Variant == Sparkline {
		result.PlotArea = Rect{X: 0, Y: 0, Width: cfg.Width, Height: cfg.Height}
	} else {
		m := cfg.Margins
		result.PlotArea = Rect{
			X:      m.Left,
			Y:      m.Top,
			Width:  cfg.Width - m.Left - m.Right,
			Height: cfg.Height - m.Top - m.Bottom,
		}
		if cfg.Title != "" {
			result.Title = &TextLayout{
				X:        cfg.Width / 2,
				Y:        m.Top - 10,
				Text:     cfg.Title,
				Anchor:   "middle",
				FontSize: cfg.Theme.FontSize + 4,
			}
			result.PlotArea.Y += 10
			result.PlotArea.Height -= 10
		}
	}

	var points []DataPoint
	total := 0.0
	for _, s := range series {
		for _, p := range s.Points {
			if p.Y > 0 {
				points = append(points, p)
				total += p.Y
			}
		}
	}
	if total == 0 {
		return
	}

	// Legend entries with their percentage, then reserve a column for them.
	showLegend := cfg.ShowLegend && cfg.Variant != Sparkline
	legendWidth := 0.0
	if showLegend {
		result.Legend = &LegendLayout{}
		for i, p := range points {
			name := fmt.Sprintf("%s %.1f%%", p.Label, 100*p.Y/total)
			result.Legend.Entries = append(result.Legend.Entries, LegendEntry{
				Name:     name,
				Color:    sliceColor(p, i, cfg),
				CSSClass: fmt.Sprintf("slice-%d", i),
			})
			legendWidth = math.Max(legendWidth, float64(len(name))*7+24)
		}
	}

	pieArea := result.PlotArea
	pieArea.Width -= legendWidth
	r := math.Min(pieArea.Width, pieArea.Height) / 2
	if cfg.Variant != Sparkline {
		r -= 2 // keep the slice outline inside the plot area
	}
	if r < 1 {
		r = 1
	}
	result.Pie = &PieLayout{
		CX: pieArea.X + pieArea.Width/2,
		CY: pieArea.Y + pieArea.Height/2,
		R:  r,
	}

	if showLegend {
		x := pieArea.X + pieArea.Width + 8
		y := result.Pie.CY - float64(len(points))*pieLegendRowHeight/2
		result.Legend.Rect = Rect{X: x, Y: y, Width: legendWidth, Height: float64(len(points)) * pieLegendRowHeight}
		for i := range result.Legend.Entries {
			result.Legend.Entries[i].X = x
			result.Legend.Entries[i].Y = y + float64(i)*pieLegendRowHeight
		}
	}

	angle := -math.Pi / 2 // start at 12 o'clock, sweep clockwise
	for i, p := range points {
		fraction := p.Y / total
		result.Slices = append(result.Slices, SliceLayout{
			Label:    p.Label,
			Value:    fmt.Sprintf("%.4g", p.Y),
			Percent:  100 * fraction,
			Color:    sliceColor(p, i, cfg),
			Path:     sectorPath(result.Pie, angle, fraction),
			CSSClass: fmt.Sprintf("slice-%d", i),
		})
		angle += 2 * math.Pi * fraction
	}
}

func sliceColor(p DataPoint, index int, cfg *ChartConfig) string {
	if p.Color != "" {
		return p.Color
	}
	return cfg.seriesColor(index)
}

// sectorPath returns the SVG path for a sector starting at angle
// (radians, clockwise from 3 o'clock) covering fraction of the circle.
func sectorPath(pie *PieLayout, start, fraction float64) string {
	if fraction >= 1-1e-9 {
		// A single arc from a point to itself draws nothing; use two half arcs.
		return fmt.Sprintf("M%.2f,%.2f A%.2f,%.2f 0 1,1 %.2f,%.2f A%.2f,%.2f 0 1,1 %.2f,%.2f Z",
			pie.CX-pie.R, pie.CY, pie.R, pie.R, pie.CX+pie.R, pie.CY, pie.R, pie.R, pie.CX-pie.R, pie.CY)
	}
	end := start + 2*math.Pi*fraction
	x0, y0 := pie.CX+pie.R*math.Cos(start), pie.CY+pie.R*math.Sin(start)
	x1, y1 := pie.CX+pie.R*math.Cos(end), pie.CY+pie.R*math.Sin(end)
	large := 0
	if fraction > 0.5 {
		large = 1
	}
	return fmt.Sprintf("M%.2f,%.2f L%.2f,%.2f A%.2f,%.2f 0 %d,1 %.2f,%.2f Z",
		pie.CX, pie.CY, x0, y0, pie.R, pie.R, large, x1, y1)
}
