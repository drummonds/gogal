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
<strong>03</strong> |
<a href="../04_bar_chart/">04</a> |
<a href="../05_htmx/">05</a> |
<a href="../06_live_sse/">06</a> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a>
</div>

# 03 — Multi-Series

Multiple data series on one chart with automatic colour cycling and legend. The two sensors sample at different rates (temperature ~hourly, humidity ~2-hourly) so the series have different point counts. This example compares two axis modes: **Temporal** (X proportional to time, gaps visible) and **Ordinal** (equal spacing by index).

<div class="columns">
<div class="column is-6">
<figure class="screenshot">
<img src="03_temporal.svg" alt="Multi-series chart with temporal axis">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Temporal axis — time-proportional spacing</figcaption>
</figure>
</div>
<div class="column is-6">
<figure class="screenshot">
<img src="03_ordinal.svg" alt="Multi-series chart with ordinal axis">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Ordinal axis — equal spacing by index</figcaption>
</figure>
</div>
</div>

<div class="buttons">
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/03_multi_series" class="button is-light">Source on Codeberg</a>
</div>

---

## The code

Adding multiple series is just multiple `Add()` calls. Colours cycle through the theme palette automatically.

```go
chart := gogal.NewLineChart(
    gogal.WithVariant(gogal.Static),
    gogal.WithAxisMode(gogal.Temporal),
    gogal.WithTitle("Weather — Temporal Axis"),
    gogal.WithGrid(true),
    gogal.WithTooltips(true),
    gogal.WithSmooth(true),
    gogal.WithTimeFormat("15:04"),
)
chart.Add("Temperature (°C)", temp)
chart.Add("Humidity (%)", humidity)
```

<div class="annotation">
<strong>Colour cycling</strong> — each <code>Add()</code> call gets the next colour from the theme palette (10 colours). The legend appears automatically when there are 2+ series.
</div>

<div class="annotation">
<strong>gogal.Temporal</strong> (default) — X positions are proportional to wall-clock time. If your data has a 6-hour gap, the chart shows that gap. Use this when time relationships matter.
</div>

<div class="annotation">
<strong>gogal.Ordinal</strong> — X positions are equally spaced by event index. The 1st data point is at the left edge, the last at the right, regardless of timestamps. Use this when you care about sequence, not absolute time.
</div>

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/03_multi_series/go/main.go)

---

## Running it

```bash
task example:03
```

Serves both charts at `http://localhost:1342`.

---

## API reference

| Function | Purpose |
|----------|---------|
| `chart.Add(name, points)` | Add a named series (call multiple times) |
| `WithAxisMode(gogal.Temporal)` | Time-proportional X spacing (default) |
| `WithAxisMode(gogal.Ordinal)` | Equal X spacing by index |
| `theme.SeriesColor(i)` | Get colour for series index `i` |

---

## Data

<div class="columns">
<div class="column is-6">

**Temperature** (21 readings)

| Time  | °C   |
|-------|------|
| 00:00 | 11.2 |
| 01:20 | 10.4 |
| 02:50 |  9.7 |
| 04:10 |  8.9 |
| 05:30 |  9.1 |
| 06:45 | 10.8 |
| 07:15 | 12.5 |
| 08:40 | 14.8 |
| 09:50 | 17.2 |
| 10:30 | 19.0 |
| 11:10 | 20.6 |
| 12:25 | 22.1 |
| 13:40 | 23.0 |
| 14:15 | 22.5 |
| 15:30 | 21.3 |
| 16:50 | 19.7 |
| 18:00 | 17.4 |
| 19:20 | 15.6 |
| 20:40 | 14.0 |
| 22:00 | 12.3 |
| 23:30 | 11.5 |

</div>
<div class="column is-6">

**Humidity** (13 readings)

| Time  |  %  |
|-------|-----|
| 00:00 |  82 |
| 02:30 |  85 |
| 04:45 |  88 |
| 06:30 |  84 |
| 08:15 |  72 |
| 10:00 |  58 |
| 11:30 |  50 |
| 13:00 |  45 |
| 15:00 |  48 |
| 17:00 |  55 |
| 19:00 |  65 |
| 21:30 |  76 |
| 23:00 |  80 |

</div>
</div>

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/03_multi_series) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
