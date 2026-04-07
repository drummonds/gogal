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
<strong>09</strong>
</div>

# 09 — Dashboard

Multiple charts on one page with different sizes, themes, and configurations. This shows how gogal charts compose — each is an independent SVG with its own scale, theme, and options.

<div class="columns">
<div class="column is-6">
<figure class="screenshot">
<img src="09_light.svg" alt="Temperature chart — light theme">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Light theme (default)</figcaption>
</figure>
</div>
<div class="column is-6">
<figure class="screenshot" style="background: #2b2b2b;">
<img src="09_dark.svg" alt="Humidity chart — dark theme">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Dark theme</figcaption>
</figure>
</div>
</div>

<div class="buttons">
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/09_dashboard" class="button is-light">Source on Codeberg</a>
</div>

---

## The pattern

Each chart is created independently. Mix and match sizes, themes, and configurations — gogal produces self-contained SVGs that don't interfere with each other.

```go
// Temperature — light theme, 500x300
tc := gogal.NewLineChart(
    gogal.WithSize(500, 300),
    gogal.WithGrid(true),
    gogal.WithSmooth(true),
    gogal.WithYFormat("%.0f°C"),
)
tc.Add("Temperature", temp)
tc.Render(w)

// Humidity — dark theme, 500x300
hc := gogal.NewLineChart(
    gogal.WithSize(500, 300),
    gogal.WithTheme(gogal.ThemeDark),
    gogal.WithGrid(true),
    gogal.WithSmooth(true),
    gogal.WithYFormat("%.0f%%"),
)
hc.Add("Humidity", humidity)
hc.Render(w)
```

<div class="annotation">
<strong>WithTheme(gogal.ThemeDark)</strong> switches to a dark background with light text and grid lines. The 10-colour data palette is shared between themes. You can also create custom themes by constructing a <code>Theme</code> struct.
</div>

<div class="annotation">
<strong>WithSize(w, h)</strong> sets the SVG viewBox. Use CSS <code>width</code>/<code>max-width</code> on the container to control the rendered size — the SVG scales responsively.
</div>

### Mixed chart types on one page

The dashboard also includes:
- A **multi-series ordinal chart** spanning the full width
- **Sparklines** for quick at-a-glance metrics

```go
// Full-width multi-series with ordinal axis
all := gogal.NewLineChart(
    gogal.WithSize(1060, 350),
    gogal.WithAxisMode(gogal.Ordinal),
    gogal.WithLegend(true),
)
all.Add("Temperature (°C)", temp)
all.Add("Humidity (%)", humidity)
all.Add("Pressure (hPa)", pressure)

// Compact sparklines
sp := gogal.NewLineChart(
    gogal.WithVariant(gogal.Sparkline),
    gogal.WithSize(150, 25),
    gogal.WithSmooth(true),
)
sp.Add("Temp", temp)
```

<div class="annotation">
<strong>Layout is your responsibility</strong> — gogal produces SVGs, you arrange them with HTML/CSS. Use CSS Grid, Flexbox, or tables. Each chart is a self-contained block element.
</div>

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/09_dashboard/go/main.go)

---

## Running it

```bash
task example:09
```

Serves at `http://localhost:1348`.

---

## API reference

| Function | Purpose |
|----------|---------|
| `WithTheme(gogal.ThemeDark)` | Dark background theme |
| `WithTheme(gogal.ThemeLight)` | Light theme (default) |
| `WithSize(w, h)` | SVG viewBox dimensions |
| `WithLegend(true)` | Force legend display |
| `gogal.ThemeDark` / `gogal.ThemeLight` | Built-in theme presets |

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/09_dashboard) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
