package wasm_go

// https://webassembly.github.io/spec/core/exec/instructions.html#exec-storen
type opStore struct {
	offset  int32
	align   int32
	storeFn func(m *memInst, addr, align int32, v Value)
}

func (o *opStore) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	mem := store.mems[frame.mod.defaultMemAddr()]
	value, _ := valueStack.Pop()
	addr := value.I32() + o.offset
	o.storeFn(&mem, addr, o.align, value)
	frame.NextStep()
	return nil
}

func i32store(m *memInst, addr, align int32, v Value) {
	m.store32(addr, align, uint32(v.I32()))
}
func i64store(m *memInst, addr, align int32, v Value) {
	m.store64(addr, align, uint64(v.I64()))
}

func f32store(m *memInst, addr, align int32, v Value) {
	m.store32(addr, align, uint32(v.F32()))
}

func f64store(m *memInst, addr, align int32, v Value) {
	m.store64(addr, align, uint64(v.F64()))
}
func i32store8(m *memInst, addr, align int32, v Value) {
	m.store8(addr, align, uint8(v.I32()))
}
func i32store16(m *memInst, addr, align int32, v Value) {
	m.store16(addr, align, uint16(v.I32()))
}
func i64store8(m *memInst, addr, align int32, v Value) {
	m.store8(addr, align, uint8(v.I64()))
}
func i64store16(m *memInst, addr, align int32, v Value) {
	m.store16(addr, align, uint16(v.I64()))
}
func i64store32(m *memInst, addr, align int32, v Value) {
	m.store32(addr, align, uint32(v.I64()))
}

// https://webassembly.github.io/spec/core/exec/instructions.html#exec-loadn
type opLoad struct {
	align  int32
	offset int32
	loadFn func(m *memInst, addr, align int32) (Value, error)
}

func (o *opLoad) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	mem := store.mems[frame.mod.defaultMemAddr()]
	baseAddr, _ := valueStack.Pop()
	baseAddrI32 := baseAddr.I32()
	if baseAddrI32 < 0 || o.offset < 0 {
		return errOutOfBounds
	}
	addr := baseAddrI32 + o.offset
	value, err := o.loadFn(&mem, addr, o.align)
	if err != nil {
		return err
	}
	valueStack.Push(value)
	frame.NextStep()
	return nil
}

func i32load(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load32(addr, align)
	return ValueFromI32(int32(v)), err
}

func i64load(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load64(addr, align)
	return ValueFromI64(int64(v)), err
}

func f32load(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load32(addr, align)
	return ValueFrom(v, F32), err
}

func f64load(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load64(addr, align)
	return ValueFrom(v, F64), err
}

func i32load8S(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load8(addr, align)
	return ValueFromI32(extendS8_32(int32(v))), err
}

func i32load8U(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load8(addr, align)
	return ValueFromI32(int32(v)), err
}

func i32load16S(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load16(addr, align)
	return ValueFromI32(extendS16_32(int32(v))), err
}

func i32load16U(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load16(addr, align)
	return ValueFromI32(int32(v)), err
}

func i64Load8S(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load8(addr, align)
	return ValueFromI64(extendS8_64(int64(v))), err
}

func i64Load8U(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load8(addr, align)
	return ValueFromI64(int64(v)), err
}

func i64load16S(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load16(addr, align)
	return ValueFromI64(extendS16_64(int64(v))), err
}

func i64load16U(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load16(addr, align)
	return ValueFromI64(int64(v)), err
}

func i64load32S(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load32(addr, align)
	return ValueFromI64(extendS32_64(int64(v))), err
}

func i64load32U(m *memInst, addr, align int32) (Value, error) {
	v, err := m.load32(addr, align)
	return ValueFromI64(int64(v)), err
}

type opMemorySize struct{}

func (o *opMemorySize) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	mem := store.mems[frame.mod.defaultMemAddr()]
	valueStack.Push(ValueFrom(int32(mem.size()), I32))
	frame.NextStep()
	return nil
}

type opMemoryGrow struct{}

func (o *opMemoryGrow) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	mem := store.mems[frame.mod.defaultMemAddr()]

	v, _ := valueStack.Pop()
	currentPages := mem.pages()
	pagesWant := int(v.I32())
	err := mem.grow(pagesWant)
	if err != nil {
		valueStack.Push(ValueFrom(-1, I32))
	} else {
		valueStack.Push(ValueFrom(currentPages, I32))
	}
	frame.NextStep()
	return nil
}

type opMemoryCopy struct {
}

func (o *opMemoryCopy) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	len, _ := valueStack.Pop()
	src, _ := valueStack.Pop()
	dst, _ := valueStack.Pop()
	frame, _ := frameStack.Top()
	mem := store.mems[frame.mod.defaultMemAddr()]
	copy(mem.data[dst.I32():], mem.data[src.I32():src.I32()+len.I32()])
	frame.NextStep()
	return nil
}

// https://webassembly.github.io/spec/core/bikeshed/#-hrefsyntax-instr-memorymathsfmemoryfill%E2%91%A0
type opMemoryFill struct {
}

func (o *opMemoryFill) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	return nil
}

func extendS8_32(v int32) int32 {
	return v << 24 >> 24
}

func extendS8_64(v int64) int64 {
	return v << 56 >> 56
}

func extendS16_32(v int32) int32 {
	return v << 16 >> 16
}
func extendS16_64(v int64) int64 {
	return v << 48 >> 48
}

func extendS32_64(v int64) int64 {
	return v << 32 >> 32
}
