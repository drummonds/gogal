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
pie.go              # Pie chart layout (sectors, vertical legend)
scale_test.go       # Scale tests (13, incl. temporal ladder)
layout_test.go      # Layout + render tests (9)
bar_test.go         # Bar chart tests (5)
pie_test.go         # Pie chart tests (6)
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
task docs:build     # build documentation HTML (markdown -> docs/)
task docs:all       # screenshots + WASM demos + HTML
task clean          # remove build artifacts (wasm, example binaries, generated HTML)
tp pages deploy     # rsync docs/ to the docs site (run docs:all first)
```

## Implementation Status

### Done
- Core: Scale, Layout, Render pipeline
- Line chart (sparkline + static variants)
- Bar chart (vertical, grouped, value labels, sparkline)
- Pie chart (sectors, legend with percentages, WithColors)
- Step chart (balance-style hold-then-jump line)
- Temporal ticks in calendar units, formatted in the data's location
- CSS tooltips, hover effects, animations
- Smooth paths (Catmull-Rom to Bezier)
- Light/Dark themes
- 36 tests passing

### Next Steps (from plan)
- Step 4: Static line chart example (examples/02_static_line, examples/03_multi_series)
- Step 5: Bar chart — horizontal and stacked variants
- Step 6: Scatter chart
- Step 7: HTMX integration (htmx.go, HandleChart, legend toggling)
- Step 8: Live/SSE (RenderSSEHandler)
- Step 9: WASM example
- Step 10: Documentation site + deploy

## Docs Site & Deploy

- Source: https://git.bytestone.uk/hum3/gogal (origin, Forgejo); mirror https://github.com/drummonds/gogal
- Docs: https://gogal.docs.bytestone.uk/ — Caddy static host on woodpecker-ci, served from `/srv/sites/gogal`
- Deploy: run `task docs:all` first, then `tp pages deploy`, which only rsyncs `docs/` over SSH as the `deploy` user (it does NOT run the `pages_build` steps in task-plus.yml; those run during `tp release`)
- Site registration lives in the Caddyfile in `~/Cloudstation/IT/statichost` (already applied for gogal)

Build artifacts are never committed — rebuild them with `task docs:all` before deploying:
- `docs/*/main.wasm` and `docs/*/wasm_exec.js` (from `docs:build-wasm`)
- Native example executables in `examples/*/go/` (left behind by `go build`)
- `docs/*.html` rendered from the top-level markdown files

Committed under `docs/`: `index.md`, `examples.md`, per-example `index.html`, `demo.html`, `app.js` and screenshot SVGs.

## Reference
- Full research: RESEARCH.md
- Implementation plan: .claude/plans/curious-tumbling-newt.md
- Follows lofigui patterns for project structure, docs, examples
