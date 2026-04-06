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
<li class="is-active"><a href="#" aria-current="page">05 — Interactive HTMX</a></li>
</ul>
</nav>

<div class="example-nav">
<a href="../01_sparkline/">01</a> |
<a href="../02_static_line/">02</a> |
<a href="../03_multi_series/">03</a> |
<a href="../04_bar_chart/">04</a> |
<strong>05</strong> |
<a href="../06_live_sse/">06</a> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a>
</div>

# 05 — Interactive HTMX

Toggle series visibility without a page reload. Clicking a series name fires an HTMX request; the server re-renders the chart with `WithHiddenSeries` and returns just the SVG fragment. No custom JavaScript — just HTMX attributes and gogal's built-in hidden series support.

<div class="columns">
<div class="column is-6">
<figure class="screenshot">
<img src="05_all.svg" alt="All series visible">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">All series visible</figcaption>
</figure>
</div>
<div class="column is-6">
<figure class="screenshot">
<img src="05_toggled.svg" alt="Humidity series hidden">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Humidity toggled off</figcaption>
</figure>
</div>
</div>

<div class="buttons">
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/05_htmx" class="button is-light">Source on Codeberg</a>
</div>

---

## The pattern

The HTMX page loads the chart fragment via `hx-get="/chart"`. Toggle buttons below the chart send requests with `?hidden=Humidity,Wind Speed` to control which series are hidden.

```go
chart := gogal.NewLineChart(
    gogal.WithVariant(gogal.Static),
    gogal.WithTitle("Weather Station"),
    gogal.WithGrid(true),
    gogal.WithSmooth(true),
    gogal.WithHiddenSeries(hidden...),
)
chart.Add("Temperature", temp)
chart.Add("Humidity", humidity)
chart.Add("Wind Speed", wind)
```

<div class="annotation">
<strong>WithHiddenSeries(names...)</strong> tells the renderer to skip those series' paths but keep them in the legend (shown with strikethrough). The axis scaling only considers visible series, so the chart rescales automatically when you toggle.
</div>

### The HTML shell

```html
<div id="chart-container"
     hx-get="/chart"
     hx-trigger="load"
     hx-swap="innerHTML">
  Loading...
</div>
```

### The toggle endpoint

```go
http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
    hidden := strings.Split(r.URL.Query().Get("hidden"), ",")
    renderChart(hidden).Render(w)
    // ... render toggle buttons with hx-get attributes
})
```

<div class="annotation">
<strong>No framework needed</strong> — the toggle logic is ~30 lines of Go. Each button computes the new <code>hidden</code> parameter (add or remove its series name) and sets <code>hx-get="/chart?hidden=..."</code>. The server does the rest.
</div>

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/05_htmx/go/main.go)

---

## Running it

```bash
task example:05
```

Serves at `http://localhost:1344`. Click series names below the chart to toggle them.

---

## API reference

| Function | Purpose |
|----------|---------|
| `WithHiddenSeries(names...)` | Hide named series (paths hidden, legend shows strikethrough) |
| `Series.Hidden` | Per-series visibility flag (set via config) |

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/05_htmx) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
