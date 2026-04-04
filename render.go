package gogal

import (
	"fmt"
	"io"
)

// render writes the SVG for a computed layout.
func render(w io.Writer, layout *LayoutResult) error {
	sw := newSVGWriter(w)
	cfg := layout.Config

	sw.openSVG(layout.Width, layout.Height, cfg.Accessible, cfg.Title)

	// Accessibility
	if cfg.Accessible && cfg.Title != "" {
		sw.title(cfg.Title)
		sw.desc(fmt.Sprintf("Chart: %s", cfg.Title))
	}

	// CSS
	if cfg.Variant != Sparkline {
		css := generateCSS(cfg)
		if css != "" {
			sw.style(css)
		}
	}

	// Background
	if cfg.Variant != Sparkline {
		sw.rect(0, 0, layout.Width, layout.Height,
			fmt.Sprintf(` fill="%s"`, cfg.Theme.Background))
	}

	// Grid lines
	if cfg.ShowGrid && layout.YAxis != nil {
		sw.openGroup(` class="grid"`)
		for _, y := range layout.YAxis.GridLines {
			sw.line(layout.PlotArea.X, y,
				layout.PlotArea.X+layout.PlotArea.Width, y,
				fmt.Sprintf(` stroke="%s" stroke-width="0.5"`, cfg.Theme.Grid))
		}
		sw.closeGroup()
	}

	// Y axis
	if layout.YAxis != nil {
		sw.openGroup(` class="y-axis"`)
		// Axis line
		sw.line(layout.YAxis.Position, layout.PlotArea.Y,
			layout.YAxis.Position, layout.PlotArea.Y+layout.PlotArea.Height,
			fmt.Sprintf(` stroke="%s" stroke-width="1"`, cfg.Theme.Axis))
		// Ticks
		for _, tick := range layout.YAxis.Ticks {
			sw.line(layout.YAxis.Position-4, tick.Position,
				layout.YAxis.Position, tick.Position,
				fmt.Sprintf(` stroke="%s" stroke-width="1"`, cfg.Theme.Axis))
			sw.text(layout.YAxis.Position-8, tick.Position+4, tick.Label,
				fmt.Sprintf(` text-anchor="end" fill="%s" font-size="%.0f" font-family="%s"`,
					cfg.Theme.Text, cfg.Theme.FontSize, cfg.Theme.Font))
		}
		// Title
		if layout.YAxis.Title != nil {
			t := layout.YAxis.Title
			sw.text(t.X, t.Y, t.Text,
				fmt.Sprintf(` text-anchor="%s" fill="%s" font-size="%.0f" font-family="%s" transform="rotate(%.0f %.2f %.2f)"`,
					t.Anchor, cfg.Theme.Text, t.FontSize, cfg.Theme.Font, t.Rotation, t.X, t.Y))
		}
		sw.closeGroup()
	}

	// X axis
	if layout.XAxis != nil {
		sw.openGroup(` class="x-axis"`)
		// Axis line
		sw.line(layout.PlotArea.X, layout.XAxis.Position,
			layout.PlotArea.X+layout.PlotArea.Width, layout.XAxis.Position,
			fmt.Sprintf(` stroke="%s" stroke-width="1"`, cfg.Theme.Axis))
		// Ticks
		for _, tick := range layout.XAxis.Ticks {
			sw.line(tick.Position, layout.XAxis.Position,
				tick.Position, layout.XAxis.Position+4,
				fmt.Sprintf(` stroke="%s" stroke-width="1"`, cfg.Theme.Axis))
			sw.text(tick.Position, layout.XAxis.Position+16, tick.Label,
				fmt.Sprintf(` text-anchor="middle" fill="%s" font-size="%.0f" font-family="%s"`,
					cfg.Theme.Text, cfg.Theme.FontSize, cfg.Theme.Font))
		}
		// Title
		if layout.XAxis.Title != nil {
			t := layout.XAxis.Title
			sw.text(t.X, t.Y, t.Text,
				fmt.Sprintf(` text-anchor="%s" fill="%s" font-size="%.0f" font-family="%s"`,
					t.Anchor, cfg.Theme.Text, t.FontSize, cfg.Theme.Font))
		}
		sw.closeGroup()
	}

	// Title
	if layout.Title != nil {
		t := layout.Title
		sw.text(t.X, t.Y, t.Text,
			fmt.Sprintf(` text-anchor="%s" fill="%s" font-size="%.0f" font-weight="bold" font-family="%s"`,
				t.Anchor, cfg.Theme.Text, t.FontSize, cfg.Theme.Font))
	}

	// Series
	for _, sl := range layout.Series {
		sw.openGroup(fmt.Sprintf(` class="%s"`, sl.CSSClass))

		// Path
		if sl.Path != "" {
			sw.path(sl.Path,
				fmt.Sprintf(` fill="none" stroke="%s" stroke-width="2" class="line"`, sl.Color))
		}

		// Data points
		if cfg.Variant != Sparkline {
			for _, p := range sl.Points {
				sw.openGroup(fmt.Sprintf(` class="data-point" aria-label="%s"`, xmlEscape(p.Label)))
				sw.circle(p.X, p.Y, 3, fmt.Sprintf(` fill="%s" class="point"`, sl.Color))
				if cfg.ShowTooltips {
					// Tooltip group: hidden by default, shown on hover via CSS
					sw.openGroup(` class="tooltip"`)
					// Background rect
					sw.rect(p.X+8, p.Y-20, 60, 18, ` fill="white" stroke="#ccc" rx="3"`)
					sw.text(p.X+12, p.Y-7, p.Label,
						fmt.Sprintf(` fill="%s" font-size="10" font-family="%s"`,
							cfg.Theme.Text, cfg.Theme.Font))
					sw.closeGroup()
				}
				sw.closeGroup()
			}
		}

		sw.closeGroup()
	}

	// Legend
	if layout.Legend != nil {
		sw.openGroup(` class="legend"`)
		xOff := layout.Legend.Rect.X
		for _, entry := range layout.Legend.Entries {
			sw.rect(xOff, layout.Legend.Rect.Y, 12, 12,
				fmt.Sprintf(` fill="%s"`, entry.Color))
			textAttr := fmt.Sprintf(` fill="%s" font-size="%.0f" font-family="%s"`,
				cfg.Theme.Text, cfg.Theme.FontSize, cfg.Theme.Font)
			if entry.Hidden {
				textAttr += ` text-decoration="line-through" opacity="0.5"`
			}
			sw.text(xOff+16, layout.Legend.Rect.Y+11, entry.Name, textAttr)
			xOff += float64(len(entry.Name))*8 + 30
		}
		sw.closeGroup()
	}

	sw.closeSVG()
	return sw.err
}

// generateCSS returns the embedded CSS for chart interactivity.
func generateCSS(cfg *ChartConfig) string {
	css := ""
	css += "    .tooltip { display: none; pointer-events: none; }\n"
	css += "    .data-point:hover .tooltip { display: block; }\n"
	css += "    .point { transition: r 0.15s; }\n"
	css += "    .data-point:hover .point { r: 5; }\n"
	css += "    .line { transition: stroke-width 0.2s; }\n"

	if cfg.Animate {
		css += `    .line {
      stroke-dasharray: 2000;
      stroke-dashoffset: 2000;
      animation: draw 1.5s ease forwards;
    }
    @keyframes draw {
      to { stroke-dashoffset: 0; }
    }
`
	}

	return css
}
