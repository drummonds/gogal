# gogal - Pure Go SVG Chart Library

## What is gogal

Pure Go SVG chart library with CSS-only interactivity and HTMX integration. Zero dependencies, no cgo, gokrazy compatible.

## Project Structure

```
gogal.go            # Public API: NewLineChart, Add, Render, Layout
option.go           # ChartConfig, functional options, Variant/AxisMode types
scale.go            # Scale interface, LinearScale, OrdinalScale
layout.go           # Layout computation (sparkline + full chart)
render.go           # SVG renderer (layout → io.Writer)
svg.go              # Low-level SVG writing helpers
series.go           # DataPoint, Series types
theme.go            # Theme struct, ThemeLight, ThemeDark
scale_test.go       # Scale tests (9)
layout_test.go      # Layout + render tests (8)
examples/
  01_sparkline/go/  # Minimal sparkline example
docs/
  index.md          # Documentation site index
RESEARCH.md         # Design research and decisions
```

## Architecture

```
Data → Scale(x,y) → Layout() → LayoutResult{} → Render(io.Writer) → SVG
```

Key separation: `Layout()` returns a plain struct testable without SVG parsing.

### Scale Interface
- `LinearScale` — continuous quantitative mapping with nice tick generation
- `OrdinalScale` — discrete indices to equal-width bands

### Variants (per chart type)
1. **Sparkline** — minimal inline, no axes/labels/legend
2. **Static** — full chart with axes, labels, legend, CSS tooltips
3. **Interactive** — HTMX-driven legend toggling, axis switching (not yet built)
4. **Live** — SSE-driven re-rendering (not yet built)

### Axis Modes
- `Temporal` — X proportional to wall-clock time
- `Ordinal` — X proportional to event index (equal spacing)

## Design Constraints
- **No JavaScript** — HTMX is the JS ceiling
- **Pure Go** — zero deps, gokrazy compatible
- **CSS-in-SVG** — tooltips, hover, transitions via embedded `<style>`
- **Accessible** — ARIA, `<title>`, `<desc>` by default

## Running

```
task check          # fmt + vet + test
task example:01     # run sparkline example
task docs:build     # build documentation HTML
```

## Implementation Status

### Done
- Core: Scale, Layout, Render pipeline
- Line chart (sparkline + static variants)
- CSS tooltips, hover effects, animations
- Smooth paths (Catmull-Rom to Bezier)
- Light/Dark themes
- 17 tests passing

### Next Steps (from plan)
- Step 4: Static line chart example (examples/02_static_line, examples/03_multi_series)
- Step 5: Bar chart (BandScale, vertical/horizontal/grouped/stacked)
- Step 6: Scatter chart
- Step 7: HTMX integration (htmx.go, HandleChart, legend toggling)
- Step 8: Live/SSE (RenderSSEHandler)
- Step 9: WASM example
- Step 10: Documentation site + deploy

### Repo Setup Needed
- Create Codeberg repo: codeberg.org/hum3/gogal
- Create GitHub mirror: github.com/drummonds/gogal
- Create statichost site: h3-gogal at https://builder.statichost.eu
```
Site name: h3-gogal
Repository: https://codeberg.org/hum3/gogal
Branch: main
Publish directory: docs
```

## Reference
- Full research: RESEARCH.md
- Implementation plan: .claude/plans/curious-tumbling-newt.md
- Follows lofigui patterns for project structure, docs, examples
