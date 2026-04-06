<style>
/* Hide the broken default breadcrumb from md2html */
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
<li class="is-active"><a href="#" aria-current="page">01 — Sparkline</a></li>
</ul>
</nav>

<div class="example-nav">
<strong>01</strong> |
<a href="../02_static_line/">02</a> |
<a href="../03_multi_series/">03</a> |
<a href="../04_bar_chart/">04</a> |
<a href="../05_htmx/">05</a> |
<a href="../06_live_sse/">06</a> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a>
</div>

# 01 — Sparkline

A sparkline is a tiny inline chart — no axes, no labels, no legend. It shows a trend at a glance, designed to sit inside a paragraph or table cell. This is the simplest thing you can build with gogal: generate data, create a chart, render SVG.

<div class="columns">
<div class="column is-7">
<figure class="screenshot">
<p style="font-family: system-ui; font-size: 14px; margin: 0;">
  Current temperature: <img src="01_regular.svg" alt="Regular sparkline" style="vertical-align: middle; width: 200px; height: 30px;">
</p>
<p style="font-family: system-ui; font-size: 14px; margin: 0.5em 0 0;">
  Smooth variant: <img src="01_smooth.svg" alt="Smooth sparkline" style="vertical-align: middle; width: 200px; height: 30px;">
</p>
<p style="font-family: system-ui; font-size: 14px; margin: 0.5em 0 0;">
  Wider: <img src="01_wide.svg" alt="Wide sparkline" style="vertical-align: middle; width: 400px; height: 40px;">
</p>
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Three sparkline variants: regular, smooth (Catmull-Rom), and wider</figcaption>
</figure>
</div>
<div class="column">
<div class="buttons">
<a href="demo.html" class="button is-primary">Launch WASM Demo</a>
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/01_sparkline" class="button is-light">Source on Codeberg</a>
</div>
</div>
</div>

---

## The code

The entire example fits in one file. Generate data, configure the chart, render SVG:

```go
package main

import (
    "fmt"
    "math"
    "net/http"

    "codeberg.org/hum3/gogal"
)

func main() {
    // Generate sample data: a sine wave with noise
    var points []gogal.DataPoint
    for i := 0; i < 50; i++ {
        x := float64(i)
        y := math.Sin(x*0.3)*10 + 20 + math.Sin(x*1.1)*3
        points = append(points, gogal.DataPoint{X: x, Y: y})
    }

    // Create a sparkline
    chart := gogal.NewLineChart(
        gogal.WithVariant(gogal.Sparkline),
        gogal.WithSize(200, 30),
    )
    chart.Add("temperature", points)

    // Write SVG to stdout
    svg, _ := chart.RenderString()
    fmt.Println(svg)

    // Also serve on HTTP with three variants
    fmt.Println("Serving sparkline at http://localhost:1340")
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, `<p>Current temperature: `)
        chart.Render(w)
        fmt.Fprint(w, `</p>`)
        // ... smooth and wide variants follow
    })
    http.ListenAndServe(":1340", nil)
}
```

<div class="annotation">
<strong>gogal.NewLineChart()</strong> creates a line chart. <code>WithVariant(gogal.Sparkline)</code> strips it down to the bare minimum — no axes, labels, or legend. <code>WithSize(200, 30)</code> sets the SVG viewBox dimensions in pixels.
</div>

<div class="annotation">
<strong>chart.Add()</strong> adds a named data series. For sparklines, the name is only used internally — there's no legend to display it. Each <code>DataPoint</code> has X (position) and Y (value) fields.
</div>

<div class="annotation">
<strong>Render vs RenderString</strong> — <code>Render(w)</code> writes SVG directly to any <code>io.Writer</code> (HTTP response, file, etc.). <code>RenderString()</code> returns the SVG as a string. Both produce identical output. Use <code>Render</code> when streaming to avoid allocating the full SVG in memory.
</div>

<div class="annotation">
<strong>WithSmooth(true)</strong> enables Catmull-Rom to Bezier curve interpolation, producing smooth paths instead of straight line segments between points. Compare the regular and smooth variants above — same data, different visual feel.
</div>

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/01_sparkline/go/main.go)

---

## Running it

```bash
task example:01
```

This outputs SVG to stdout and starts an HTTP server at `http://localhost:1340` showing three sparkline variants (regular, smooth, wider).

---

## API reference

| Function | Purpose |
|----------|---------|
| `gogal.NewLineChart(opts...)` | Create a line chart with functional options |
| `gogal.WithVariant(gogal.Sparkline)` | Minimal inline chart, no axes/labels/legend |
| `gogal.WithSize(w, h)` | Set SVG viewBox dimensions |
| `gogal.WithSmooth(true)` | Enable Catmull-Rom curve smoothing |
| `chart.Add(name, points)` | Add a named data series |
| `chart.Render(w)` | Write SVG to an `io.Writer` |
| `chart.RenderString()` | Return SVG as a string |

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/01_sparkline) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
