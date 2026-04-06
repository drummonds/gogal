//go:build js && wasm

package main

import (
	"math"
	"syscall/js"

	"codeberg.org/hum3/gogal"
)

func render() string {
	var points []gogal.DataPoint
	for i := 0; i < 50; i++ {
		x := float64(i)
		y := math.Sin(x*0.3)*10 + 20 + math.Sin(x*1.1)*3
		points = append(points, gogal.DataPoint{X: x, Y: y})
	}

	regular := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(200, 30),
	)
	regular.Add("temperature", points)
	regSVG, _ := regular.RenderString()

	smooth := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(200, 30),
		gogal.WithSmooth(true),
	)
	smooth.Add("temperature", points)
	smoothSVG, _ := smooth.RenderString()

	wide := gogal.NewLineChart(
		gogal.WithVariant(gogal.Sparkline),
		gogal.WithSize(400, 40),
	)
	wide.Add("temperature", points)
	wideSVG, _ := wide.RenderString()

	return `<p>Current temperature: ` + regSVG + `</p>` +
		`<p>Smooth variant: ` + smoothSVG + `</p>` +
		`<p>Wider: ` + wideSVG + `</p>`
}

func main() {
	js.Global().Set("goRender", js.FuncOf(func(this js.Value, args []js.Value) any {
		return render()
	}))
	js.Global().Call("wasmReady")
	select {} // block forever
}
