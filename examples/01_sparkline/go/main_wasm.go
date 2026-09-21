//go:build js && wasm

package main

import "git.bytestone.uk/hum3/lofigui"

func main() { lofigui.RunWASM(model) }
