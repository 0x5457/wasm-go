package wasm_go

import "fmt"

type Interpreter struct {
	frameStack stack[frame]
	valueStack stack[Value]
	store      store
	mod        moduleInst
}

func NewInterpreter(bytes []byte) (Interpreter, error) {
	p := newParser(bytes)
	m, err := p.parse()
	i := Interpreter{}
	if err != nil {
		return i, err
	}

	store, modInst, err := newStoreAndModuleInst(&i.valueStack, m)
	if err != nil {
		return i, err
	}
	i.store = store
	i.mod = modInst
	return i, nil
}

func (i *Interpreter) Execute() error {
	for !i.frameStack.isEmpty() {
		frame, _ := i.frameStack.Peek(0)
		instr := frame.insts[frame.pc]
		if err := instr.exec(&i.frameStack, &i.valueStack, &i.store); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) GetFunc(fnName string) (func(args []Value) ([]Value, error), error) {
	fnIdx := -1
	for _, export := range i.mod.exports {
		if export.name == fnName {
			if export.value.kind != exportImportKindFunc {
				return nil, fmt.Errorf("%s not a func", fnName)
			}
			fnIdx = int(export.value.idx)
			break
		}
	}
	if fnIdx < 0 {
		return nil, fmt.Errorf("can't find %s func", fnName)
	}

	fnAddr := i.mod.funcAddrs[fnIdx]
	fn := i.store.funcs[fnAddr]
	if fn.kind == externalFunc {
		// TODO: external func
	}

	return func(args []Value) ([]Value, error) {
		i.frameStack.Push(frame{
			pc:    0,
			sp:    i.valueStack.Len(),
			insts: fn.internalFunc.code.body,
			mod:   &i.mod,
		})

		for x := len(args) - 1; x >= 0; x-- {
			i.valueStack.Push(args[x])
		}

		err := i.Execute()
		if err != nil {
			// cleanup valueStack and frameStack
			i.frameStack = stack[frame]{}
			i.valueStack = stack[Value]{}
			return nil, err
		}

		results := make([]Value, len(fn.funcType.results))
		for x := 0; x < len(fn.funcType.results); x++ {
			ret, _ := i.valueStack.Pop()
			results[x] = ret
		}
		return results, nil
	}, nil
}

// https://webassembly.github.io/spec/core/exec/runtime.html#store
type store struct {
	funcs   []funcInst
	tables  []tableInst
	mems    []memInst
	globals []globalInst
	elems   []elemInst
	datas   []dataInst
}

func newStoreAndModuleInst(
	valueStack *stack[Value],
	m module,
) (store, moduleInst, error) {
	s := store{}
	modInst := moduleInst{}

	eval := func(expr expr) (Value, error) {
		frameStack := stack[frame]{}
		// mock frame
		frameStack.Push(frame{
			pc:  0,
			sp:  valueStack.Len(),
			mod: &modInst,
		})
		for _, i := range expr {
			if err := i.exec(&frameStack, valueStack, &s); err != nil {
				return Value{}, err
			}
		}
		v, _ := valueStack.Pop()
		frameStack.Pop()
		return v, nil
	}

	for i, g := range m.globals {
		gv, err := eval(g.initExpr)
		if err != nil {
			return s, modInst, err
		}
		modInst.globalAddrs = append(modInst.globalAddrs, uint32(i))
		s.globals = append(s.globals, globalInst{
			globalType: g.type_,
			value:      gv,
		})
	}

	for i, f := range m.funcs {
		modInst.funcAddrs = append(modInst.funcAddrs, uint32(i))
		s.funcs = append(s.funcs, funcInst{
			funcType: m.types[f.typeIdx],
			kind:     internalFunc,
			internalFunc: internalFuncInst{
				module: &modInst,
				code:   f,
			},
		})
	}

	for i, mem := range m.mems {
		min := mem.limits.Min * uint32(PAGE_SIZE)
		modInst.memAddrs = append(modInst.memAddrs, uint32(i))
		s.mems = append(s.mems, memInst{
			memType: memType{limits: mem.limits},
			data:    make([]byte, min),
		})
	}

	for i := range m.elems {
		modInst.elemAddrs = append(modInst.elemAddrs, uint32(i))
	}
	for i, tab := range m.tables {
		elems := make([]ref, tab.limits.Min)
		modInst.tableAddrs = append(modInst.tableAddrs, uint32(i))
		for _, elem := range m.elems {
			offsetVal, err := eval(elem.offset)
			offset := int(offsetVal.I32())
			if err != nil {
				return s, modInst, err
			}
			if len(elems) <= offset+len(elem.init) {
				originalElems := elems
				elems = make([]ref, offset+len(elem.init))
				copy(elems, originalElems)
			}

			for i, funcIdx := range elem.init {
				elems[i+offset] = ref{addr: int(funcIdx), kind: refFunc}
			}
		}
		s.tables = append(s.tables, tableInst{
			tableType: tableType{
				limits:   tab.limits,
				elemType: tab.elemType,
			},
			elems: elems,
		})
	}

	for i, data := range m.datas {
		modInst.dataAddrs = append(modInst.dataAddrs, uint32(i))
		offsetVal, err := eval(data.offset)
		if err != nil {
			return s, modInst, err
		}
		offset := int(offsetVal.I32())
		mem := s.mems[data.memIdx]
		if len(mem.data) < offset+len(data.init) {
			return s, modInst, fmt.Errorf("data is too large to fit in memory")
		}
		copy(mem.data[offset:], data.init)
	}
	for _, export := range m.exports {
		modInst.exports = append(modInst.exports, exportInst{
			name: export.name,
			value: externalVal{
				kind: export.kind,
				idx:  uint32(export.idx),
			},
		})
	}
	modInst.signatures = m.types
	return s, modInst, nil
}

type frame struct {
	// current instruction position.
	pc int
	// value stack pointer
	sp int
	// function instructions
	insts []instr

	// labels for if, loop, block
	labels stack[label]
	mod    *moduleInst
}

func (f *frame) NextStep() {
	f.pc += 1
}
