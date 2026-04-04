# gogal — SVG Chart Library Research

## Design Constraints

- **No JavaScript** — HTMX is the JS ceiling (it's a framework, not custom JS)
- **Pure SVG** — server-side rendered, clean markup
- **CSS-only interactivity** — hover effects, tooltips, transitions via embedded `<style>`
- **HTMX for discrete interactions** — legend toggling, re-rendering, filter/period controls
- **lofigui alignment** — HTML/CSS-first, single-page re-rendering, HTMX multipage
- **Pure Go** — zero or minimal deps, no cgo, gokrazy compatible
- **Dual axis modes** — time-proportional and event-proportional X axes

---

## 1. Existing Go SVG Chart Libraries

### Viable candidates (pure SVG, no JS)

| Library | Deps | Chart Types | API | Last Release | Adoption |
|---------|------|-------------|-----|-------------|----------|
| **margaid** (erkkah) | Zero | Line, Smooth, Bar | Functional options | Mar 2023 | 4 importers |
| **go-chart** (wcharczuk) | freetype | Line, Bar, Pie, Donut, Financial (SMA/EMA/MACD/Bollinger) | Struct config | Aug 2024 | 567 importers |
| **go-charts** (vicanso) | go-chart | Line, Bar, HBar, Pie, Radar, Funnel, Table | Functional options | Aug 2024 | 115 importers |
| **gosvgchart** (riclib) | Zero | Line, Bar, Pie/Donut, Heatmap | Builder/chaining | Feb 2025 | New |
| **gonum/plot** | gonum ecosystem | 15+ scientific types | OOP/interfaces | Mar 2025 | Mature |
| **go-svg-charts** (fabienmasson) | 9 pkgs | Line, Bar, Area, Pie, Treemap, HeatMap, GeoMap | Builder/chaining | Sep 2023 | 0 importers |

### Ruled out

- **go-echarts** — wraps Apache ECharts JS. Requires browser. Not pure SVG.
- **tomarus/chart** — embedded JS. Abandoned (2019).
- **gopie** — pie/donut only. Unmaintained.
- **svgPlot** — minimal, basic.

### Assessment of enhance-vs-build candidates

**margaid** — Closest philosophically: zero deps, pure SVG, `io.Writer` output, functional options, time-series aware with `TimeTicker`, streaming data support (series capping by age/size). Weaknesses: only 3 chart types, no CSS interactivity, pre-v1, last release 2023. Would need substantial expansion for chart types, CSS tooltips, HTMX patterns, and dual-axis support.

**gosvgchart** — Most modern: zero deps, responsive SVG (`viewBox`, `width="100%"`), CSS dark mode via `prefers-color-scheme`, markdown chart format, Goldmark extension. Weaknesses: pre-v1, limited chart types, unfamiliar maintainer, no interactivity beyond pie tooltips.

**go-chart** — Most battle-tested: 567 importers, financial chart series, stable v2. Weaknesses: freetype dependency (cgo risk), struct-config API is less ergonomic than functional options, no CSS interactivity, no HTMX awareness.

**go-charts** — Best API of the wrappers: functional options, ECharts JSON compatibility, fast (~3.3ms SVG). Inherits go-chart's freetype dependency.

**gonum/plot** — Most capable for scientific plotting, publication quality, multiple output backends. Too verbose for web charting, heavy dependency tree.

---

## 2. Pygal (Python reference)

### What to adopt
- **Builder pattern API**: `chart = pygal.Bar()` → `.add('Series', data)` → `.render()` — simple, low ceremony
- **Legend toggling**: clicking legend items toggles series visibility — most beloved feature
- **Theming via style objects**: named themes + custom style structs with colors, opacity, transitions
- **Secondary Y axis**: `chart.add('B', data, secondary=True)`
- **Rich data points**: `{'value': 10, 'label': 'Ten', 'color': '#f00'}` — per-point overrides
- **`human_readable` number formatting**: 1K, 1M, etc.
- **Multiple render targets**: string, file, `io.Writer`, data URI

### What to avoid
- **Embedded JavaScript**: pygal's `pygal-tooltips.js` (~15KB) breaks `<img>` embedding, bloats output, and violates the no-JS constraint
- **Sprawling kwargs**: dozens of constructor options discoverable only via source code
- **Poor accessibility**: no ARIA, no `<desc>` on data elements
- **SVG bloat**: 50KB+ for trivial charts due to embedded JS/CSS

### Key lesson
Pygal's interactivity (tooltips, legend toggling) is implemented via embedded JS. **The same interactions can be achieved with CSS-only (tooltips/hover) and HTMX (legend toggling with server re-render and axis rescaling)** — without any embedded JS.

---

## 3. HTMX + SVG Interactivity

### Critical rule
**Always swap the entire `<svg>` element into an HTML `<div>` container.** Never swap inner SVG content — browser's `innerHTML` parser creates elements in the wrong namespace and they won't render.

```html
<div id="chart-container" hx-get="/chart" hx-trigger="...">
  <svg viewBox="0 0 800 400"><!-- server-rendered --></svg>
</div>
```

### Interaction feasibility matrix

| Interaction | Approach | Feasibility |
|---|---|---|
| Filter/period controls | Pure HTMX | Excellent |
| Legend toggle (show/hide series) | Pure HTMX | Excellent — server rescales axes |
| Axis mode switch (time ↔ event) | Pure HTMX | Excellent |
| Live data updates | HTMX SSE extension | Good (1-2 Hz) |
| Zoom buttons (in/out/reset) | Pure HTMX | Excellent |
| Click-drag range select | HTMX + ~20 lines JS | Good |
| Hover tooltips (single chart) | CSS-only | Excellent |
| Hover highlight (single chart) | CSS-only | Excellent |
| Cross-chart crosshair | ~30 lines JS | Not for HTMX |
| Scroll-to-zoom | Client JS only | Too latency-sensitive |
| Continuous pan | Client JS only | Too latency-sensitive |

### CSS-only interactivity (no JS, no HTMX)

```css
/* Tooltips */
.data-point:hover .tooltip { display: block; }

/* Hover highlight */
.bar { transition: opacity 0.2s; }
.bar:hover { opacity: 0.8; }

/* Line drawing animation */
path.line {
  stroke-dasharray: 1000;
  stroke-dashoffset: 1000;
  animation: draw 2s forwards;
}

/* Legend toggling (CSS-only, no axis rescale) */
#toggle-b:not(:checked) ~ svg .series-b { display: none; }
```

### HTMX legend toggling (with axis rescale)

```html
<div id="chart-area">
  <svg><!-- chart --></svg>
  <div class="legend">
    <button hx-get="/chart?hide=B" hx-target="#chart-area">Series A</button>
    <button hx-get="/chart?hide=A" hx-target="#chart-area">Series B</button>
  </div>
</div>
```

Server tracks visible series, re-renders with recalculated Y-axis scale. CSS-only toggling can't rescale axes — HTMX approach is superior.

### Live updates via SSE

```html
<div hx-ext="sse" sse-connect="/events/chart" sse-swap="update">
  <div id="live-chart">
    <svg><!-- initial --></svg>
  </div>
</div>
```

Works well at 1-2 updates/second. Each event sends complete `<svg>`.

---

## 4. Time vs Event Ordering

### The problem
Same dataset, two valid X-axis interpretations:
- **Time-proportional**: X maps to wall-clock time. Gaps visible. 3-hour gap = 3× space.
- **Event-proportional**: X maps to event index. Equal spacing regardless of timing.

### Architectural pattern
The X-axis is a **pluggable scale**. Both implement the same interface:

```go
type Scale interface {
    Map(value float64) float64    // data → pixel
    Inverse(pixel float64) float64 // pixel → data
    Ticks() []Tick
}
```

- `TimeScale` — maps timestamps proportionally
- `EventScale` — maps event indices to equal-width slots

The rendering pipeline stays identical: `data → scale(x) → pixel position`. Switching modes = swapping the scale.

### Gold-standard API (Highcharts pattern)
A single boolean: `ordinal: true/false`. In Go:

```go
chart := gogal.NewLineChart(data,
    gogal.WithAxisMode(gogal.Ordinal), // or gogal.Temporal
)
```

### Linked views
Link on **data identity**, not position. Hover on point in View A → look up same point by ID → highlight in View B. The axis mapping is irrelevant to linking.

### Gap handling
- **Remove gaps entirely** (ordinal mode)
- **Compress gaps** to fixed width (hybrid)
- **Show break indicators** (zigzag on axis)

---

## 5. Broader Design Patterns

### From Vega-Lite / Observable Plot
- **Data type system**: `Quantitative`, `Ordinal`, `Nominal`, `Temporal` — drives automatic scale selection, axis formatting, legend generation. Single most impactful pattern.
- **Channel-based encoding**: `x`, `y`, `color`, `size` map to data fields
- **Automatic scales and axes**: inferred from data types

### From D3
- **Scales as first-class objects**: composable, reusable, testable
- **Layout computation separated from rendering**: enables testing without parsing SVG
- **Shape generators**: functions that compute SVG path `d` attributes

### From matplotlib
- **Backend/renderer pattern**: scene graph → visitor → SVG output
- **`<defs>`/`<use>` deduplication**: define repeated elements (markers, patterns) once, reference via `<use>`. Critical for SVG file size.

### SVG best practices
- **Accessibility**: `role="img"`, `<title>`, `<desc>`, `aria-labelledby`, `aria-label` on data elements
- **Responsive**: `viewBox` + no fixed width/height + CSS `max-width: 100%`
- **Styling**: presentation attributes for defaults, CSS classes for overridability, embedded `<style>` for hover/transitions
- **Animation**: CSS `@keyframes` for line drawing, bar entrance — no JS needed

---

## 6. Build vs Enhance Decision

### Option A: Enhance margaid
**Pros**: Zero deps, pure SVG, functional options API, time-series aware, streaming data support, `io.Writer` output — philosophically aligned.
**Cons**: Only 3 chart types, unmaintained since 2023, no CSS interactivity, no HTMX patterns, ISC license (fine but unusual). Would need: 5+ new chart types, CSS tooltip layer, HTMX integration helpers, dual-axis support, accessibility, theming. Essentially a rewrite wearing margaid's clothes.

### Option B: Enhance gosvgchart
**Pros**: Zero deps, responsive SVG, CSS dark mode, modern design, actively developed.
**Cons**: Pre-v1, limited chart types, someone else's project direction, no dual-axis support. Same expansion scope as margaid.

### Option C: Build new (gogal)
**Pros**: Full control over architecture, purpose-built for lofigui/HTMX integration, can adopt best patterns from all research (Vega-Lite encoding, D3 scales, pygal ergonomics, HTMX-first interactivity). Zero deps from day one. Dual-axis as first-class feature. CSS-only interactivity by design.
**Cons**: More initial work. No existing community/adoption.

### Option D: Fork + extend go-chart
**Pros**: Battle-tested, 567 importers, financial series.
**Cons**: freetype dependency (cgo risk for gokrazy), struct-config API less ergonomic, would need significant refactoring to add CSS interactivity and HTMX patterns.

### Recommendation: **Build new (Option C)**

The gap between existing libraries and the requirements is large enough that enhancing any of them would be a near-rewrite anyway. No existing Go library has:
- CSS-only interactivity
- HTMX integration patterns
- Time/event dual-axis support
- Accessibility (ARIA)
- Zero deps + responsive SVG

Building new allows adopting the best architectural patterns from the research:
1. **Scale interface** (from D3) — pluggable TimeScale/EventScale
2. **Functional options API** (from margaid, Go idiom) — `WithAxisMode()`, `WithTheme()`
3. **Data type-driven automation** (from Vega-Lite) — auto scale/axis/legend from data types
4. **CSS-in-SVG interactivity** (from SVG best practices) — tooltips, hover, transitions
5. **HTMX integration helpers** — `RenderToHandler()`, query param parsing for visible series/axis mode/zoom range
6. **Accessibility by default** (from WAI guidelines) — `<title>`, `<desc>`, ARIA labels
7. **`io.Writer` streaming** (Go idiom) — no intermediate string allocation
8. **`<defs>`/`<use>` deduplication** (from matplotlib) — compact SVG output

Draw inspiration from margaid (zero deps, streaming), gosvgchart (responsive SVG), and pygal (API ergonomics) without being constrained by their architectures.

---

## 7. Suggested Initial Scope

### Phase 1 — Core
- Line chart (straight + smooth/bezier)
- Bar chart (vertical, horizontal, grouped, stacked)
- Scatter/XY plot
- Scale interface with TimeScale and EventScale
- Functional options API
- CSS-only tooltips and hover effects
- Responsive SVG (`viewBox`, no fixed dimensions)
- Accessibility (`<title>`, `<desc>`, ARIA)
- Basic theming (light/dark)

### Phase 2 — HTMX Integration
- `http.Handler` helpers for chart endpoints
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
- Cross-chart crosshair (minimal JS module, ~30 lines)
- Drag-to-select zoom (minimal JS module, ~20 lines)
- Gap compression / break axes
- Financial series (SMA, EMA, Bollinger)
