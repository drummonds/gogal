# Roadmap

## Done

- Research: surveyed Go SVG libraries, pygal, HTMX+SVG patterns, time/event ordering
- Decision: build new library (no existing Go lib meets all requirements)
- Project scaffolding

## In Progress

- Core architecture (Scale, Layout, Render)
- Line chart (sparkline variant)

## Planned

### Phase 1 — Core
- Line chart (all 4 variants: sparkline, static, interactive, live)
- Bar chart (vertical, horizontal, grouped, stacked)
- Scatter/XY chart
- Scale interface with TimeScale and EventScale
- CSS-only tooltips and hover effects
- Responsive SVG (viewBox, no fixed dimensions)
- Accessibility (ARIA, title, desc)
- Basic theming (light/dark)

### Phase 2 — HTMX Integration
- http.Handler helpers for chart endpoints
- Query param parsing (visible series, axis mode, zoom range)
- HTMX attribute generation for legend controls
- SSE helper for live chart updates

### Phase 3 — Extended Chart Types
- Area chart (filled line)
- Pie/donut
- Heatmap
- Box plot

### Phase 4 — Advanced
- Linked/synchronized views
- Cross-chart crosshair (minimal JS module)
- Drag-to-select zoom (minimal JS module)
- Gap compression / break axes
- Financial series (SMA, EMA, Bollinger)

### Phase 5 — Documentation
- Documentation site with examples increasing in complexity
- WASM demos for interactive examples
- Architecture diagrams (d2)
