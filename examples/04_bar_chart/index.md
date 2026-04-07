<style>
.annotation { border-left: 3px solid #3273dc; background: #f0f4ff; padding: 0.75em 1em; margin: 0.75em 0; border-radius: 0 4px 4px 0; font-size: 0.9em; }
.annotation strong { color: #3273dc; }
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
<a href="../09_dashboard/">09</a>
</div>

# 04 — Bar Chart

**Planned for Phase 2.** Bar charts will add vertical, horizontal, grouped, and stacked variants using a new `BandScale` for categorical axes.

---

## Planned API

```go
chart := gogal.NewBarChart(
    gogal.WithTitle("Monthly Sales"),
    gogal.WithGrid(true),
)
chart.Add("Revenue", []gogal.DataPoint{
    {Label: "Jan", Y: 120},
    {Label: "Feb", Y: 95},
    {Label: "Mar", Y: 140},
})
```

<div class="annotation">
<strong>BandScale</strong> will map categorical labels to equal-width bands, with configurable padding between bars. Grouped and stacked modes will subdivide bands for multi-series data.
</div>

## Variants planned

| Variant | Description |
|---------|-------------|
| Vertical | Standard column chart |
| Horizontal | Rotated bars |
| Grouped | Side-by-side bars per category |
| Stacked | Cumulative bars per category |

See [ROADMAP](../ROADMAP.html) for timeline.

---

[Back to examples](../index.html) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
