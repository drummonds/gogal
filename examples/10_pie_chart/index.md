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
<a href="../04_bar_chart/">04</a> |
<a href="../05_htmx/">05</a> |
<a href="../06_live_sse/">06</a> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a> |
<strong>10</strong>
</div>

# 10 — Pie Chart

Proportional slices starting at 12 o'clock and sweeping clockwise, with a legend to the right listing each label and its percentage. Hovering a slice shows its value in a native tooltip.

<div class="columns">
<div class="column is-6">
<figure class="screenshot">
<img src="10_theme.svg" alt="Pie chart of language popularity in theme colours">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Theme palette</figcaption>
</figure>
</div>
<div class="column is-6">
<figure class="screenshot">
<img src="10_colours.svg" alt="Pie chart of language popularity with an explicit palette">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Explicit colours via WithColors</figcaption>
</figure>
</div>
</div>

<div class="buttons">
<a target="_blank" href="https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/10_pie_chart" class="button is-light">Source on Forgejo</a>
</div>

---

## The code

```go
chart := gogal.NewPieChart(
    gogal.WithTitle("Language Popularity"),
    gogal.WithSize(450, 300),
    gogal.WithColors("#3298dc", "#48c774", "#f14668", "#ffdd57", "#b86bff"),
)
chart.AddSlice("Go", 35)
chart.AddSlice("Python", 28)
chart.AddSlice("Rust", 18)
chart.AddSlice("JS", 12)
chart.AddSlice("Other", 7)
svg, _ := chart.RenderString()
```

<div class="annotation">
<strong>AddSlice(label, value)</strong> adds one slice. Percentages are computed from the total of all positive values; slices with a zero or negative value are skipped.
</div>

<div class="annotation">
<strong>WithColors(...)</strong> replaces the theme palette. Slices take colours in order and cycle if there are more slices than colours. Leave it out to use the theme's ten-colour palette.
</div>

<div class="annotation">
<strong>Legend and tooltips.</strong> The legend reads "Go 35.0%". Each slice carries an SVG <code>&lt;title&gt;</code>, so browsers show "Go: 35 (35.0%)" on hover with no JavaScript.
</div>

[Full source on Forgejo](https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/10_pie_chart/go/main.go)

---

## Running it

```bash
task example:10
```

Outputs SVG to stdout and serves at `http://localhost:1349`.

---

## API reference

| Function | Purpose |
|----------|---------|
| `NewPieChart(opts...)` | Pie chart |
| `AddSlice(label, value)` | One slice; values are shares of the total |
| `WithColors(colours...)` | Palette override, in slice order |
| `WithLegend(false)` | Hide the legend |

Donut charts are not yet implemented — see the [ROADMAP](../ROADMAP.html).

---

[Back to examples](../index.html) | [Source on Forgejo](https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/10_pie_chart) | [API docs](https://pkg.go.dev/git.bytestone.uk/hum3/gogal)
