package wasm_go

import (
	"errors"
	"math"
	"math/bits"
)

var (
	errIntegerDivideByZero = errors.New("integer divide by zero")
	errIntegerOverflow     = errors.New("integer overflow")
)

// clz | ctz | popcnt
// abs ∣ neg ∣ sqrt ∣ ceil ∣ floor ∣ trunc ∣ nearest
type opUn struct {
	unOpFn func(v Value) Value
}

func (o *opUn) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	v, _ := valueStack.Pop()
	valueStack.Push(o.unOpFn(v))
	frame, _ := frameStack.Top()
	frame.NextStep()
	return nil
}

// https://webassembly.github.io/spec/core/exec/numerics.html#op-iclz
func i32Clz(v Value) Value {
	return ValueFrom(int32(bits.LeadingZeros32(uint32(v.I32()))), I32)
}
func i64Clz(v Value) Value {
	return ValueFrom(int64(bits.LeadingZeros64(uint64(v.I64()))), I64)
}

// https://webassembly.github.io/spec/core/exec/numerics.html#op-ictz
func i32Ctz(v Value) Value {
	return ValueFrom(int32(bits.TrailingZeros32(uint32(v.I32()))), I32)
}
func i64Ctz(v Value) Value {
	return ValueFrom(int64(bits.TrailingZeros64(uint64(v.I64()))), I64)
}

// https://webassembly.github.io/spec/core/exec/numerics.html#xref-exec-numerics-op-ipopcnt-mathrm-ipopcnt-n-i
func i32Popcnt(v Value) Value {
	return ValueFrom(int32(bits.OnesCount32(uint32(v.I32()))), I32)
}
func i64Popcnt(v Value) Value {
	return ValueFrom(int64(bits.OnesCount64(uint64(v.I64()))), I64)
}
func f32Abs(v Value) Value {
	return ValueFrom(float32(math.Abs(float64(v.F32()))), F32)
}
func f64Abs(v Value) Value {
	return ValueFrom(math.Abs(float64(v.F64())), F64)
}

func f32Neg(v Value) Value {
	return ValueFrom(-v.F32(), F32)
}
func f64Neg(v Value) Value {
	return ValueFrom(-v.F64(), F64)
}

func f32Sqrt(v Value) Value {
	return ValueFrom(float32(math.Sqrt(float64(v.F32()))), F32)
}
func f64Sqrt(v Value) Value {
	return ValueFrom(math.Sqrt(float64(v.F64())), F64)
}

func f32Ceil(v Value) Value {
	return ValueFrom(float32(math.Ceil(float64(v.F32()))), F32)
}
func f64Ceil(v Value) Value {
	return ValueFrom(math.Ceil(float64(v.F64())), F64)
}
func f32Floor(v Value) Value {
	return ValueFrom(float32(math.Floor(float64(v.F32()))), F32)
}
func f64Floor(v Value) Value {
	return ValueFrom(math.Floor(float64(v.F64())), F64)
}

func f32Trunc(v Value) Value {
	return ValueFrom(float32(math.Trunc(float64(v.F32()))), F32)
}
func f64Trunc(v Value) Value {
	return ValueFrom(math.Trunc(float64(v.F64())), F64)
}

func nearest(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) > 0.5 {
		t = t + math.Copysign(1, x)
	}
	return t
}
func f32Nearest(v Value) Value {
	return ValueFrom(float32(nearest(float64(v.F32()))), F32)
}

func f64Nearest(v Value) Value {
	return ValueFrom(nearest(v.F64()), F64)
}

func i32Extend8S(v Value) Value {
	return ValueFrom(extendS8_32(v.I32()), I32)
}

func i32Extend16S(v Value) Value {
	return ValueFrom(extendS16_32(v.I32()), I32)
}

func i64Extend8S(v Value) Value {
	return ValueFrom(extendS8_64(v.I64()), I64)
}
func i64Extend16S(v Value) Value {
	return ValueFrom(extendS16_64(v.I64()), I64)
}

func i64Extend32S(v Value) Value {
	return ValueFrom(extendS32_64(v.I64()), I64)
}

// add ∣ sub ∣ mul ∣ div_u | div_s ∣ rem_u | rem_s
// and ∣ or ∣ xor ∣ shl ∣ shr_s | shr_u ∣ rotl ∣ rotr
// div ∣ min ∣ max ∣ copysign
type opBin struct {
	binFn func(a, b Value) (Value, error)
}

func (o *opBin) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	a, _ := valueStack.Pop()
	b, _ := valueStack.Pop()

	ret, err := o.binFn(a, b)
	if err != nil {
		return err
	}
	valueStack.Push(ret)
	frame, _ := frameStack.Top()
	frame.NextStep()
	return nil
}

func i32Add(a, b Value) (Value, error) {
	return ValueFrom(a.I32()+b.I32(), I32), nil
}

func i64Add(a, b Value) (Value, error) {
	return ValueFromI64(a.I64() + b.I64()), nil
}

func f32Add(a, b Value) (Value, error) {
	return ValueFrom(a.F32()+b.F32(), F32), nil
}

func f64Sub(a, b Value) (Value, error) {
	return ValueFrom(a.F64()-b.F64(), F64), nil
}
func i32Sub(a, b Value) (Value, error) {
	return ValueFrom(a.I32()-b.I32(), I32), nil
}

func i64Sub(a, b Value) (Value, error) {
	return ValueFrom(a.I64()-b.I64(), I64), nil
}

func f32Sub(a, b Value) (Value, error) {
	return ValueFrom(a.F32()-b.F32(), F32), nil
}

func f64Add(a, b Value) (Value, error) {
	return ValueFrom(a.F64()+b.F64(), F64), nil
}

func f32Mul(a, b Value) (Value, error) {
	return ValueFrom(a.F32()*b.F32(), F32), nil
}

func f64Mul(a, b Value) (Value, error) {
	return ValueFrom(a.F64()*b.F64(), F64), nil
}

func i32Mul(a, b Value) (Value, error) {
	return ValueFrom(a.I32()*b.I32(), I32), nil
}

func i64Mul(a, b Value) (Value, error) {
	return ValueFrom(a.I64()*b.I64(), I64), nil
}

func f32Div(a, b Value) (Value, error) {
	return ValueFrom(a.F32()/b.F32(), F32), nil
}

func f64Div(a, b Value) (Value, error) {
	return ValueFrom(a.F64()/b.F64(), F64), nil
}

func i32DivU(a, b Value) (Value, error) {
	aI32 := a.I32()
	bI32 := b.I32()
	if bI32 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	return ValueFrom(uint32(aI32)/uint32(bI32), I32), nil
}
func i32DivS(a, b Value) (Value, error) {
	aI32 := a.I32()
	bI32 := b.I32()
	if bI32 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	if aI32 == math.MinInt32 && bI32 == -1 {
		return Value{}, errIntegerOverflow
	}
	return ValueFrom(aI32/bI32, I32), nil
}

func i64DivU(a, b Value) (Value, error) {
	aI64 := a.I64()
	bI64 := b.I64()
	if bI64 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	return ValueFrom(uint64(aI64)/uint64(bI64), I64), nil
}

func i64DivS(a, b Value) (Value, error) {
	aI64 := a.I64()
	bI64 := b.I64()
	if bI64 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	if aI64 == math.MinInt64 && bI64 == -1 {
		return Value{}, errIntegerOverflow
	}
	return ValueFrom(aI64/bI64, I64), nil
}

func i32RemU(a, b Value) (Value, error) {
	aI32 := a.I32()
	bI32 := b.I32()
	if bI32 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	return ValueFrom(uint32(aI32)%uint32(bI32), I32), nil
}
func i32RemS(a, b Value) (Value, error) {
	aI32 := a.I32()
	bI32 := b.I32()
	if bI32 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	return ValueFrom(aI32%bI32, I32), nil
}

func i64RemU(a, b Value) (Value, error) {
	aI64 := a.I64()
	bI64 := b.I64()
	if bI64 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	return ValueFrom(uint64(aI64)%uint64(bI64), I64), nil
}

func i64RemS(a, b Value) (Value, error) {
	aI64 := a.I64()
	bI64 := b.I64()
	if bI64 == 0 {
		return Value{}, errIntegerDivideByZero
	}
	return ValueFrom(aI64%bI64, I64), nil
}

func i32And(a, b Value) (Value, error) {
	return ValueFrom(a.I32()&b.I32(), I32), nil
}

func i64And(a, b Value) (Value, error) {
	return ValueFrom(a.I64()&b.I64(), I64), nil
}

func i32Or(a, b Value) (Value, error) {
	return ValueFrom(a.I32()|b.I32(), I32), nil
}

func i64Or(a, b Value) (Value, error) {
	return ValueFrom(a.I64()|b.I64(), I64), nil
}

func i32Xor(a, b Value) (Value, error) {
	return ValueFrom(a.I32()^b.I32(), I32), nil
}

func i64Xor(a, b Value) (Value, error) {
	return ValueFrom(a.I64()^b.I64(), I64), nil
}

func i32Shl(a, b Value) (Value, error) {
	return ValueFrom(a.I32()<<(uint32(b.I32())%32), I32), nil
}

func i64Shl(a, b Value) (Value, error) {
	return ValueFrom(a.I64()<<(uint64(b.I64())%64), I64), nil
}

func i32ShrS(a, b Value) (Value, error) {
	return ValueFrom(a.I32()>>(uint32(b.I32())%32), I32), nil
}

func i32ShrU(a, b Value) (Value, error) {
	return ValueFrom(uint32(a.I32())>>(uint32(b.I32())%32), I32), nil
}

func i64ShrS(a, b Value) (Value, error) {
	return ValueFrom(a.I64()>>(uint64(b.I64())%64), I64), nil
}

func i64ShrU(a, b Value) (Value, error) {
	return ValueFrom(uint64(a.I64())>>(uint64(b.I64())%64), I64), nil
}

func i32RotL(a, b Value) (Value, error) {
	return ValueFrom(bits.RotateLeft32(uint32(a.I32()), int(b.I32())), I32), nil
}

func i64RotL(a, b Value) (Value, error) {
	return ValueFrom(bits.RotateLeft64(uint64(a.I64()), int(b.I64())), I64), nil
}

func i32RotR(a, b Value) (Value, error) {
	return ValueFrom(rotateRight32(uint32(a.I32()), int(b.I32())), I32), nil
}

func i64RotR(a, b Value) (Value, error) {
	return ValueFrom(rotateRight64(uint64(a.I64()), int(b.I64())), I64), nil
}

func f32Min(a, b Value) (Value, error) {
	aF32 := a.F32()
	bF32 := b.F32()
	if math.IsNaN(float64(aF32)) || math.IsNaN(float64(bF32)) {
		return ValueFrom(float32(math.NaN()), F32), nil
	}
	return ValueFrom(float32(math.Min(float64(aF32), float64(bF32))), F32), nil
}

func f64Min(a, b Value) (Value, error) {
	aF64 := a.F64()
	bF64 := b.F64()
	if math.IsNaN(aF64) || math.IsNaN(bF64) {
		return ValueFrom(math.NaN(), F64), nil
	}
	return ValueFrom(math.Min(aF64, bF64), F64), nil
}

func f32Max(a, b Value) (Value, error) {
	aF32 := a.F32()
	bF32 := b.F32()
	if math.IsNaN(float64(aF32)) || math.IsNaN(float64(bF32)) {
		return ValueFrom(float32(math.NaN()), F32), nil
	}
	return ValueFrom(float32(math.Max(float64(aF32), float64(bF32))), F32), nil
}

func f64Max(a, b Value) (Value, error) {
	aF64 := a.F64()
	bF64 := b.F64()
	if math.IsNaN(aF64) || math.IsNaN(bF64) {
		return ValueFrom(math.NaN(), F64), nil
	}
	return ValueFrom(math.Max(aF64, bF64), F64), nil
}

func f32Copysign(a, b Value) (Value, error) {
	return ValueFrom(math.Copysign(float64(a.F32()), float64(b.F32())), F32), nil
}

func f64Copysign(a, b Value) (Value, error) {
	return ValueFrom(math.Copysign(a.F64(), b.F64()), F64), nil
}

// https://webassembly.github.io/spec/core/exec/instructions.html#t-mathsf-xref-syntax-instructions-syntax-instr-numeric-mathsf-const-c
type opConst struct {
	val Value
}

func (o *opConst) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	valueStack.Push(o.val)
	frame.NextStep()
	return nil
}

// https://webassembly.github.io/spec/core/syntax/instructions.html#syntax-relop
// eq ∣ ne ∣ lt_sx ∣ gt_sx ∣ le_sx ∣ ge_sx
// lt ∣ gt ∣ le ∣ ge
type opRel struct {
	relFn func(a, b Value) bool
}

func (o *opRel) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	a, _ := valueStack.Pop()
	b, _ := valueStack.Pop()

	valueStack.Push(numericBool(o.relFn(a, b)))

	frame, _ := frameStack.Top()
	frame.NextStep()
	return nil
}

func i32Eq(a, b Value) bool {
	return a.I32() == b.I32()
}
func i64Eq(a, b Value) bool {
	return a.I64() == b.I64()
}
func f32Eq(a, b Value) bool {
	return a.F32() == b.F32()
}
func f64Eq(a, b Value) bool {
	return a.F64() == b.F64()
}
func i32Ne(a, b Value) bool {
	return a.I32() != b.I32()
}
func i64Ne(a, b Value) bool {
	return a.I64() != b.I64()
}
func f32Ne(a, b Value) bool {
	return a.F32() != b.F32()
}
func f64Ne(a, b Value) bool {
	return a.F64() != b.F64()
}

func i32LtS(a, b Value) bool {
	return a.I32() < b.I32()
}
func i64LtS(a, b Value) bool {
	return a.I64() < b.I64()
}
func i32LtU(a, b Value) bool {
	return uint32(a.I32()) < uint32(b.I32())
}
func i64LtU(a, b Value) bool {
	return uint64(a.I64()) < uint64(b.I64())
}

func f32Lt(a, b Value) bool {
	return a.F32() < b.F32()
}
func f64Lt(a, b Value) bool {
	return a.F64() < b.F64()
}

func i32GtS(a, b Value) bool {
	return a.I32() > b.I32()
}
func i64GtS(a, b Value) bool {
	return a.I64() > b.I64()
}
func i32GtU(a, b Value) bool {
	return uint32(a.I32()) > uint32(b.I32())
}
func i64GtU(a, b Value) bool {
	return uint64(a.I64()) > uint64(b.I64())
}
func f32Gt(a, b Value) bool {
	return a.F32() > b.F32()
}
func f64Gt(a, b Value) bool {
	return a.F64() > b.F64()
}

func i32LeS(a, b Value) bool {
	return a.I32() <= b.I32()
}
func i64LeS(a, b Value) bool {
	return a.I64() <= b.I64()
}
func i32LeU(a, b Value) bool {
	return uint32(a.I32()) <= uint32(b.I32())
}
func i64LeU(a, b Value) bool {
	return uint64(a.I64()) <= uint64(b.I64())
}

func i32GeS(a, b Value) bool {
	return a.I32() >= b.I32()
}
func i64GeS(a, b Value) bool {
	return a.I64() >= b.I64()
}
func i32GeU(a, b Value) bool {
	return uint32(a.I32()) >= uint32(b.I32())
}
func i64GeU(a, b Value) bool {
	return uint64(a.I64()) >= uint64(b.I64())
}

func f32Le(a, b Value) bool {
	return a.F32() <= b.F32()
}
func f64Le(a, b Value) bool {
	return a.F64() <= b.F64()
}

func f32Ge(a, b Value) bool {
	return a.F32() >= b.F32()
}
func f64Ge(a, b Value) bool {
	return a.F64() >= b.F64()
}

// https://webassembly.github.io/spec/core/syntax/instructions.html#syntax-testop
// eqz
type opTest struct {
	testFn func(v Value) bool
}

func (o *opTest) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	v, _ := valueStack.Pop()
	valueStack.Push(numericBool(o.testFn(v)))
	frame, _ := frameStack.Top()
	frame.NextStep()
	return nil
}

func i32Eqz(v Value) bool {
	return v.I32() == 0
}
func i64Eqz(v Value) bool {
	return v.I64() == 0
}

func rotateRight32(x uint32, k int) uint32 {
	const n = 32
	s := uint(k) & (n - 1)
	return x>>s | x<<(n-s)
}

func rotateRight64(x uint64, k int) uint64 {
	const n = 64
	s := uint(k) & (n - 1)
	return x>>s | x<<(n-s)
}

// https://webassembly.github.io/spec/core/exec/numerics.html#integer-operations
func numericBool(b bool) Value {
	v := int32(0)
	if b {
		v = int32(1)
	}
	return ValueFrom(v, I32)
}
