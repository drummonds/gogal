<style>
.annotation { border-left: 3px solid #3273dc; background: #f0f4ff; padding: 0.75em 1em; margin: 0.75em 0; border-radius: 0 4px 4px 0; font-size: 0.9em; }
.annotation strong { color: #3273dc; }
.screenshot { border: 1px solid #dbdbdb; border-radius: 4px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); overflow: hidden; padding: 1em; background: #fff; }
.example-nav { background: #f5f5f5; border-radius: 4px; padding: 0.5em 1em; margin-bottom: 1.5em; font-size: 0.9em; }
.example-nav a { margin-right: 0.75em; }
.example-nav strong { color: #363636; }
</style>

<div class="example-nav">
<a href="../01_sparkline/">01</a> |
<a href="../01a_axis_formats/">01a</a> |
<a href="../02_static_line/">02</a> |
<a href="../03_multi_series/">03</a> |
<strong>04</strong> |
<a href="../05_htmx/">05</a> |
<a href="../06_live_sse/">06</a> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a> |
<a href="../10_pie_chart/">10</a>
</div>

# 04 — Bar Chart

A vertical bar chart: one bar per category, category labels on the X axis, the value printed above each bar, and a Y axis that always starts at zero. Adding a second series groups the bars side by side within each category.

<div class="columns">
<div class="column is-6">
<figure class="screenshot">
<img src="04_fibonacci.svg" alt="Bar chart of the first ten Fibonacci numbers">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Single series with value labels</figcaption>
</figure>
</div>
<div class="column is-6">
<figure class="screenshot">
<img src="04_grouped.svg" alt="Grouped bar chart of quarterly sales for two years">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Two series grouped per category</figcaption>
</figure>
</div>
</div>

<div class="buttons">
<a target="_blank" href="https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/04_bar_chart" class="button is-light">Source on Forgejo</a>
</div>

---

## The code

```go
chart := gogal.NewBarChart(
    gogal.WithTitle("Fibonacci (n=10)"),
    gogal.WithSize(500, 220),
    gogal.WithYFormat("%.0f"),
)
chart.AddCategories("Fibonacci",
    []string{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10"},
    []float64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55})
svg, _ := chart.RenderString()
```

<div class="annotation">
<strong>AddCategories(name, categories, values)</strong> takes parallel slices. The categories become the X-axis labels; call it once per series to get grouped bars.
</div>

<div class="annotation">
<strong>Baseline at zero.</strong> Unlike line charts, the Y scale always includes 0 so bar heights compare honestly. Negative values hang below the baseline. A little headroom is reserved above the tallest bar for its value label.
</div>

<div class="annotation">
<strong>WithValueLabels(false)</strong> turns off the value printed above each bar. Hovering a bar dims it slightly — pure CSS, no JavaScript.
</div>

<div class="annotation">
<strong>Sparkline bars.</strong> <code>NewBarChart(gogal.WithVariant(gogal.Sparkline))</code> renders a bare 100×20 bar strip with no axes, for inline use in tables.
</div>

[Full source on Forgejo](https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/04_bar_chart/go/main.go)

---

## Running it

```bash
task example:04
```

Outputs SVG to stdout and serves at `http://localhost:1343`.

---

## API reference

| Function | Purpose |
|----------|---------|
| `NewBarChart(opts...)` | Vertical bar chart |
| `AddCategories(name, categories, values)` | One series; repeat for grouped bars |
| `WithValueLabels(on)` | Print each value above its bar (default on) |
| `WithYFormat(fmt)` | Y-axis tick format (Printf) |

Horizontal and stacked bars are not yet implemented — see the [ROADMAP](../ROADMAP.html).

---

[Back to examples](../index.html) | [Source on Forgejo](https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/04_bar_chart) | [API docs](https://pkg.go.dev/git.bytestone.uk/hum3/gogal)
