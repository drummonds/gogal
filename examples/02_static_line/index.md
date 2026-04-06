<style>
section.section > .container > nav.breadcrumb { display: none; }
.annotation { border-left: 3px solid #3273dc; background: #f0f4ff; padding: 0.75em 1em; margin: 0.75em 0; border-radius: 0 4px 4px 0; font-size: 0.9em; }
.annotation strong { color: #3273dc; }
.screenshot { border: 1px solid #dbdbdb; border-radius: 4px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); overflow: hidden; padding: 1em; background: #fff; }
.example-nav { background: #f5f5f5; border-radius: 4px; padding: 0.5em 1em; margin-bottom: 1.5em; font-size: 0.9em; }
.example-nav a { margin-right: 0.75em; }
.example-nav strong { color: #363636; }
</style>

<nav class="breadcrumb" aria-label="breadcrumbs">
<ul>
<li><a href="../index.html">gogal</a></li>
<li><a href="../index.html#examples">Examples</a></li>
<li class="is-active"><a href="#" aria-current="page">02 — Static Line</a></li>
</ul>
</nav>

<div class="example-nav">
<a href="../01_sparkline/">01</a> |
<strong>02</strong> |
<a href="../03_multi_series/">03</a> |
<a href="../04_bar_chart/">04</a> |
<a href="../05_htmx/">05</a> |
<a href="../06_live_sse/">06</a> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a>
</div>

# 02 — Static Line

A full line chart with all the trimmings: axes with titles, grid lines, legend, CSS tooltips on hover, and entry animation. The data points are irregularly spaced in time (real weather-station readings, not one-per-hour) to show how gogal handles uneven temporal spacing on the X axis.

<div class="columns">
<div class="column is-8">
<figure class="screenshot">
<img src="02_static_line.svg" alt="Static line chart with axes, grid, and tooltips">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Temperature over 24 hours — full static chart</figcaption>
</figure>
</div>
<div class="column">
<div class="buttons">
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/02_static_line" class="button is-light">Source on Codeberg</a>
</div>
</div>
</div>

---

## The code

```go
// Irregular time spacing — not every hour
measurements := []struct{ hour, min int; temp float64 }{
    {0, 0, 11.2}, {1, 15, 10.5}, {2, 45, 9.8}, /* ... */
}
for _, m := range measurements {
    t := time.Date(2024, 6, 15, m.hour, m.min, 0, 0, time.UTC)
    points = append(points, gogal.DataPoint{Time: t, Y: m.temp, Label: fmt.Sprintf("%.1f°C", m.temp)})
}

chart := gogal.NewLineChart(
    gogal.WithVariant(gogal.Static),
    gogal.WithTitle("Temperature — 15 June 2024"),
    gogal.WithXTitle("Hour"),
    gogal.WithYTitle("°C"),
    gogal.WithGrid(true),
    gogal.WithTooltips(true),
    gogal.WithAnimate(true),
    gogal.WithSmooth(true),
    gogal.WithTimeFormat("15:04"),
    gogal.WithYFormat("%.0f"),
)
chart.Add("Temperature", points)
```

<div class="annotation">
<strong>WithVariant(gogal.Static)</strong> is the default — it enables axes, legend, grid, and tooltips. You don't need to specify it explicitly, but it makes the intent clear when reading the code.
</div>

<div class="annotation">
<strong>WithTitle / WithXTitle / WithYTitle</strong> add text labels to the chart and axes. Titles are positioned automatically based on the chart margins.
</div>

<div class="annotation">
<strong>WithGrid(true)</strong> draws horizontal grid lines aligned to Y-axis ticks. Grid lines use the theme's grid colour and sit behind the data.
</div>

<div class="annotation">
<strong>WithTooltips(true)</strong> enables CSS-only tooltips — hover over any data point to see its label. No JavaScript required. The tooltip content comes from the <code>Label</code> field of each <code>DataPoint</code>.
</div>

<div class="annotation">
<strong>WithAnimate(true)</strong> adds a CSS stroke-dashoffset animation — the line draws itself on page load. Pure CSS, no JavaScript.
</div>

<div class="annotation">
<strong>WithTimeFormat / WithYFormat</strong> control tick label formatting. <code>TimeFormat</code> uses Go's time layout; <code>YFormat</code> uses Printf syntax.
</div>

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/02_static_line/go/main.go)

---

## Running it

```bash
task example:02
```

Outputs SVG to stdout and serves at `http://localhost:1341`.

---

## API reference

| Function | Purpose |
|----------|---------|
| `WithVariant(gogal.Static)` | Full chart with axes, legend, tooltips |
| `WithTitle(s)` | Chart title |
| `WithXTitle(s)` / `WithYTitle(s)` | Axis titles |
| `WithGrid(true)` | Horizontal grid lines |
| `WithTooltips(true)` | CSS hover tooltips |
| `WithAnimate(true)` | Stroke-dash entry animation |
| `WithTimeFormat(layout)` | X-axis tick format (Go time) |
| `WithYFormat(fmt)` | Y-axis tick format (Printf) |

---

## Data

| Time  | Temp (°C) |
|-------|-----------|
| 00:00 | 11.2      |
| 01:15 | 10.5      |
| 02:45 |  9.8      |
| 03:10 |  9.1      |
| 04:50 |  8.7      |
| 06:05 |  9.4      |
| 07:00 | 11.0      |
| 07:40 | 12.3      |
| 08:30 | 14.1      |
| 09:55 | 16.8      |
| 10:20 | 18.5      |
| 11:00 | 20.1      |
| 12:35 | 22.4      |
| 13:10 | 23.0      |
| 14:00 | 22.7      |
| 14:45 | 21.9      |
| 15:50 | 20.6      |
| 16:30 | 19.2      |
| 17:15 | 17.8      |
| 18:40 | 15.9      |
| 19:00 | 14.3      |
| 20:25 | 13.1      |
| 21:50 | 12.0      |
| 23:10 | 11.4      |

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/02_static_line) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
