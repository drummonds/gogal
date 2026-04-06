<style>
section.section > .container > nav.breadcrumb { display: none; }
.annotation { border-left: 3px solid #3273dc; background: #f0f4ff; padding: 0.75em 1em; margin: 0.75em 0; border-radius: 0 4px 4px 0; font-size: 0.9em; }
.annotation strong { color: #3273dc; }
.example-nav { background: #f5f5f5; border-radius: 4px; padding: 0.5em 1em; margin-bottom: 1.5em; font-size: 0.9em; }
.example-nav a { margin-right: 0.75em; }
.example-nav strong { color: #363636; }
</style>

<nav class="breadcrumb" aria-label="breadcrumbs">
<ul>
<li><a href="../index.html">gogal</a></li>
<li><a href="../index.html#examples">Examples</a></li>
<li class="is-active"><a href="#" aria-current="page">07 — Scatter</a></li>
</ul>
</nav>

<div class="example-nav">
<a href="../01_sparkline/">01</a> |
<a href="../02_static_line/">02</a> |
<a href="../03_multi_series/">03</a> |
<a href="../04_bar_chart/">04</a> |
<a href="../05_htmx/">05</a> |
<a href="../06_live_sse/">06</a> |
<strong>07</strong> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a>
</div>

# 07 — Scatter

**Planned for Phase 3.** Scatter charts will render XY data as points without connecting lines, with optional size and colour mapping per point.

---

## Planned API

```go
chart := gogal.NewScatterChart(
    gogal.WithTitle("Height vs Weight"),
    gogal.WithXTitle("Height (cm)"),
    gogal.WithYTitle("Weight (kg)"),
    gogal.WithGrid(true),
)
chart.Add("Male", malePoints)
chart.Add("Female", femalePoints)
```

<div class="annotation">
<strong>Per-point styling</strong> — scatter charts will use <code>DataPoint.Color</code> for individual point colours and a new <code>Size</code> field for bubble charts.
</div>

## Features planned

| Feature | Description |
|---------|-------------|
| Basic scatter | Points plotted at X,Y coordinates |
| Multi-series | Colour-coded series with legend |
| Bubble | Point size mapped to a third variable |
| Tooltips | Hover to see point labels |

See [ROADMAP](../ROADMAP.html) for timeline.

---

[Back to examples](../index.html) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
