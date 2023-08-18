package tests

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"wasm_go"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	runTest(t, "./suite/json/address.json")
}

func TestBlock(t *testing.T) {
	runTest(t, "./suite/json/block.json")
}

func TestI32(t *testing.T) {
	runTest(t, "./suite/json/i32.json")
}

func TestI64(t *testing.T) {
	runTest(t, "./suite/json/i64.json")
}

func TestF32(t *testing.T) {
	runTest(t, "./suite/json/f32.json")
}

func TestF64(t *testing.T) {
	runTest(t, "./suite/json/f64.json")
}

func runTest(t *testing.T, jsonPath string) {
	config := loadConfigFromFile(jsonPath)
	dir, _ := filepath.Split(jsonPath)
	var i wasm_go.Interpreter
	for _, cmd := range config.Commands {
		fmt.Println(cmd.Line)
		switch cmd.Type {
		case "module":
			wasm, err := os.ReadFile(path.Join(dir, cmd.Filename))
			assert.NoError(t, err)
			i, err = wasm_go.NewInterpreter(wasm)
			assert.NoError(t, err)
		case "assert_return":
			switch cmd.Action.Type {
			case "invoke":
				ret, err := invoke(t, &i, cmd)
				assert.NoError(t, err)
				expected := wasmValue(cmd.Expected)
				if len(cmd.Expected) > 0 && (cmd.Expected[0].Value == "nan:canonical" || cmd.Expected[0].Value == "nan:arithmetic") {
					var isNaN bool
					if cmd.Expected[0].Type == "f32" {
						isNaN = math.IsNaN(float64(ret[0].F32()))
					} else {
						isNaN = math.IsNaN(ret[0].F64())
					}
					assert.Truef(t, isNaN, "line: %d ret[0] should be NaN but got %f", cmd.Line, ret[0].F32())
				} else {
					eq := assert.Equal(t, expected, ret, "line: %d; %s(%s) expected: %s, got: %s", cmd.Line, cmd.Action.Field, goValue(wasmValue(cmd.Action.Args)), goValue(expected), goValue(ret))
					if !eq {
						return
					}
				}
			default:
				t.Errorf("unknown action: %s", cmd.Action.Type)
			}
		case "assert_trap":
			switch cmd.Action.Type {
			case "invoke":
				_, err := invoke(t, &i, cmd)
				if assert.NotNil(t, err, "line: %d; %s(%s) expected tarp: %s, got: nil", cmd.Line, cmd.Action.Field, cmd.Action.Args, cmd.Text) {
					assert.Equal(t, cmd.Text, err.Error(), "line: %d; %s(%s) expected tarp: %s, got: %s", cmd.Line, cmd.Action.Field, cmd.Action.Args, cmd.Text, err.Error())
				}
			default:
				t.Errorf("unknown action: %s", cmd.Action.Type)
			}
		}
	}
}

type config struct {
	SourceFilename string    `json:"source_filename"`
	Commands       []command `json:"commands"`
}

type command struct {
	Type       string      `json:"type"`
	Line       int         `json:"line"`
	Filename   string      `json:"filename"`
	Name       string      `json:"name"`
	Action     cmdAction   `json:"action"`
	Text       string      `json:"text"`
	ModuleType string      `json:"module_type"`
	Expected   []valueInfo `json:"expected"`
}

type cmdAction struct {
	Type     string      `json:"type"`
	Module   string      `json:"module"`
	Field    string      `json:"field"`
	Args     []valueInfo `json:"args"`
	Expected []valueInfo `json:"expected"`
}

type valueInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func loadConfigFromFile(filename string) config {
	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var cfg config
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func wasmValue(vs []valueInfo) []wasm_go.Value {
	values := make([]wasm_go.Value, len(vs))
	for i, value := range vs {
		v, _ := strconv.ParseUint(value.Value, 10, 0)
		switch value.Type {
		case "i32":
			values[i] = wasm_go.ValueFrom(int32(v), wasm_go.I32)
		case "i64":
			values[i] = wasm_go.ValueFrom(int64(v), wasm_go.I64)
		case "f32":
			values[i] = wasm_go.ValueFrom(uint32(v), wasm_go.F32)
		case "f64":
			values[i] = wasm_go.ValueFrom(v, wasm_go.F64)
		}
	}
	return values
}

func goValue(values []wasm_go.Value) []any {
	vs := make([]any, len(values))
	for i, value := range values {
		switch value.ValType {
		case wasm_go.I32:
			vs[i] = value.I32()
		case wasm_go.I64:
			vs[i] = value.I64()
		case wasm_go.F32:
			vs[i] = value.F32()
		case wasm_go.F64:
			vs[i] = value.F64()
		}
	}
	return vs
}

func invoke(t *testing.T, i *wasm_go.Interpreter, cmd command) ([]wasm_go.Value, error) {
	fn, err := i.GetFunc(cmd.Action.Field)
	if err != nil {
		return nil, err
	}
	return fn(wasmValue(cmd.Action.Args))
}
