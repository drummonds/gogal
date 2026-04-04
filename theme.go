package gogal

// Theme defines the visual appearance of a chart.
type Theme struct {
	Name       string
	Background string
	PlotBg     string
	Text       string
	Grid       string
	Axis       string
	Colors     []string // series color palette
	Font       string   // CSS font-family
	FontSize   float64
}

// ThemeLight is the default light theme.
var ThemeLight = &Theme{
	Name:       "light",
	Background: "#ffffff",
	PlotBg:     "#ffffff",
	Text:       "#333333",
	Grid:       "#e0e0e0",
	Axis:       "#666666",
	Colors: []string{
		"#4e79a7", "#f28e2b", "#e15759", "#76b7b2",
		"#59a14f", "#edc948", "#b07aa1", "#ff9da7",
		"#9c755f", "#bab0ac",
	},
	Font:     "system-ui, -apple-system, sans-serif",
	FontSize: 12,
}

// ThemeDark is a dark background theme.
var ThemeDark = &Theme{
	Name:       "dark",
	Background: "#1e1e1e",
	PlotBg:     "#2d2d2d",
	Text:       "#cccccc",
	Grid:       "#404040",
	Axis:       "#888888",
	Colors: []string{
		"#4e79a7", "#f28e2b", "#e15759", "#76b7b2",
		"#59a14f", "#edc948", "#b07aa1", "#ff9da7",
		"#9c755f", "#bab0ac",
	},
	Font:     "system-ui, -apple-system, sans-serif",
	FontSize: 12,
}

// SeriesColor returns the color for a series by index, cycling through the palette.
func (t *Theme) SeriesColor(index int) string {
	if len(t.Colors) == 0 {
		return "#4e79a7"
	}
	return t.Colors[index%len(t.Colors)]
}
