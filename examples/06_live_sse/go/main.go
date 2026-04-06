//go:build !js

// Example 06: Live SSE
//
// Server-sent events push a re-rendered chart every 2 seconds as
// new data arrives. The browser uses EventSource to replace the
// SVG without a full page reload. Serves on http://localhost:1345.
package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"codeberg.org/hum3/gogal"
)

func main() {
	fmt.Println("Serving live SSE chart at http://localhost:1345")

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		var points []gogal.DataPoint
		tick := time.NewTicker(2 * time.Second)
		defer tick.Stop()

		i := 0
		for {
			select {
			case <-r.Context().Done():
				return
			case now := <-tick.C:
				y := 50 + 20*math.Sin(float64(i)*0.3) + rand.Float64()*10 - 5
				points = append(points, gogal.DataPoint{
					Time:  now,
					Y:     y,
					Label: fmt.Sprintf("%.1f", y),
				})
				// Keep last 30 points
				if len(points) > 30 {
					points = points[len(points)-30:]
				}

				chart := gogal.NewLineChart(
					gogal.WithVariant(gogal.Static),
					gogal.WithTitle("Live Sensor Data"),
					gogal.WithGrid(true),
					gogal.WithSmooth(true),
					gogal.WithTooltips(true),
					gogal.WithTimeFormat("15:04:05"),
					gogal.WithYFormat("%.0f"),
				)
				chart.Add("Sensor", points)

				var buf bytes.Buffer
				chart.Render(&buf)
				// SSE data lines: replace newlines with \ndata:
				svg := strings.ReplaceAll(buf.String(), "\n", "\ndata:")
				fmt.Fprintf(w, "data:%s\n\n", svg)
				flusher.Flush()
				i++
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Live SSE</title></head>
<body style="font-family: system-ui; padding: 2em; max-width: 900px; margin: 0 auto;">
<h1>Live Chart (SSE)</h1>
<p>Chart updates every 2 seconds via Server-Sent Events.</p>
<div id="chart"><p>Connecting...</p></div>
<script>
const source = new EventSource('/events');
const chart = document.getElementById('chart');
source.onmessage = function(e) {
    chart.innerHTML = e.data;
};
source.onerror = function() {
    chart.innerHTML = '<p>Connection lost. Reload to reconnect.</p>';
};
</script>
</body></html>`)
	})

	http.ListenAndServe(":1345", nil)
}
