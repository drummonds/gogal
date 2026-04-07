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
<strong>01a</strong> |
<a href="../02_static_line/">02</a> |
<a href="../03_multi_series/">03</a> |
<a href="../04_bar_chart/">04</a> |
<a href="../05_htmx/">05</a> |
<a href="../06_live_sse/">06</a> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a>
</div>

# 01a — Axis Formats

Visual test of axis tick formatting across different data ranges. This gallery shows how gogal formats axis ticks for numeric, scientific, temporal, and ordinal data.

<div class="columns">
<div class="column is-8">
<figure class="screenshot">
<img src="01a_complete.svg" alt="Axis format gallery">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-1">All axis format examples</figcaption>
</figure>
</div>
<div class="column">
<div class="buttons">
<a href="demo.html" class="button is-primary">Launch WASM Demo</a>
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/01a_axis_formats" class="button is-light">Source on Codeberg</a>
</div>
</div>
</div>

---

## What it demonstrates

| Category | Input Range | Expected tick labels |
|----------|-------------|---------------------|
| Small integers | 0-10 | 0, 2, 4, 6, 8, 10 |
| Decimals | 0-1 | 0, 0.2, 0.4, ... |
| Hundreds | 0-500 | 0, 100, 200, ... |
| Large | 0-100k | 0, 20000, 40000, ... |
| Scientific small | 1e-6 to 1e-3 | 1e-06, ... |
| Scientific large | 1e6 to 1e9 | 1e+06, ... |
| Negative | -50 to 50 | -40, -20, 0, 20, 40 |
| Narrow | 10.0-10.5 | 10, 10.1, 10.2, ... |
| Temporal (hours) | 24h | 00:00, 04:00, ... |
| Temporal (days) | 30d | Jun 01, Jun 10, ... |
| Temporal (months) | 12 months | Jan, Feb, Mar, ... |
| Temporal (years) | 2020-2030 | 2020, 2022, 2024, ... |
| Ordinal (weekdays) | Mon-Sun | Mon, Tue, ... |
| Ordinal (categories) | Named items | Alpha, Beta, ... |

<div class="annotation">
<strong>Temporal axis</strong> uses <code>TemporalScale</code> which formats Unix timestamp ticks as human-readable times via Go's <code>time.Format()</code>. Set the format with <code>WithTimeFormat("15:04")</code>.
</div>

<div class="annotation">
<strong>Ordinal axis</strong> uses <code>OrdinalScale</code> which maps discrete indices to equal-width bands. Labels come from <code>DataPoint.Label</code> when no <code>Time</code> is set.
</div>

---

## Running

```bash
task example:01a
```

Serves at [http://localhost:1339](http://localhost:1339).

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/01a_axis_formats)
