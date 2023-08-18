package wasm_go

import (
	"errors"
	"fmt"
	"io"
)

var errInvalidWASMBinary = errors.New("invalid wasm binary magic")

const WASM_MAGIC uint32 = 0x6d736100

// https://webassembly.github.io/spec/core/binary/modules.html#sections
type SectionID uint8

const (
	CustomSection   SectionID = 0x00
	TypeSection     SectionID = 0x01
	ImportSection   SectionID = 0x02
	FunctionSection SectionID = 0x03
	TableSection    SectionID = 0x04
	MemorySection   SectionID = 0x05
	GlobalSection   SectionID = 0x06
	ExportSection   SectionID = 0x07
	StartSection    SectionID = 0x08
	ElementSection  SectionID = 0x09
	CodeSection     SectionID = 0x0a
	DataSection     SectionID = 0x0b
)

type parser struct {
	r leb128Reader
}

func newParser(bytes []byte) parser {
	return parser{
		r: leb128Reader{bytes: bytes, pos: 0},
	}
}

// https://webassembly.github.io/spec/core/binary/modules.html#binary-module
func (p *parser) parse() (module, error) {
	m := module{}
	magic, version, err := p.header()
	if err != nil {
		return m, err
	}
	if magic != WASM_MAGIC || version != 1 {
		return m, errInvalidWASMBinary
	}

	for {
		sid, length, err := p.sectionHeader()
		if err == io.EOF {
			break
		}

		if err != nil {
			return m, err
		}

		switch sid {
		case CustomSection:
			m.custom, err = p.customSection(length)
		case TypeSection:
			m.types, err = p.typeSection()
		case ImportSection:
			m.imports, err = p.importSection()
		case FunctionSection:
			m.funcs, err = p.funcSection()
		case TableSection:
			m.tables, err = p.tableSection()
		case MemorySection:
			m.mems, err = p.memorySection()
		case GlobalSection:
			m.globals, err = p.globalSection()
		case ExportSection:
			m.exports, err = p.exportSection()
		case StartSection:
			m.start, err = p.startSection()
		case ElementSection:
			m.elems, err = p.elemSection()
		case CodeSection:
			err = p.codeSection(m.funcs)
		case DataSection:
			m.datas, err = p.dataSection()
		}
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func (p *parser) header() (magic, version uint32, err error) {
	magicBytes, err := p.r.eatBytes(4)
	if err != nil {
		return
	}
	magic = uint32(magicBytes[0]) | (uint32(magicBytes[1]) << 8) | (uint32(magicBytes[2]) << 16) | (uint32(magicBytes[3]) << 24)

	versionBytes, err := p.r.eatBytes(4)
	if err != nil {
		return
	}
	version = uint32(versionBytes[0]) | (uint32(versionBytes[1]) << 8) | (uint32(versionBytes[2]) << 16) | (uint32(versionBytes[3]) << 24)
	return
}

// https://webassembly.github.io/spec/core/binary/modules.html#sections
func (p *parser) sectionHeader() (sid SectionID, length uint32, err error) {
	sidU8, err := p.r.eatU8()
	sid = SectionID(sidU8)
	if err != nil {
		return
	}
	length, err = p.r.eatU32()
	return
}

// https://webassembly.github.io/spec/core/binary/modules.html#custom-section
func (p *parser) customSection(length uint32) (custom, error) {
	c, err := custom{}, error(nil)
	c.name, err = p.name()
	if err != nil {
		return c, err
	}
	c.data, err = p.r.eatBytes(length - (uint32(len(c.name) + 4)))
	return c, err
}

// https://webassembly.github.io/spec/core/binary/modules.html#type-section
func (p *parser) typeSection() ([]funcType, error) {
	var funcTypes []funcType
	count, err := p.r.eatU32()
	if err != nil {
		return funcTypes, err
	}
	funcTypes = make([]funcType, count)
	for i := uint32(0); i < count; i++ {
		ft, err := p.r.eatU8()
		if err != nil {
			return funcTypes, err
		}
		const FUNC_TYPE_LEADING_BYTE = 0x60
		if FUNC_TYPE_LEADING_BYTE != ft {
			return funcTypes, fmt.Errorf("invalid func type %x", ft)
		}
		funcTypes[i] = funcType{}

		// param types
		paramsCount, err := p.r.eatU32()
		if err != nil {
			return funcTypes, err
		}

		for j := uint32(0); j < paramsCount; j++ {
			valType, err := p.r.eatU8()
			if err != nil {
				return funcTypes, err
			}
			funcTypes[i].params = append(funcTypes[i].params, type_(valType))
		}

		// result types
		resultsCount, err := p.r.eatU32()
		if err != nil {
			return funcTypes, err
		}

		for j := uint32(0); j < resultsCount; j++ {
			valType, err := p.r.eatU8()
			if err != nil {
				return funcTypes, err
			}
			funcTypes[i].results = append(funcTypes[i].results, type_(valType))
		}
	}
	return funcTypes, nil
}

// https://webassembly.github.io/spec/core/binary/modules.html#function-section
// The and fields of the respective functions are encoded separately in the code section.
func (p *parser) funcSection() ([]function, error) {
	var funcs []function
	count, err := p.r.eatU32()
	if err != nil {
		return funcs, err
	}

	funcs = make([]function, count)
	for i := uint32(0); i < count; i++ {
		idx, err := p.r.eatU32()
		if err != nil {
			return funcs, err
		}
		funcs[i].typeIdx = idx
	}
	return funcs, nil
}

// https://webassembly.github.io/spec/core/binary/modules.html#table-section
func (p *parser) tableSection() ([]table, error) {
	var tables []table
	count, err := p.r.eatU32()
	if err != nil {
		return tables, err
	}
	tables = make([]table, count)
	for i := uint32(0); i < count; i++ {
		tables[i], err = p.table()
		if err != nil {
			return tables, err
		}
	}
	return tables, nil
}

func (p *parser) table() (table, error) {
	t := table{}
	elemType, err := p.r.eatU8()
	if err != nil {
		return t, err
	}
	t.elemType = type_(elemType)
	t.limits, err = p.limits()
	return t, err
}

// https://webassembly.github.io/spec/core/binary/modules.html#memory-section
// (memory 1)
func (p *parser) memorySection() ([]mem, error) {
	var mems []mem
	count, err := p.r.eatU32()
	if err != nil {
		return mems, err
	}
	mems = make([]mem, count)
	for i := uint32(0); i < count; i++ {
		mems[i], err = p.memory()
		if err != nil {
			return mems, err
		}
	}
	return mems, nil
}

func (p *parser) memory() (mem, error) {
	m := mem{}
	limits, err := p.limits()
	m.limits = limits
	return m, err
}

// https://webassembly.github.io/spec/core/binary/modules.html#global-section
func (p *parser) globalSection() ([]global, error) {
	var globals []global
	count, err := p.r.eatU32()
	if err != nil {
		return globals, err
	}
	globals = make([]global, count)

	for i := uint32(0); i < count; i++ {
		globals[i].type_, err = p.globalType()
		if err != nil {
			return globals, err
		}

		globals[i].initExpr, err = p.expr()
		if err != nil {
			return globals, err
		}
	}
	return globals, nil
}

// elem ::= { table tableidx, offset expr, init vec(funcidx) }
func (p *parser) elemSection() ([]elem, error) {
	var elems []elem
	count, err := p.r.eatU32()
	if err != nil {
		return elems, err
	}
	elems = make([]elem, count)

	for i := uint32(0); i < count; i++ {
		tableIdx, err := p.r.eatU32()
		if err != nil {
			return elems, err
		}
		elems[i].tableIdx = tableIdx
		elems[i].offset, err = p.expr()
		if err != nil {
			return elems, err
		}
		funcIdxCount, err := p.r.eatU32()
		if err != nil {
			return elems, err
		}

		for j := uint32(0); j < funcIdxCount; j++ {
			funcIdx, err := p.r.eatU32()
			if err != nil {
				return elems, err
			}
			elems[i].init = append(elems[i].init, funcIdx)
		}
	}
	return elems, nil
}

// https://www.w3.org/TR/wasm-core-1/#data-segments%E2%91%A0
// data ::= {data memidx, offset expr, init vec(byte)}
func (p *parser) dataSection() ([]data, error) {
	var datas []data
	count, err := p.r.eatU32()
	if err != nil {
		return datas, err
	}
	datas = make([]data, count)

	for i := uint32(0); i < count; i++ {
		memIdx, err := p.r.eatU32()
		if err != nil {
			return datas, err
		}
		datas[i].memIdx = memIdx
		datas[i].offset, err = p.expr()
		if err != nil {
			return datas, err
		}

		initCount, err := p.r.eatU32()
		if err != nil {
			return datas, err
		}
		datas[i].init, err = p.r.eatBytes(initCount)
		if err != nil {
			return datas, err
		}
	}

	return datas, nil
}

func (p *parser) importSection() ([]import_, error) {
	var imports []import_
	count, err := p.r.eatU32()
	if err != nil {
		return imports, err
	}
	imports = make([]import_, count)

	for i := uint32(0); i < count; i++ {
		imports[i].module, err = p.name()
		if err != nil {
			return imports, err
		}
		imports[i].name, err = p.name()
		if err != nil {
			return imports, err
		}

		kind, err := p.r.eatU8()
		if err != nil {
			return imports, err
		}

		switch exportImportKind(kind) {
		case exportImportKindFunc:
			imports[i].importDesc.typeIdx, err = p.r.eatU32()
		case exportImportKindTable:
			imports[i].importDesc.table, err = p.table()
		case exportImportKindMem:
			imports[i].importDesc.mem, err = p.memory()
		case exportImportKindGlobal:
			imports[i].importDesc.global, err = p.globalType()
		}
		if err != nil {
			return imports, err
		}
	}
	return imports, nil
}

// https://webassembly.github.io/spec/core/binary/modules.html#export-section
func (p *parser) exportSection() ([]export, error) {
	var exports []export
	count, err := p.r.eatU32()
	if err != nil {
		return exports, err
	}
	exports = make([]export, count)

	for i := uint32(0); i < count; i++ {
		exports[i].name, err = p.name()
		if err != nil {
			return exports, err
		}
		kind, err := p.r.eatU8()
		if err != nil {
			return exports, err
		}
		exports[i].kind = exportImportKind(kind)
		idx, err := p.r.eatU32()
		if err != nil {
			return exports, err
		}
		exports[i].idx = idx
	}
	return exports, nil
}

func (p *parser) startSection() (start, error) {
	s, err := p.r.eatU32()
	return start{funcIdx: s}, err
}

// https://webassembly.github.io/spec/core/binary/modules.html#code-section
func (p *parser) codeSection(fs []function) error {
	count, err := p.r.eatU32()
	if err != nil {
		return err
	}
	if count != uint32(len(fs)) {
		return fmt.Errorf("function count mismatch: codeLen(%d) != funcLen(%d)", count, len(fs))
	}

	for i := uint32(0); i < count; i++ {
		// func size
		_, err := p.r.eatU32()
		if err != nil {
			return err
		}
		localsCount, err := p.r.eatU32()
		if err != nil {
			return nil
		}
		fs[i].locals = make([]locals, localsCount)
		for j := uint32(0); j < localsCount; j++ {
			typeCount, err := p.r.eatU32()
			if err != nil {
				return nil
			}
			fs[i].locals[j].count = typeCount
			valType, err := p.r.eatU8()
			if err != nil {
				return nil
			}
			fs[i].locals[j].valType = type_(valType)
		}
		fs[i].body, err = p.expr()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) globalType() (globalType, error) {
	gt := globalType{}
	valueType, err := p.r.eatU8()
	if err != nil {
		return gt, err
	}
	gt.valueType = type_(valueType)
	mut, err := p.r.eatU8()
	gt.mut = mutability(mut)
	return gt, err
}

// https://webassembly.github.io/spec/core/binary/types.html#limits
func (p *parser) limits() (limits, error) {
	var l limits
	limits, err := p.r.eatU32()
	if err != nil {
		return l, err
	}

	l.Min, err = p.r.eatU32()
	if err != nil {
		return l, err
	}
	if limits == 0 {
		// -1 means there is no maximum value
		l.Max = -1
	} else {
		max, err := p.r.eatU32()
		if err != nil {
			return l, err
		}
		l.Max = int32(max)
	}

	return l, nil
}

// https://webassembly.github.io/spec/core/binary/values.html#names
func (p *parser) name() (string, error) {
	length, err := p.r.eatU32()
	if err != nil {
		return "", err
	}
	name, err := p.r.eatString(length)
	return name, err
}

// https://webassembly.github.io/spec/core/binary/instructions.html#expressions
func (p *parser) expr() (expr, error) {
	e := expr{}
	for {
		instr, isEnd, err := p.instr()
		if err != nil {
			return e, err
		}
		e = append(e, instr)
		if isEnd {
			break
		}
	}
	return e, nil
}

func (p *parser) instr() (i instr, isEnd bool, err error) {
	op, err := p.r.eatU8()
	if err != nil {
		return nil, false, err
	}
	switch opcode(op) {
	case opCodeUnreachable:
		i = &opUnreachable{}
	case opCodeNop:
		i = &opNop{}
	case opCodeBlock:
		block, err := p.eatBlock()
		if err != nil {
			return nil, false, err
		}
		i = &opBlock{block}
	case opCodeLoop:
		block, err := p.eatBlock()
		if err != nil {
			return nil, false, err
		}
		i = &opLoop{block}
	case opCodeIf:
		block, err := p.eatBlock()
		if err != nil {
			return nil, false, err
		}
		i = &opIf{block}
	case opCodeElse:
		i = &opElse{}
	case opCodeEnd:
		i = &opEnd{}
		return i, true, nil
	case opCodeBr:
	case opCodeBrIf:
	case opCodeBrTable:
	case opCodeLocalGet:
		idx, err := p.r.eatU32()
		if err != nil {
			return nil, false, err
		}
		i = &opLocalGet{localIdx: int(idx)}
	case opCodeLocalSet:
		idx, err := p.r.eatU32()
		if err != nil {
			return nil, false, err
		}
		i = &opLocalSet{localIdx: int(idx)}
	case opCodeLocalTee:
	case opCodeGlobalGet:
	case opCodeGlobalSet:
	case opCodeCall:
	case opCodeCallIndirect:
	case opCodeI32Const:
		v, err := p.r.eatI32()
		if err != nil {
			return nil, false, err
		}
		i = &opConst{val: ValueFromI32(v)}
	case opCodeI32Eqz:
		i = &opTest{testFn: i32Eqz}
	case opCodeI32Eq:
		i = &opRel{relFn: i32Eq}
	case opCodeI32Ne:
		i = &opRel{relFn: i32Ne}
	case opCodeI32LtS:
		i = &opRel{relFn: i32LtS}
	case opCodeI32LtU:
		i = &opRel{relFn: i32LtU}
	case opCodeI32GtS:
		i = &opRel{relFn: i32GtS}
	case opCodeI32GtU:
		i = &opRel{relFn: i32GtU}
	case opCodeI32LeS:
		i = &opRel{relFn: i32LeS}
	case opCodeI32LeU:
		i = &opRel{relFn: i32LeU}
	case opCodeI32GeS:
		i = &opRel{relFn: i32GeS}
	case opCodeI32GeU:
		i = &opRel{relFn: i32GeU}
	case opCodeI32Add:
		i = &opBin{binFn: i32Add}
	case opCodeI32Sub:
		i = &opBin{binFn: i32Sub}
	case opCodeI32Mul:
		i = &opBin{binFn: i32Mul}
	case opCodeI32Clz:
		i = &opUn{unOpFn: i32Clz}
	case opCodeI32Ctz:
		i = &opUn{unOpFn: i32Ctz}
	case opCodeI32Popcnt:
		i = &opUn{unOpFn: i32Popcnt}
	case opCodeI32DivS:
		i = &opBin{binFn: i32DivS}
	case opCodeI32DivU:
		i = &opBin{binFn: i32DivU}
	case opCodeI32RemS:
		i = &opBin{binFn: i32RemS}
	case opCodeI32RemU:
		i = &opBin{binFn: i32RemU}
	case opCodeI32And:
		i = &opBin{binFn: i32And}
	case opCodeI32Or:
		i = &opBin{binFn: i32Or}
	case opCodeI32Xor:
		i = &opBin{binFn: i32Xor}
	case opCodeI32ShL:
		i = &opBin{binFn: i32Shl}
	case opCodeI32ShrS:
		i = &opBin{binFn: i32ShrS}
	case opCodeI32ShrU:
		i = &opBin{binFn: i32ShrU}
	case opCodeI32RtoL:
		i = &opBin{binFn: i32RotL}
	case opCodeI32RtoR:
		i = &opBin{binFn: i32RotR}
	case opCodeI32Extend8S:
		i = &opUn{unOpFn: i32Extend8S}
	case opCodeI32Extend16S:
		i = &opUn{unOpFn: i32Extend16S}
	case opCodeI64Const:
		v, err := p.r.eatI64()
		if err != nil {
			return nil, false, err
		}
		i = &opConst{val: ValueFromI64(v)}
	case opCodeI64Eqz:
		i = &opTest{testFn: i64Eqz}
	case opCodeI64Eq:
		i = &opRel{relFn: i64Eq}
	case opCodeI64Ne:
		i = &opRel{relFn: i64Ne}
	case opCodeI64LtS:
		i = &opRel{relFn: i64LtS}
	case opCodeI64LtU:
		i = &opRel{relFn: i64LtU}
	case opCodeI64GtS:
		i = &opRel{relFn: i64GtS}
	case opCodeI64GtU:
		i = &opRel{relFn: i64GtU}
	case opCodeI64LeS:
		i = &opRel{relFn: i64LeS}
	case opCodeI64LeU:
		i = &opRel{relFn: i64LeU}
	case opCodeI64GeS:
		i = &opRel{relFn: i64GeS}
	case opCodeI64GeU:
		i = &opRel{relFn: i64GeU}
	case opCodeI64Clz:
		i = &opUn{unOpFn: i64Clz}
	case opCodeI64Ctz:
		i = &opUn{unOpFn: i64Ctz}
	case opCodeI64Popcnt:
		i = &opUn{unOpFn: i64Popcnt}
	case opCodeI64Add:
		i = &opBin{binFn: i64Add}
	case opCodeI64Sub:
		i = &opBin{binFn: i64Sub}
	case opCodeI64Mul:
		i = &opBin{binFn: i64Mul}
	case opCodeI64DivS:
		i = &opBin{binFn: i64DivS}
	case opCodeI64DivU:
		i = &opBin{binFn: i64DivU}
	case opCodeI64RemS:
		i = &opBin{binFn: i64RemS}
	case opCodeI64RemU:
		i = &opBin{binFn: i64RemU}
	case opCodeI64And:
		i = &opBin{binFn: i64And}
	case opCodeI64Or:
		i = &opBin{binFn: i64Or}
	case opCodeI64Xor:
		i = &opBin{binFn: i64Xor}
	case opCodeI64ShL:
		i = &opBin{binFn: i64Shl}
	case opCodeI64ShrS:
		i = &opBin{binFn: i64ShrS}
	case opCodeI64ShrU:
		i = &opBin{binFn: i64ShrU}
	case opCodeI64RtoL:
		i = &opBin{binFn: i64RotL}
	case opCodeI64RtoR:
		i = &opBin{binFn: i64RotR}
	case opCodeI64Extend8S:
		i = &opUn{unOpFn: i64Extend8S}
	case opCodeI64Extend16S:
		i = &opUn{unOpFn: i64Extend16S}
	case opCodeI64Extend32S:
		i = &opUn{unOpFn: i64Extend32S}
	case opCodeF32Const:
	case opCodeF64Const:
	case opCodeF32Eq:
		i = &opRel{relFn: f32Eq}
	case opCodeF32Ne:
		i = &opRel{relFn: f32Ne}
	case opCodeF32Lt:
		i = &opRel{relFn: f32Lt}
	case opCodeF32Gt:
		i = &opRel{relFn: f32Gt}
	case opCodeF32Le:
		i = &opRel{relFn: f32Le}
	case opCodeF32Ge:
		i = &opRel{relFn: f32Ge}
	case opCodeF32Abs:
		i = &opUn{unOpFn: f32Abs}
	case opCodeF32Neg:
		i = &opUn{unOpFn: f32Neg}
	case opCodeF32Ceil:
		i = &opUn{unOpFn: f32Ceil}
	case opCodeF32Floor:
		i = &opUn{unOpFn: f32Floor}
	case opCodeF32Trunc:
		i = &opUn{unOpFn: f32Trunc}
	case opCodeF32Nearest:
		i = &opUn{unOpFn: f32Nearest}
	case opCodeF32Sqrt:
		i = &opUn{unOpFn: f32Sqrt}
	case opCodeF32Add:
		i = &opBin{binFn: f32Add}
	case opCodeF32Sub:
		i = &opBin{binFn: f32Sub}
	case opCodeF32Mul:
		i = &opBin{binFn: f32Mul}
	case opCodeF32Div:
		i = &opBin{binFn: f32Div}
	case opCodeF32Min:
		i = &opBin{binFn: f32Min}
	case opCodeF32Max:
		i = &opBin{binFn: f32Max}
	case opCodeF64Abs:
		i = &opUn{unOpFn: f64Abs}
	case opCodeF64Neg:
		i = &opUn{unOpFn: f64Neg}
	case opCodeF64Ceil:
		i = &opUn{unOpFn: f64Ceil}
	case opCodeF64Floor:
		i = &opUn{unOpFn: f64Floor}
	case opCodeF64Trunc:
		i = &opUn{unOpFn: f64Trunc}
	case opCodeF64Nearest:
		i = &opUn{unOpFn: f64Nearest}
	case opCodeF64Sqrt:
		i = &opUn{unOpFn: f64Sqrt}
	case opCodeF64Add:
		i = &opBin{binFn: f64Add}
	case opCodeF64Sub:
		i = &opBin{binFn: f64Sub}
	case opCodeF64Mul:
		i = &opBin{binFn: f64Mul}
	case opCodeF64Div:
		i = &opBin{binFn: f64Div}
	case opCodeF64Min:
		i = &opBin{binFn: f64Min}
	case opCodeF64Max:
		i = &opBin{binFn: f64Max}
	case opCodeF64Copysign:
		i = &opBin{binFn: f64Copysign}
	case opCodeI32WrapI64:
	case opCodeF64Eq:
		i = &opRel{relFn: f64Eq}
	case opCodeF64Ne:
		i = &opRel{relFn: f64Ne}
	case opCodeF64Lt:
		i = &opRel{relFn: f64Lt}
	case opCodeF64Gt:
		i = &opRel{relFn: f64Gt}
	case opCodeF64Le:
		i = &opRel{relFn: f64Le}
	case opCodeF64Ge:
		i = &opRel{relFn: f64Ge}
	case opCodeF32Copysign:
		i = &opBin{binFn: f32Copysign}
	case opCodeReturn:
		i = &opReturn{}
	case opCodeI32Load:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i32load}
	case opCodeI64Load:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i64load}
	case opCodeF32Load:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: f32load}
	case opCodeF64Load:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: f64load}
	case opCodeI32Load8S:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i32load8S}
	case opCodeI32Load8U:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i32load8U}
	case opCodeI32Load16S:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i32load16S}
	case opCodeI32Load16U:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i32load16U}
	case opCodeI64Load8S:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i64Load8S}
	case opCodeI64Load8U:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i64Load8U}
	case opCodeI64Load16S:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i64load16S}
	case opCodeI64Load16U:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i64load16U}
	case opCodeI64Load32S:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i64load32S}
	case opCodeI64Load32U:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opLoad{align: align, offset: offset, loadFn: i64load32U}
	case opCodeI32Store:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: i32store}
	case opCodeI64Store:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: i64store}
	case opCodeF32Store:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: f32store}
	case opCodeF64Store:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: f64store}
	case opCodeI32Store8:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: i32store8}
	case opCodeI32Store16:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: i32store16}
	case opCodeI64Store8:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: i64store8}
	case opCodeI64Store16:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: i64store16}
	case opCodeI64Store32:
		align, offset, err := p.memoryArgs()
		if err != nil {
			return nil, false, err
		}
		i = &opStore{align: align, offset: offset, storeFn: i64store32}
	case opCodeMemorySize:
		i = &opMemorySize{}
	case opCodeMemoryGrow:
		i = &opMemoryGrow{}
	case opCodeMemoryCopyOrFill:
		kind, err := p.r.eatU8()
		if err != nil {
			return nil, false, err
		}
		if kind == 10 {
			// 0xFC 10:U32 0x00 0x00
			p.r.eatU32()
			p.r.eatU32()
			i = &opMemoryCopy{}
		} else if kind == 11 {
			// 0xFC 11:U32 0x00
			p.r.eatU32()
			i = &opMemoryFill{}
		} else {
			return nil, false, fmt.Errorf("unknown memory copy or fill kind: %d", kind)
		}
	case opCodeSelect:
		i = &opSelect{}
	case opCodeDrop:
		i = &opDrop{}
	case opCodeI32TruncF32S:
	case opCodeI32TruncF32U:
	case opCodeI32TruncF64S:
	case opCodeI32TruncF64U:
	case opCodeI64ExtendI32S:
	case opCodeI64ExtendI32U:
	case opCodeI64TruncF32S:
	case opCodeI64TruncF32U:
	case opCodeI64TruncF64S:
	case opCodeI64TruncF64U:
	case opCodeF32ConvertI32S:
	case opCodeF32ConvertI32U:
	case opCodeF32ConvertI64S:
	case opCodeF32ConvertI64U:
	case opCodeF32DemoteF64:
	case opCodeF64ConvertI32S:
	case opCodeF64ConvertI32U:
	case opCodeF64ConvertI64S:
	case opCodeF64ConvertI64U:
	case opCodeF64PromoteF32:
	case opCodeI32ReinterpretF32:
	case opCodeI64ReinterpretF64:
	case opCodeF32ReinterpretI32:
	case opCodeF64ReinterpretI64:
	}

	return i, false, nil
}

// eat align and offset two i32 values
func (p *parser) memoryArgs() (align, offset int32, err error) {
	align, err = p.r.eatI32()
	if err != nil {
		return
	}
	offset, err = p.r.eatI32()
	if err != nil {
		return
	}
	return
}

func (p *parser) eatBlock() (block, error) {
	blockType, err := p.r.eatU8()
	if err != nil {
		return block{}, err
	}
	if blockType == 0x40 {
		return block{blockType: blockTypeEmpty}, nil
	} else {
		return block{blockType: blockTypeValue, valType: []type_{type_(blockType)}}, nil
	}
}
