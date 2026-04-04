package gogal

import (
	"fmt"
	"math"
	"strings"
)

// LayoutResult holds the computed layout of a chart, ready for rendering.
// This is the key separation point: layout math is testable without SVG parsing.
type LayoutResult struct {
	Width    float64
	Height   float64
	PlotArea Rect
	Title    *TextLayout
	XAxis    *AxisLayout
	YAxis    *AxisLayout
	Legend   *LegendLayout
	Series   []SeriesLayout
	Config   *ChartConfig
}

// Rect defines a rectangle in SVG coordinates.
type Rect struct {
	X, Y, Width, Height float64
}

// TextLayout positions a text element.
type TextLayout struct {
	X, Y     float64
	Text     string
	Anchor   string  // "start", "middle", "end"
	Rotation float64 // degrees
	FontSize float64
}

// AxisLayout positions an axis with its ticks and optional title.
type AxisLayout struct {
	Position  float64 // pixel position of the axis line
	Ticks     []TickLayout
	Title     *TextLayout
	GridLines []float64 // pixel positions for grid lines
}

// TickLayout positions a single tick mark.
type TickLayout struct {
	Position float64
	Label    string
}

// SeriesLayout holds the computed visual representation of a data series.
type SeriesLayout struct {
	Name     string
	Color    string
	Points   []PointLayout
	Path     string // SVG path d attribute
	CSSClass string // e.g. "series-0"
}

// PointLayout positions a data point.
type PointLayout struct {
	X, Y  float64
	Label string // tooltip text
	DataX string // formatted original X value
	DataY string // formatted original Y value
}

// LegendLayout positions the legend.
type LegendLayout struct {
	Entries []LegendEntry
	Rect    Rect
}

// LegendEntry represents one item in the legend.
type LegendEntry struct {
	Name     string
	Color    string
	CSSClass string
	Hidden   bool
}

// computeLayout calculates the full chart layout from config and series data.
func computeLayout(cfg *ChartConfig, series []Series, chartType string) *LayoutResult {
	result := &LayoutResult{
		Width:  cfg.Width,
		Height: cfg.Height,
		Config: cfg,
	}

	visibleSeries := filterVisibleSeries(series, cfg.HiddenSeries)

	if cfg.Variant == Sparkline {
		computeSparklineLayout(result, visibleSeries, cfg)
		return result
	}

	computeFullLayout(result, visibleSeries, cfg, chartType)
	return result
}

func filterVisibleSeries(series []Series, hidden []string) []Series {
	if len(hidden) == 0 {
		return series
	}
	hiddenSet := make(map[string]bool, len(hidden))
	for _, name := range hidden {
		hiddenSet[name] = true
	}
	var visible []Series
	for _, s := range series {
		if !hiddenSet[s.Name] {
			visible = append(visible, s)
		}
	}
	return visible
}

func computeSparklineLayout(result *LayoutResult, series []Series, cfg *ChartConfig) {
	result.PlotArea = Rect{X: 0, Y: 0, Width: cfg.Width, Height: cfg.Height}

	if len(series) == 0 {
		return
	}

	// Collect all Y values for scale
	allY := collectYValues(series)
	if len(allY) == 0 {
		return
	}

	yScale := NewLinearScaleFromData(allY)
	// Invert Y range: SVG y=0 is top
	yScale.SetRange(cfg.Height-1, 1)

	for i, s := range series {
		sl := computeSeriesLayout(s, i, nil, yScale, cfg)
		result.Series = append(result.Series, sl)
	}
}

func computeFullLayout(result *LayoutResult, series []Series, cfg *ChartConfig, chartType string) {
	m := cfg.Margins
	plotArea := Rect{
		X:      m.Left,
		Y:      m.Top,
		Width:  cfg.Width - m.Left - m.Right,
		Height: cfg.Height - m.Top - m.Bottom,
	}

	// Reserve space for legend at bottom
	if cfg.ShowLegend && len(series) > 1 {
		plotArea.Height -= 30
	}

	result.PlotArea = plotArea

	// Title
	if cfg.Title != "" {
		result.Title = &TextLayout{
			X:        cfg.Width / 2,
			Y:        m.Top - 10,
			Text:     cfg.Title,
			Anchor:   "middle",
			FontSize: cfg.Theme.FontSize + 4,
		}
		// Shift plot area down a bit for title
		result.PlotArea.Y += 10
		result.PlotArea.Height -= 10
	}

	allY := collectYValues(series)
	if len(allY) == 0 {
		return
	}

	// Build scales
	yScale := NewLinearScaleFromData(allY)
	yScale.SetRange(result.PlotArea.Y+result.PlotArea.Height, result.PlotArea.Y)

	var xScale Scale
	if cfg.Axis == Ordinal || chartType == "bar" {
		labels := collectXLabels(series)
		ordScale := NewOrdinalScale(labels)
		ordScale.SetRange(result.PlotArea.X, result.PlotArea.X+result.PlotArea.Width)
		xScale = ordScale
	} else {
		allX := collectXValues(series)
		linScale := NewLinearScaleFromData(allX)
		linScale.SetRange(result.PlotArea.X, result.PlotArea.X+result.PlotArea.Width)
		xScale = linScale
	}

	// Y axis
	yTicks := yScale.Ticks()
	yAxis := &AxisLayout{
		Position: result.PlotArea.X,
	}
	for _, tick := range yTicks {
		yAxis.Ticks = append(yAxis.Ticks, TickLayout{
			Position: tick.Position,
			Label:    tick.Label,
		})
		if cfg.ShowGrid {
			yAxis.GridLines = append(yAxis.GridLines, tick.Position)
		}
	}
	if cfg.YTitle != "" {
		yAxis.Title = &TextLayout{
			X:        15,
			Y:        result.PlotArea.Y + result.PlotArea.Height/2,
			Text:     cfg.YTitle,
			Anchor:   "middle",
			Rotation: -90,
			FontSize: cfg.Theme.FontSize,
		}
	}
	result.YAxis = yAxis

	// X axis
	xTicks := xScale.Ticks()
	xAxis := &AxisLayout{
		Position: result.PlotArea.Y + result.PlotArea.Height,
	}
	for _, tick := range xTicks {
		xAxis.Ticks = append(xAxis.Ticks, TickLayout{
			Position: tick.Position,
			Label:    tick.Label,
		})
	}
	if cfg.XTitle != "" {
		xAxis.Title = &TextLayout{
			X:        result.PlotArea.X + result.PlotArea.Width/2,
			Y:        cfg.Height - 5,
			Text:     cfg.XTitle,
			Anchor:   "middle",
			FontSize: cfg.Theme.FontSize,
		}
	}
	result.XAxis = xAxis

	// Series layouts
	for i, s := range series {
		sl := computeSeriesLayout(s, i, xScale, yScale, cfg)
		result.Series = append(result.Series, sl)
	}

	// Legend
	if cfg.ShowLegend && len(series) > 1 {
		result.Legend = computeLegendLayout(series, cfg)
	}
}

func computeSeriesLayout(s Series, index int, xScale Scale, yScale *LinearScale, cfg *ChartConfig) SeriesLayout {
	color := s.Color
	if color == "" {
		color = cfg.Theme.SeriesColor(index)
	}

	sl := SeriesLayout{
		Name:     s.Name,
		Color:    color,
		CSSClass: fmt.Sprintf("series-%d", index),
	}

	if len(s.Points) == 0 {
		return sl
	}

	// Build points
	for _, p := range s.Points {
		var px float64
		if xScale != nil {
			px = xScale.Map(p.X)
		} else {
			// Sparkline: linear mapping of index
			n := float64(len(s.Points) - 1)
			if n == 0 {
				px = cfg.Width / 2
			} else {
				idx := findPointIndex(s.Points, p)
				px = float64(idx) / n * cfg.Width
			}
		}
		py := yScale.Map(p.Y)

		label := p.Label
		if label == "" {
			label = fmt.Sprintf("%.4g", p.Y)
		}

		sl.Points = append(sl.Points, PointLayout{
			X:     px,
			Y:     py,
			Label: label,
			DataX: fmt.Sprintf("%.4g", p.X),
			DataY: fmt.Sprintf("%.4g", p.Y),
		})
	}

	// Build SVG path
	if cfg.Smooth {
		sl.Path = buildSmoothPath(sl.Points)
	} else {
		sl.Path = buildLinePath(sl.Points)
	}

	return sl
}

func findPointIndex(points []DataPoint, p DataPoint) int {
	for i, pt := range points {
		if pt.X == p.X && pt.Y == p.Y {
			return i
		}
	}
	return 0
}

func buildLinePath(points []PointLayout) string {
	if len(points) == 0 {
		return ""
	}
	var b strings.Builder
	for i, p := range points {
		if i == 0 {
			fmt.Fprintf(&b, "M%.2f,%.2f", p.X, p.Y)
		} else {
			fmt.Fprintf(&b, " L%.2f,%.2f", p.X, p.Y)
		}
	}
	return b.String()
}

func buildSmoothPath(points []PointLayout) string {
	if len(points) < 2 {
		return buildLinePath(points)
	}
	if len(points) == 2 {
		return buildLinePath(points)
	}

	// Catmull-Rom to Bezier conversion
	var b strings.Builder
	fmt.Fprintf(&b, "M%.2f,%.2f", points[0].X, points[0].Y)

	for i := 0; i < len(points)-1; i++ {
		var p0, p1, p2, p3 PointLayout
		if i == 0 {
			p0 = points[0]
		} else {
			p0 = points[i-1]
		}
		p1 = points[i]
		p2 = points[i+1]
		if i+2 < len(points) {
			p3 = points[i+2]
		} else {
			p3 = points[len(points)-1]
		}

		// Catmull-Rom to cubic Bezier control points
		cp1x := p1.X + (p2.X-p0.X)/6
		cp1y := p1.Y + (p2.Y-p0.Y)/6
		cp2x := p2.X - (p3.X-p1.X)/6
		cp2y := p2.Y - (p3.Y-p1.Y)/6

		fmt.Fprintf(&b, " C%.2f,%.2f %.2f,%.2f %.2f,%.2f",
			cp1x, cp1y, cp2x, cp2y, p2.X, p2.Y)
	}

	return b.String()
}

func computeLegendLayout(series []Series, cfg *ChartConfig) *LegendLayout {
	legend := &LegendLayout{
		Rect: Rect{
			X:      cfg.Margins.Left,
			Y:      cfg.Height - 25,
			Width:  cfg.Width - cfg.Margins.Left - cfg.Margins.Right,
			Height: 20,
		},
	}
	for i, s := range series {
		hidden := false
		for _, h := range cfg.HiddenSeries {
			if h == s.Name {
				hidden = true
				break
			}
		}
		color := s.Color
		if color == "" {
			color = cfg.Theme.SeriesColor(i)
		}
		legend.Entries = append(legend.Entries, LegendEntry{
			Name:     s.Name,
			Color:    color,
			CSSClass: fmt.Sprintf("series-%d", i),
			Hidden:   hidden,
		})
	}
	return legend
}

func collectYValues(series []Series) []float64 {
	var vals []float64
	for _, s := range series {
		for _, p := range s.Points {
			vals = append(vals, p.Y)
		}
	}
	return vals
}

func collectXValues(series []Series) []float64 {
	var vals []float64
	for _, s := range series {
		for _, p := range s.Points {
			vals = append(vals, p.X)
		}
	}
	return vals
}

func collectXLabels(series []Series) []string {
	seen := make(map[string]bool)
	var labels []string
	for _, s := range series {
		for _, p := range s.Points {
			label := p.Label
			if label == "" {
				label = fmt.Sprintf("%.4g", p.X)
			}
			if !seen[label] {
				seen[label] = true
				labels = append(labels, label)
			}
		}
	}
	return labels
}

// niceMin rounds down to a nice value for axis start.
func niceMin(value float64) float64 {
	if value == 0 {
		return 0
	}
	exp := math.Floor(math.Log10(math.Abs(value)))
	pow := math.Pow(10, exp)
	return math.Floor(value/pow) * pow
}
