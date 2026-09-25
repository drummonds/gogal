//go:build js && wasm

// WASM build of example 05. There is no server to answer HTMX requests,
// so app.js plays the server's part: it calls goRender with the current
// hidden list and swaps the returned fragment into the page.
package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

// renderFragment returns the chart plus its toggle buttons, the same
// fragment the HTTP /chart endpoint serves, with data-series attributes
// in place of hx-get so app.js can wire the clicks.
func renderFragment(hiddenCSV string) string {
	hidden := parseHiddenList(hiddenCSV)
	var b strings.Builder
	svg, _ := renderChart(hidden).RenderString()
	b.WriteString(svg)
	b.WriteString(`<div style="margin-top: 1em;">`)
	for _, name := range seriesNames {
		fmt.Fprintf(&b, `<button style="%s" data-series="%s">%s</button>`,
			toggleButtonStyle(isHidden(hidden, name)), name, name)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func main() {
	js.Global().Set("goRender", js.FuncOf(func(this js.Value, args []js.Value) any {
		hidden := ""
		if len(args) > 0 {
			hidden = args[0].String()
		}
		return renderFragment(hidden)
	}))
	js.Global().Call("wasmReady")
	select {}
}
