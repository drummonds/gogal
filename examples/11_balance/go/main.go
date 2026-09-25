//go:build !js && !screenshot

// Example 11: Balance (step chart)
//
// A running account balance: the value holds between transactions and
// jumps at each one. Serves on http://localhost:1350.
package main

import (
	"fmt"
	"net/http"
)

func main() {
	svg, _ := balanceChart().RenderString()
	fmt.Println(svg)

	fmt.Println("Serving balance chart at http://localhost:1350")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html><head><title>gogal - Balance</title></head>
<body style="font-family: system-ui; padding: 2em; max-width: 900px; margin: 0 auto;">
<h1>Balance (step chart)</h1>
`)
		balanceChart().Render(w)
		fmt.Fprint(w, "\n</body></html>")
	})
	http.ListenAndServe(":1350", nil)
}
