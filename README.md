# gogal

Pure Go SVG chart library with CSS-only interactivity and HTMX integration.

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

`gogal` generates pure SVG charts server-side with zero JavaScript. Interactivity is provided through CSS (hover tooltips, transitions) and HTMX (legend toggling, axis switching, live updates).

Designed for use with [lofigui](https://codeberg.org/hum3/lofigui) and server-rendered web applications.

## Key Features

- **Pure Go**: Zero dependencies, no cgo, gokrazy compatible
- **Pure SVG**: Server-side rendered, clean markup, responsive viewBox
- **CSS interactivity**: Hover tooltips, transitions, animations — no JavaScript
- **HTMX integration**: Legend toggling, axis switching, zoom controls, SSE live updates
- **Dual axis modes**: Time-proportional and event-proportional X axes
- **Four variants**: Sparkline, Static, Interactive, Live
- **Accessible**: ARIA labels, `<title>`, `<desc>` by default
- **Testable**: Layout computation separated from SVG rendering

## Quick Start

```go
package main

import (
    "os"
    "codeberg.org/hum3/gogal"
)

func main() {
    chart := gogal.NewLineChart(
        gogal.WithTitle("Temperature"),
        gogal.WithSize(800, 400),
    )
    chart.Add("Sensor A", []gogal.DataPoint{
        {X: 0, Y: 20.1},
        {X: 1, Y: 21.3},
        {X: 2, Y: 19.8},
        {X: 3, Y: 22.5},
    })
    chart.Render(os.Stdout)
}
```

## Chart Types

| Type | Sparkline | Static | Interactive | Live |
|------|-----------|--------|-------------|------|
| Line | Yes | Yes | Yes | Yes |
| Bar  | Yes | Yes | Yes | Yes |
| Scatter | Yes | Yes | Yes | Yes |

## Examples

See the [examples/](examples/) directory for progressively complex demos.

## Links

<!-- auto:links -->
| | |
|---|---|
| Documentation | https://h3-gogal.statichost.page/ |
| Source (Codeberg) | https://codeberg.org/hum3/gogal |
| Mirror (GitHub) | https://github.com/drummonds/gogal |
<!-- /auto:links -->
