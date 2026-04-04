package gogal

// Variant controls the level of chart decoration and interactivity.
type Variant int

const (
	// Sparkline renders a minimal inline chart: just the data path,
	// no axes, labels, legend, or tooltips.
	Sparkline Variant = iota
	// Static renders a full chart with axes, labels, legend, and CSS tooltips.
	Static
	// Interactive adds HTMX-driven controls: legend toggling, axis switching.
	Interactive
	// Live adds SSE-driven re-rendering for real-time data.
	Live
)

// AxisMode controls how the X axis maps data to positions.
type AxisMode int

const (
	// Temporal maps X values proportional to wall-clock time.
	// Gaps between events are visible.
	Temporal AxisMode = iota
	// Ordinal maps X values to equal-width slots by event index.
	// Every data point gets equal spacing regardless of timing.
	Ordinal
)

// Margins defines chart padding in SVG coordinate units.
type Margins struct {
	Top, Right, Bottom, Left float64
}

// ChartConfig holds all configuration for a chart.
type ChartConfig struct {
	Width   float64
	Height  float64
	Title   string
	XTitle  string
	YTitle  string
	Variant Variant
	Axis    AxisMode
	Theme   *Theme
	Margins Margins

	ShowGrid     bool
	ShowLegend   bool
	ShowTooltips bool
	Accessible   bool
	Animate      bool
	Smooth       bool // use bezier curves for line charts

	HiddenSeries []string // series names to hide (for HTMX toggle)
	TimeFormat   string   // Go time format for tick labels
	YFormat      string   // printf format for Y labels
}

// Option is a functional option for configuring a chart.
type Option func(*ChartConfig)

// DefaultConfig returns sensible defaults for a static chart.
func DefaultConfig() ChartConfig {
	return ChartConfig{
		Width:        800,
		Height:       400,
		Variant:      Static,
		Axis:         Temporal,
		Theme:        ThemeLight,
		ShowGrid:     true,
		ShowLegend:   true,
		ShowTooltips: true,
		Accessible:   true,
		Margins:      Margins{Top: 40, Right: 20, Bottom: 50, Left: 60},
		TimeFormat:   "2006-01-02",
		YFormat:      "%.1f",
	}
}

// SparklineConfig returns defaults for a sparkline variant.
func SparklineConfig() ChartConfig {
	return ChartConfig{
		Width:   100,
		Height:  20,
		Variant: Sparkline,
		Theme:   ThemeLight,
	}
}

func WithTitle(title string) Option {
	return func(c *ChartConfig) { c.Title = title }
}

func WithSize(w, h float64) Option {
	return func(c *ChartConfig) { c.Width = w; c.Height = h }
}

func WithVariant(v Variant) Option {
	return func(c *ChartConfig) {
		c.Variant = v
		if v == Sparkline {
			c.ShowGrid = false
			c.ShowLegend = false
			c.ShowTooltips = false
			c.Accessible = false
			c.Width = 100
			c.Height = 20
		}
	}
}

func WithAxisMode(mode AxisMode) Option {
	return func(c *ChartConfig) { c.Axis = mode }
}

func WithTheme(theme *Theme) Option {
	return func(c *ChartConfig) { c.Theme = theme }
}

func WithGrid(show bool) Option {
	return func(c *ChartConfig) { c.ShowGrid = show }
}

func WithLegend(show bool) Option {
	return func(c *ChartConfig) { c.ShowLegend = show }
}

func WithTooltips(show bool) Option {
	return func(c *ChartConfig) { c.ShowTooltips = show }
}

func WithAccessibility(on bool) Option {
	return func(c *ChartConfig) { c.Accessible = on }
}

func WithAnimate(on bool) Option {
	return func(c *ChartConfig) { c.Animate = on }
}

func WithSmooth(on bool) Option {
	return func(c *ChartConfig) { c.Smooth = on }
}

func WithHiddenSeries(names ...string) Option {
	return func(c *ChartConfig) { c.HiddenSeries = names }
}

func WithTimeFormat(format string) Option {
	return func(c *ChartConfig) { c.TimeFormat = format }
}

func WithYFormat(format string) Option {
	return func(c *ChartConfig) { c.YFormat = format }
}

func WithMargins(top, right, bottom, left float64) Option {
	return func(c *ChartConfig) {
		c.Margins = Margins{Top: top, Right: right, Bottom: bottom, Left: left}
	}
}

func WithXTitle(title string) Option {
	return func(c *ChartConfig) { c.XTitle = title }
}

func WithYTitle(title string) Option {
	return func(c *ChartConfig) { c.YTitle = title }
}
