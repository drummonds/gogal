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
<strong>08</strong> |
<a href="../09_dashboard/">09</a>
</div>

# 08 — WASM Demo

A full multi-series chart rendered entirely in the browser via WebAssembly. The gogal library compiles to WASM with zero changes — same `NewLineChart`, `Add`, `RenderString` API. No server required.

<div class="columns">
<div class="column is-8">
<figure class="screenshot">
<img src="08_wasm.svg" alt="Multi-series chart rendered via WASM">
<figcaption class="has-text-centered has-text-grey is-size-7 mt-2">Full static chart — identical output whether server-rendered or WASM</figcaption>
</figure>
</div>
<div class="column">
<div class="buttons">
<a href="demo.html" class="button is-primary">Launch WASM Demo</a>
<a target="_blank" href="https://codeberg.org/hum3/gogal/src/branch/main/examples/08_wasm" class="button is-light">Source on Codeberg</a>
</div>
</div>
</div>

---

## The WASM entry point

A `main_wasm.go` file (build-tagged `js && wasm`) exports a `goRender()` function that returns SVG HTML. The JS glue calls it once on load.

```go
//go:build js && wasm

package main

import (
    "syscall/js"
    "codeberg.org/hum3/gogal"
)

func render() string {
    chart := gogal.NewLineChart(
        gogal.WithVariant(gogal.Static),
        gogal.WithTitle("Weather — 15 June 2024 (WASM)"),
        gogal.WithGrid(true),
        gogal.WithSmooth(true),
    )
    chart.Add("Temperature (°C)", temp)
    chart.Add("Humidity (%)", humidity)
    svg, _ := chart.RenderString()
    return svg
}

func main() {
    js.Global().Set("goRender", js.FuncOf(func(this js.Value, args []js.Value) any {
        return render()
    }))
    js.Global().Call("wasmReady")
    select {}
}
```

<div class="annotation">
<strong>Same API, different target</strong> — the chart code is identical to server-rendered examples. Only the entry point differs: instead of <code>Render(w)</code> to an HTTP response, you call <code>RenderString()</code> and return it to JavaScript.
</div>

<div class="annotation">
<strong>Build tags</strong> — <code>main.go</code> has <code>//go:build !js</code> (server mode) and <code>main_wasm.go</code> has <code>//go:build js && wasm</code>. The Go compiler picks the right file based on <code>GOOS=js GOARCH=wasm</code>.
</div>

### Building

```bash
GOOS=js GOARCH=wasm go build -o main.wasm .
```

Or via the Taskfile:

```bash
task docs:build-wasm
```

<div class="annotation">
<strong>Binary size</strong> — a gogal WASM binary is ~2.6 MB. Since gogal has zero dependencies, the binary is mostly the Go runtime. TinyGo could reduce this further (not yet tested).
</div>

### JavaScript glue

```javascript
window.wasmReady = function() {
    document.getElementById('output').innerHTML = goRender();
};

const go = new Go();
WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject)
    .then(result => go.run(result.instance));
```

[Full source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/08_wasm/go/)

---

## Running it

Server mode:

```bash
task example:08
```

WASM mode: open `docs/08_wasm/demo.html` in a browser (needs an HTTP server for WASM loading).

---

[Back to examples](../index.html) | [Source on Codeberg](https://codeberg.org/hum3/gogal/src/branch/main/examples/08_wasm) | [API docs](https://pkg.go.dev/codeberg.org/hum3/gogal)
