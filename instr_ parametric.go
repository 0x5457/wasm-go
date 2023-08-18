package wasm_go

type opSelect struct{}

func (o *opSelect) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	c, _ := valueStack.Pop()
	v1, _ := valueStack.Pop()
	v2, _ := valueStack.Pop()

	if c.I32() == 0 {
		valueStack.Push(v1)
	} else {
		valueStack.Push(v2)
	}

	frame.NextStep()
	return nil
}

type opDrop struct{}

func (o *opDrop) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	valueStack.Pop()
	frame.NextStep()
	return nil
}
