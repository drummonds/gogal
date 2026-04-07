<style>
.annotation { border-left: 3px solid #3273dc; background: #f0f4ff; padding: 0.75em 1em; margin: 0.75em 0; border-radius: 0 4px 4px 0; font-size: 0.9em; }
.annotation strong { color: #3273dc; }
.screenshot { border: 1px solid #dbdbdb; border-radius: 4px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); overflow: hidden; padding: 1em; background: #fff; }
.example-nav { background: #f5f5f5; border-radius: 4px; padding: 0.5em 1em; margin-bottom: 1.5em; font-size: 0.9em; }
.example-nav a { margin-right: 0.75em; }
.example-nav strong { color: #363636; }
</style>

<div class="example-nav">
<strong>01</strong> |
<a href="../01a_axis_formats/">01a</a> |
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

<div class="columns is-vcentered">
<div class="column is-5">
<figure class="image screenshot">
<img src="01_polling.svg" alt="During generation — partial sparkline">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-1">During generation</figcaption>
</figure>
</div>
<div class="column is-narrow has-text-centered">&#x27A1;</div>
<div class="column is-5">
<figure class="image screenshot">
<img src="01_complete.svg" alt="After completion — full sparkline">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-1">Complete</figcaption>
</figure>
</div>
</div>

<div class="buttons">
<a href="demo.html" class="button is-primary">Launch WASM Demo</a>
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/01_sparkline" class="button is-light">Source on Codeberg</a>
</div>

---

## The code

The model function generates points one at a time, updating the sparkline after each. Uses lofigui for interactive WASM and server modes:

```go
func model(app *lofigui.App) {
    const nPoints = 7
    var points []gogal.DataPoint
    y := 20.0

    for i := range nPoints {
        y += rand.Float64()*10 - 5
        points = append(points, gogal.DataPoint{X: float64(i), Y: y})

        lofigui.Reset()
        lofigui.HTML(`<p>Sparkline: ` + renderSparkline(points, false) + `</p>`)
        lofigui.HTML(`<p>Smooth: ` + renderSparkline(points, true) + `</p>`)

        if i < nPoints-1 {
            app.Sleep(500 * time.Millisecond)
        }
    }
}
```

<div class="annotation">
<strong>Sequential generation</strong> — each iteration adds one point, calls <code>lofigui.Reset()</code> to clear the buffer, then re-renders the sparkline with the growing dataset. The browser polls for updates, showing the chart grow in real time.
</div>

<div class="annotation">
<strong>app.Sleep()</strong> pauses between points. The lofigui framework handles rendering in the browser during the pause, so each intermediate state is visible.
</div>

<div class="annotation">
<strong>WithSmooth(true)</strong> enables Catmull-Rom to Bezier curve interpolation, producing smooth paths instead of straight line segments between points.
</div>

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/01_sparkline/go/model.go)

---

## Running it

```bash
task example:01
```

Starts the lofigui server at `http://localhost:1340`. Click Start to generate a random sparkline — watch it grow point by point.

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
