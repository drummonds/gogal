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
<strong>06</strong> |
<a href="../07_scatter/">07</a> |
<a href="../08_wasm/">08</a> |
<a href="../09_dashboard/">09</a>
</div>

# 06 — Live SSE

Real-time chart updates via Server-Sent Events. The server pushes a re-rendered SVG every 2 seconds as new data arrives. The browser replaces the chart with a single `EventSource` listener — no polling, no WebSocket, no HTMX.

<div class="columns">
<div class="column is-8">
<figure class="screenshot">
<img src="06_live.svg" alt="Live chart with 30 data points">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Snapshot after 30 data points (chart updates every 2s)</figcaption>
</figure>
</div>
<div class="column">
<div class="buttons">
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/06_live_sse" class="button is-light">Source on Codeberg</a>
</div>
</div>
</div>

---

## The pattern

### Server: SSE endpoint

The `/events` handler keeps the connection open and pushes a new SVG each tick. The data window slides — only the last 30 points are kept.

```go
http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")

    flusher := w.(http.Flusher)
    tick := time.NewTicker(2 * time.Second)
    defer tick.Stop()

    var points []gogal.DataPoint
    for {
        select {
        case <-r.Context().Done():
            return
        case now := <-tick.C:
            points = append(points, gogal.DataPoint{Time: now, Y: newValue()})
            if len(points) > 30 {
                points = points[len(points)-30:]
            }

            chart := gogal.NewLineChart(
                gogal.WithVariant(gogal.Static),
                gogal.WithTitle("Live Sensor Data"),
                gogal.WithGrid(true),
                gogal.WithSmooth(true),
            )
            chart.Add("Sensor", points)

            var buf bytes.Buffer
            chart.Render(&buf)
            // SSE format: each line prefixed with "data:"
            svg := strings.ReplaceAll(buf.String(), "\n", "\ndata:")
            fmt.Fprintf(w, "data:%s\n\n", svg)
            flusher.Flush()
        }
    }
})
```

<div class="annotation">
<strong>Full SVG replacement</strong> — each SSE event contains a complete SVG. The chart is small (~2-5 KB) so this is efficient and avoids the complexity of incremental DOM patching. The browser simply replaces innerHTML.
</div>

<div class="annotation">
<strong>Sliding window</strong> — keeping the last 30 points prevents unbounded memory growth and keeps the chart readable. The time axis auto-scales to the visible window.
</div>

### Client: EventSource

```html
<div id="chart"></div>
<script>
const source = new EventSource('/events');
source.onmessage = function(e) {
    document.getElementById('chart').innerHTML = e.data;
};
</script>
```

<div class="annotation">
<strong>Three lines of JavaScript</strong> — <code>EventSource</code> handles reconnection automatically. If the server restarts, the browser reconnects and the chart resumes.
</div>

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/06_live_sse/go/main.go)

---

## Running it

```bash
task example:06
```

Serves at `http://localhost:1345`. The chart starts empty and fills in over ~60 seconds.

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/06_live_sse) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
