package main

import (
	"fmt"
	"wasm_go"

	"github.com/bytecodealliance/wasmtime-go/v9"
)

func main() {
	wasm, err := wasmtime.Wat2Wasm(`
		(module
			(func (param i32) (param i32) (result i32)
				local.get 0
				local.get 1
				i32.add
			)
			(export "add" (func 0))
		)
	`)
	if err != nil {
		panic(err)
	}
	i, err := wasm_go.NewInterpreter(wasm)
	if err != nil {
		panic(err)
	}
	addFn, err := i.GetFunc("add")
	if err != nil {
		panic(err)
	}
	ret, err := addFn([]wasm_go.Value{
		wasm_go.ValueFromI32(1),
		wasm_go.ValueFromI32(1),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("add(1 + 1) = ", ret[0].I32())
}
