package wasm_go

import "fmt"

type labelKind uint8

const (
	LabelKindIf    labelKind = 0x00
	LabelKindLoop  labelKind = 0x01
	LabelKindBlock labelKind = 0x02
)

type label struct {
	kind    labelKind
	startPc int
	endPc   int
}

type opUnreachable struct{}

func (o *opUnreachable) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	return fmt.Errorf("unreachable")
}

type opNop struct{}

func (o *opNop) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	return nil
}

type opIf struct {
	block block
}

func (o *opIf) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	cond, _ := valueStack.Pop()
	frame, _ := frameStack.Top()

	nextPc, err := nextEndAddr(frame.pc+1, frame.insts)
	if err != nil {
		return err
	}

	if !cond.Bool() {
		// condition is false, skip the if block
		addr, err := nextElseOrEndAddr(frame.pc+1, frame.insts)
		if err != nil {
			return err
		}
		frame.pc = addr
	}
	frame.labels.Push(label{
		kind:    LabelKindIf,
		startPc: frame.pc,
		endPc:   nextPc,
	})
	return nil
}

type opLoop struct {
	block block
}

func (o *opLoop) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	nextPc, err := nextEndAddr(frame.pc+1, frame.insts)
	if err != nil {
		return err
	}
	frame.labels.Push(label{
		kind:    LabelKindLoop,
		startPc: frame.pc,
		endPc:   nextPc,
	})
	return nil
}

type opBlock struct {
	block block
}

func (o *opBlock) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	nextPc, err := nextEndAddr(frame.pc+1, frame.insts)
	if err != nil {
		return err
	}
	frame.labels.Push(label{
		kind:    LabelKindBlock,
		startPc: frame.pc,
		endPc:   nextPc,
	})
	frame.NextStep()
	return nil
}

type opElse struct{}

func (o *opElse) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	label, ok := frame.labels.Pop()
	if !ok {
		return fmt.Errorf("no label found when else instr")
	}
	frame.pc = label.endPc
	return nil
}

type opEnd struct{}

func (o *opEnd) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	label, ok := frame.labels.Pop()
	if !ok {
		// end func
		frameStack.Pop()
	} else {
		// end label
		frame.pc = label.endPc
	}
	// TODO: restore stack
	return nil
}

type opBr struct {
	level int
}

func (o *opBr) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	var err error
	frame.pc, err = br(&frame.labels, valueStack, int(o.level))
	return err
}

type opBrIf struct {
	level int
}

func (o *opBrIf) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	cond, _ := valueStack.Pop()
	frame, _ := frameStack.Top()

	if cond.Bool() {
		var err error
		frame.pc, err = br(&frame.labels, valueStack, int(o.level))
		return err
	}
	frame.NextStep()
	return nil
}

type opBrTable struct {
	labelIdxArr []int
	defaultIdx  int
}

func (o *opBrTable) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	idxValue, _ := valueStack.Pop()
	frame, _ := frameStack.Top()
	idx := int(idxValue.I32())

	level := o.defaultIdx
	if idx < len(o.labelIdxArr) {
		level = o.labelIdxArr[idx]
	}

	var err error
	frame.pc, err = br(&frame.labels, valueStack, level)
	return err
}

type opReturn struct{}

func (o *opReturn) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	return nil
}

type opCall struct{}

func (o *opCall) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	return nil
}

type opCallIndirect struct{}

func (o *opCallIndirect) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	return nil
}

func br(labels *stack[label], valueStack *stack[Value], level int) (int, error) {
	label, ok := labels.Peek(level)
	if !ok {
		return 0, fmt.Errorf("no label found level: %d", level)
	}
	var nextPc int
	if label.kind == LabelKindLoop {
		// jump start of loop
		nextPc = label.startPc
	} else {
		nextPc = label.endPc
	}
	// TODO: restore stack
	return nextPc, nil
}

// nextEndAddr finds the next end address of a block of instructions given the current program counter `pc` and the list of instructions `insts`.
//
// pc: The current program counter.
// insts: The list of instructions.
//
// Returns the index of the next end address and an error if it is not found.
func nextEndAddr(pc int, insts []instr) (int, error) {
	depth := 0
	for ; pc < len(insts); pc++ {
		instr := insts[pc]
		switch instr.(type) {
		case *opIf:
		case *opLoop:
		case *opBlock:
			depth += 1
		case *opEnd:
			if depth == 0 {
				return pc, nil
			} else {
				depth -= 1
			}
		}
	}
	return -1, fmt.Errorf("no end instruction found")
}

// nextElseOrEndAddr returns the address of the next 'else' or 'end' instruction following the if instruction at the given address.
//
// pc: the address of the if instruction to start searching from.
// insts: the slice of instructions to search through.
// (int, error): the address of the next else or end instruction, or -1 and an error if not found.
func nextElseOrEndAddr(pc int, insts []instr) (int, error) {
	depth := 0
	for ; pc < len(insts); pc++ {
		instr := insts[pc]
		switch instr.(type) {
		case *opIf:
			depth += 1
		case *opElse:
			if depth == 0 {
				return pc, nil
			}
		case *opEnd:
			if depth == 0 {
				return pc, nil
			} else {
				depth -= 1
			}
		}
	}
	return -1, fmt.Errorf("no else or end instruction found after instr if")
}
