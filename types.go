package wasm_go

// https://webassembly.github.io/spec/core/syntax/modules.html#modules
type module struct {
	custom  custom
	types   []funcType
	funcs   []function
	tables  []table
	mems    []mem
	globals []global
	elems   []elem
	datas   []data
	start   start
	imports []import_
	exports []export
}

type custom struct {
	name string
	data []byte
}

type funcType struct {
	params  []type_
	results []type_
}

type locals struct {
	count   uint32
	valType type_
}
type function struct {
	typeIdx uint32
	locals  []locals
	body    []instr
}

type table struct {
	tableType
}

type global struct {
	type_    globalType
	initExpr expr
}
type mem struct {
	memType
}

// https://www.w3.org/TR/wasm-core-1/#data-segments%E2%91%A0
// data ::= {data memidx,offset expr,init vec(byte)}
type data struct {
	memIdx uint32
	offset expr
	init   []byte
}

type elem struct {
	tableIdx uint32
	offset   expr
	// vec<funcIdx>
	init []uint32
}

type import_ struct {
	module     string
	name       string
	kind       exportImportKind
	importDesc importDesc
}
type importDesc struct {
	typeIdx uint32
	table   table
	mem     mem
	global  globalType
}

type exportImportKind uint8

const (
	exportImportKindFunc   exportImportKind = 0x00
	exportImportKindTable  exportImportKind = 0x01
	exportImportKindMem    exportImportKind = 0x02
	exportImportKindGlobal exportImportKind = 0x03
)

type export struct {
	name string
	kind exportImportKind
	idx  uint32
}

type start struct {
	funcIdx uint32
}

type instr interface {
	exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error
}

type expr []instr

// https://webassembly.github.io/spec/core/binary/types.html#value-types
type type_ uint8

const (
	I32       type_ = 0x70
	I64       type_ = 0x6f
	F32       type_ = 0x7D
	F64       type_ = 0x7C
	V128      type_ = 0x7B
	FuncRef   type_ = 0x70
	ExternRef type_ = 0x6F
)

type limits struct {
	Min uint32
	// -1 means there is no maximum value
	Max int32
}

type tableType struct {
	limits   limits
	elemType type_
}

type memType struct {
	limits limits
}

type mutability uint8

const (
	const_ mutability = 0x00
	var_   mutability = 0x01
)

type globalType struct {
	valueType type_
	mut       mutability
}

type blockType = uint8

const (
	blockTypeEmpty blockType = 0
	blockTypeValue blockType = 1
)

type block struct {
	blockType blockType
	valType   []type_
}

type opcode uint8

const (
	opCodeUnreachable       opcode = 0x00
	opCodeNop               opcode = 0x01
	opCodeBlock             opcode = 0x02
	opCodeLoop              opcode = 0x03
	opCodeIf                opcode = 0x04
	opCodeElse              opcode = 0x05
	opCodeEnd               opcode = 0x0B
	opCodeBr                opcode = 0x0C
	opCodeBrIf              opcode = 0x0D
	opCodeBrTable           opcode = 0x0E
	opCodeLocalGet          opcode = 0x20
	opCodeLocalSet          opcode = 0x21
	opCodeLocalTee          opcode = 0x22
	opCodeGlobalGet         opcode = 0x23
	opCodeGlobalSet         opcode = 0x24
	opCodeCall              opcode = 0x10
	opCodeCallIndirect      opcode = 0x11
	opCodeI32Const          opcode = 0x41
	opCodeI32Eqz            opcode = 0x45
	opCodeI32Eq             opcode = 0x46
	opCodeI32Ne             opcode = 0x47
	opCodeI32LtS            opcode = 0x48
	opCodeI32LtU            opcode = 0x49
	opCodeI32GtS            opcode = 0x4A
	opCodeI32GtU            opcode = 0x4B
	opCodeI32LeS            opcode = 0x4C
	opCodeI32LeU            opcode = 0x4D
	opCodeI32GeS            opcode = 0x4E
	opCodeI32GeU            opcode = 0x4F
	opCodeI32Add            opcode = 0x6a
	opCodeI32Sub            opcode = 0x6b
	opCodeI32Mul            opcode = 0x6c
	opCodeI32Clz            opcode = 0x67
	opCodeI32Ctz            opcode = 0x68
	opCodeI32Popcnt         opcode = 0x69
	opCodeI32DivS           opcode = 0x6D
	opCodeI32DivU           opcode = 0x6E
	opCodeI32RemS           opcode = 0x6F
	opCodeI32RemU           opcode = 0x70
	opCodeI32And            opcode = 0x71
	opCodeI32Or             opcode = 0x72
	opCodeI32Xor            opcode = 0x73
	opCodeI32ShL            opcode = 0x74
	opCodeI32ShrS           opcode = 0x75
	opCodeI32ShrU           opcode = 0x76
	opCodeI32RtoL           opcode = 0x77
	opCodeI32RtoR           opcode = 0x78
	opCodeI32Extend8S       opcode = 0xC0
	opCodeI32Extend16S      opcode = 0xC1
	opCodeI64Const          opcode = 0x42
	opCodeI64Eqz            opcode = 0x50
	opCodeI64Eq             opcode = 0x51
	opCodeI64Ne             opcode = 0x52
	opCodeI64LtS            opcode = 0x53
	opCodeI64LtU            opcode = 0x54
	opCodeI64GtS            opcode = 0x55
	opCodeI64GtU            opcode = 0x56
	opCodeI64LeS            opcode = 0x57
	opCodeI64LeU            opcode = 0x58
	opCodeI64GeS            opcode = 0x59
	opCodeI64GeU            opcode = 0x5A
	opCodeI64Clz            opcode = 0x79
	opCodeI64Ctz            opcode = 0x7A
	opCodeI64Popcnt         opcode = 0x7B
	opCodeI64Add            opcode = 0x7C
	opCodeI64Sub            opcode = 0x7D
	opCodeI64Mul            opcode = 0x7E
	opCodeI64DivS           opcode = 0x7F
	opCodeI64DivU           opcode = 0x80
	opCodeI64RemS           opcode = 0x81
	opCodeI64RemU           opcode = 0x82
	opCodeI64And            opcode = 0x83
	opCodeI64Or             opcode = 0x84
	opCodeI64Xor            opcode = 0x85
	opCodeI64ShL            opcode = 0x86
	opCodeI64ShrS           opcode = 0x87
	opCodeI64ShrU           opcode = 0x88
	opCodeI64RtoL           opcode = 0x89
	opCodeI64RtoR           opcode = 0x8A
	opCodeI64Extend8S       opcode = 0xC2
	opCodeI64Extend16S      opcode = 0xC3
	opCodeI64Extend32S      opcode = 0xC4
	opCodeF32Const          opcode = 0x43
	opCodeF64Const          opcode = 0x44
	opCodeF32Eq             opcode = 0x5B
	opCodeF32Ne             opcode = 0x5C
	opCodeF32Lt             opcode = 0x5D
	opCodeF32Gt             opcode = 0x5E
	opCodeF32Le             opcode = 0x5F
	opCodeF32Ge             opcode = 0x60
	opCodeF32Abs            opcode = 0x8B
	opCodeF32Neg            opcode = 0x8C
	opCodeF32Ceil           opcode = 0x8D
	opCodeF32Floor          opcode = 0x8E
	opCodeF32Trunc          opcode = 0x8F
	opCodeF32Nearest        opcode = 0x90
	opCodeF32Sqrt           opcode = 0x91
	opCodeF32Add            opcode = 0x92
	opCodeF32Sub            opcode = 0x93
	opCodeF32Mul            opcode = 0x94
	opCodeF32Div            opcode = 0x95
	opCodeF32Min            opcode = 0x96
	opCodeF32Max            opcode = 0x97
	opCodeF64Abs            opcode = 0x99
	opCodeF64Neg            opcode = 0x9A
	opCodeF64Ceil           opcode = 0x9B
	opCodeF64Floor          opcode = 0x9C
	opCodeF64Trunc          opcode = 0x9D
	opCodeF64Nearest        opcode = 0x9E
	opCodeF64Sqrt           opcode = 0x9F
	opCodeF64Add            opcode = 0xA0
	opCodeF64Sub            opcode = 0xA1
	opCodeF64Mul            opcode = 0xA2
	opCodeF64Div            opcode = 0xA3
	opCodeF64Min            opcode = 0xA4
	opCodeF64Max            opcode = 0xA5
	opCodeF64Copysign       opcode = 0xA6
	opCodeI32WrapI64        opcode = 0xA7
	opCodeF64Eq             opcode = 0x61
	opCodeF64Ne             opcode = 0x62
	opCodeF64Lt             opcode = 0x63
	opCodeF64Gt             opcode = 0x64
	opCodeF64Le             opcode = 0x65
	opCodeF64Ge             opcode = 0x66
	opCodeF32Copysign       opcode = 0x98
	opCodeReturn            opcode = 0x0f
	opCodeI32Load           opcode = 0x28
	opCodeI64Load           opcode = 0x29
	opCodeF32Load           opcode = 0x2A
	opCodeF64Load           opcode = 0x2B
	opCodeI32Load8S         opcode = 0x2C
	opCodeI32Load8U         opcode = 0x2D
	opCodeI32Load16S        opcode = 0x2E
	opCodeI32Load16U        opcode = 0x2F
	opCodeI64Load8S         opcode = 0x30
	opCodeI64Load8U         opcode = 0x31
	opCodeI64Load16S        opcode = 0x32
	opCodeI64Load16U        opcode = 0x33
	opCodeI64Load32S        opcode = 0x34
	opCodeI64Load32U        opcode = 0x35
	opCodeI32Store          opcode = 0x36
	opCodeI64Store          opcode = 0x37
	opCodeF32Store          opcode = 0x38
	opCodeF64Store          opcode = 0x39
	opCodeI32Store8         opcode = 0x3A
	opCodeI32Store16        opcode = 0x3B
	opCodeI64Store8         opcode = 0x3C
	opCodeI64Store16        opcode = 0x3D
	opCodeI64Store32        opcode = 0x3E
	opCodeMemorySize        opcode = 0x3F
	opCodeMemoryGrow        opcode = 0x40
	opCodeMemoryCopyOrFill  opcode = 0xFC
	opCodeSelect            opcode = 0x1B
	opCodeDrop              opcode = 0x1A
	opCodeI32TruncF32S      opcode = 0xA8
	opCodeI32TruncF32U      opcode = 0xA9
	opCodeI32TruncF64S      opcode = 0xAA
	opCodeI32TruncF64U      opcode = 0xAB
	opCodeI64ExtendI32S     opcode = 0xAC
	opCodeI64ExtendI32U     opcode = 0xAD
	opCodeI64TruncF32S      opcode = 0xAE
	opCodeI64TruncF32U      opcode = 0xAF
	opCodeI64TruncF64S      opcode = 0xB0
	opCodeI64TruncF64U      opcode = 0xB1
	opCodeF32ConvertI32S    opcode = 0xB2
	opCodeF32ConvertI32U    opcode = 0xB3
	opCodeF32ConvertI64S    opcode = 0xB4
	opCodeF32ConvertI64U    opcode = 0xB5
	opCodeF32DemoteF64      opcode = 0xB6
	opCodeF64ConvertI32S    opcode = 0xB7
	opCodeF64ConvertI32U    opcode = 0xB8
	opCodeF64ConvertI64S    opcode = 0xB9
	opCodeF64ConvertI64U    opcode = 0xBA
	opCodeF64PromoteF32     opcode = 0xBB
	opCodeI32ReinterpretF32 opcode = 0xBC
	opCodeI64ReinterpretF64 opcode = 0xBD
	opCodeF32ReinterpretI32 opcode = 0xBE
	opCodeF64ReinterpretI64 opcode = 0xBF
)
