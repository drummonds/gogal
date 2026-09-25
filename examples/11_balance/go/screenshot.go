//go:build screenshot

package main

import "os"

func main() {
	f, err := os.Create("../../../docs/11_balance/11_balance.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	balanceChart().Render(f)
}
