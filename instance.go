package wasm_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var errOutOfBounds = errors.New("out of bounds memory access")

const DEFAULT_MEM_ADDR_IDX = 0

// https://webassembly.github.io/spec/core/exec/runtime.html#module-instances
type moduleInst struct {
	signatures  []funcType
	funcAddrs   []uint32
	tableAddrs  []uint32
	memAddrs    []uint32
	globalAddrs []uint32
	elemAddrs   []uint32
	dataAddrs   []uint32
	exports     []exportInst
}

func (m *moduleInst) defaultMemAddr() uint32 {
	return m.memAddrs[DEFAULT_MEM_ADDR_IDX]
}

// https://webassembly.github.io/spec/core/exec/runtime.html#function-instances
type funcInst struct {
	funcType     funcType
	kind         funcKind
	internalFunc internalFuncInst
	externalFunc externalFuncInst
}
type funcKind uint8

const (
	internalFunc funcKind = 0x01
	externalFunc funcKind = 0x02
)

type internalFuncInst struct {
	module *moduleInst
	code   function
}

type externalFuncInst struct {
	// TODO:
}

// https://webassembly.github.io/spec/core/exec/runtime.html#table-instances
type tableInst struct {
	tableType
	elems []ref
}

const PAGE_SIZE int = 65536

type memInst struct {
	memType memType
	data    []byte
}

func (m *memInst) size() int {
	return len(m.data)
}

func (m *memInst) pages() int {
	return int(m.size() / PAGE_SIZE)
}

func (m *memInst) grow(n int) error {
	toPages := m.pages() + n
	if m.memType.limits.Max >= 0 && toPages > int(m.memType.limits.Max) {
		return fmt.Errorf("memory page is overflow. max is %d, grow size is %d", toPages, m.memType.limits.Max)
	}
	data := make([]byte, toPages*PAGE_SIZE)
	copy(data, m.data)
	m.data = data
	return nil
}

func (m *memInst) load8(addr, align int32) (uint8, error) {
	if addr < 0 || addr+1 > int32(len(m.data)) {
		return 0, errOutOfBounds
	}
	var v uint8
	err := binary.Read(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, &v)
	return v, err
}

func (m *memInst) load16(addr, align int32) (uint16, error) {
	if addr < 0 || addr+2 > int32(len(m.data)) {
		return 0, errOutOfBounds
	}
	var v uint16
	err := binary.Read(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, &v)
	return v, err
}

func (m *memInst) load32(addr, align int32) (uint32, error) {
	if addr < 0 || addr+4 > int32(len(m.data)) {
		return 0, errOutOfBounds
	}
	var v uint32
	err := binary.Read(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, &v)
	return v, err
}

func (m *memInst) load64(addr, align int32) (uint64, error) {
	if addr < 0 || addr+8 > int32(len(m.data)) {
		return 0, errOutOfBounds
	}
	var v uint64
	err := binary.Read(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, &v)
	return v, err
}

func (m *memInst) store8(addr, align int32, v uint8) error {
	if addr < 0 || addr+1 > int32(len(m.data)) {
		return errOutOfBounds
	}
	return binary.Write(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, v)
}

func (m *memInst) store16(addr, align int32, v uint16) error {
	if addr < 0 || addr+2 > int32(len(m.data)) {
		return errOutOfBounds
	}
	return binary.Write(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, v)
}

func (m *memInst) store32(addr, align int32, v uint32) error {
	if addr < 0 || addr+4 > int32(len(m.data)) {
		return errOutOfBounds
	}
	return binary.Write(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, v)
}

func (m *memInst) store64(addr, align int32, v uint64) error {
	if addr < 0 || addr+8 > int32(len(m.data)) {
		return errOutOfBounds
	}
	return binary.Write(bytes.NewBuffer(m.data[addr:]), binary.LittleEndian, v)
}

type globalInst struct {
	globalType globalType
	value      Value
}

// https://webassembly.github.io/spec/core/exec/runtime.html#element-instances
type elemInst struct {
	elemType type_
	elem     []type_
}

// https://webassembly.github.io/spec/core/exec/runtime.html#data-instances
type dataInst struct {
	data []byte
}

type exportInst struct {
	name  string
	value externalVal
}

type Value struct {
	ValType type_
	data    []byte
}

func ValueFrom(v any, t type_) Value {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, v)
	return Value{
		ValType: t,
		data:    buffer.Bytes(),
	}
}

func ValueFromI32(v int32) Value {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, v)
	return Value{
		ValType: I32,
		data:    buffer.Bytes(),
	}
}

func ValueFromI64(v int64) Value {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, v)
	return Value{
		ValType: I64,
		data:    buffer.Bytes(),
	}
}

func ValueFromF32(v float32) Value {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, v)
	return Value{
		ValType: F32,
		data:    buffer.Bytes(),
	}
}

func ValueFromF64(v float64) Value {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, v)
	return Value{
		ValType: F64,
		data:    buffer.Bytes(),
	}
}

func (v *Value) F32() float32 {
	var f float32
	binary.Read(bytes.NewReader(v.data), binary.LittleEndian, &f)
	return f
}

func (v *Value) F64() float64 {
	var u float64
	binary.Read(bytes.NewReader(v.data), binary.LittleEndian, &u)
	return u
}

func (v *Value) I32() int32 {
	var i int32
	binary.Read(bytes.NewReader(v.data), binary.LittleEndian, &i)
	return i
}
func (v *Value) I64() int64 {
	var i int64
	binary.Read(bytes.NewReader(v.data), binary.LittleEndian, &i)
	return i
}

func (v *Value) Bool() bool {
	if v.ValType == I32 {
		return int32(0) != v.I32()
	} else if v.ValType == I64 {
		return int64(0) != v.I64()
	}
	panic("value must be int type")
}

type refKind uint8

const (
	refExtern refKind = 0x00
	refFunc   refKind = 0x01
	refNull   refKind = 0x03
)

type ref struct {
	addr int
	kind refKind
}

func (r *ref) isNull() bool {
	return r.addr == 0
}

type externalVal struct {
	kind exportImportKind
	idx  uint32
}
