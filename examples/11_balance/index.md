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
<a href="../10_pie_chart/">10</a> |
<strong>11</strong>
</div>

# 11 — Balance

A running account balance drawn as a step chart. A balance does not drift between transactions: it holds its value until a posting lands, then changes instantly. A straight line between points would invent a gradual change that never happened, so the step chart draws a horizontal hold followed by a vertical jump at each data point.

<div class="columns">
<div class="column is-8">
<figure class="screenshot">
<img src="11_balance.svg" alt="Step chart of a current account balance over a month">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Balance holds between transactions and jumps at each one</figcaption>
</figure>
</div>
<div class="column">
<div class="buttons">
<a target="_blank" href="https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/11_balance" class="button is-light">Source on Forgejo</a>
</div>
</div>
</div>

---

## The code

```go
chart := gogal.NewStepChart(
    gogal.WithTitle("Current account — September 2026"),
    gogal.WithYTitle("£"),
    gogal.WithYFormat("%.0f"),
    gogal.WithTimeFormat("2 Jan"),
)
for _, tx := range ledger {
    points = append(points, gogal.DataPoint{
        Time:  tx.date,
        Y:     tx.balance, // balance after this posting
        Label: fmt.Sprintf("%s: £%.0f", tx.desc, tx.balance),
    })
}
chart.Add("Balance", points)
```

<div class="annotation">
<strong>NewStepChart</strong> takes the same options and data as <code>NewLineChart</code>; only the path differs. Each point is the value <em>after</em> the event, so the line holds the previous value up to the point's X and then steps to the new one. <code>WithSmooth</code> is ignored.
</div>

<div class="annotation">
<strong>Temporal axis.</strong> Transactions are unevenly spaced, so the default temporal axis places each step at its real date and the long flat runs show exactly how long the balance sat unchanged. Use <code>WithAxisMode(gogal.Ordinal)</code> to space postings evenly instead.
</div>

[Full source on Forgejo](https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/11_balance/go/chart.go)

---

## Running it

```bash
task example:11
```

Outputs SVG to stdout and serves at `http://localhost:1350`.

---

## API reference

| Function | Purpose |
|----------|---------|
| `NewStepChart(opts...)` | Line that holds each value until the next point |
| `Add` / `AddXY` / `AddTimeSeries` | Same series API as line charts |
| `WithAxisMode(mode)` | Temporal (default) or Ordinal spacing |

---

[Back to examples](../index.html) | [Source on Forgejo](https://git.bytestone.uk/hum3/gogal/src/branch/main/examples/11_balance) | [API docs](https://pkg.go.dev/git.bytestone.uk/hum3/gogal)
