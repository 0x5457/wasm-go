package wasm_go

import "fmt"

type opLocalGet struct {
	localIdx int
}

func (o *opLocalGet) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	v, ok := valueStack.Get(frame.sp, o.localIdx)
	if !ok {
		return fmt.Errorf("local variable[%d] not found", o.localIdx)
	}
	valueStack.Push(*v)

	frame.NextStep()
	return nil
}

type opLocalSet struct {
	localIdx int
}

func (o *opLocalSet) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	v, _ := valueStack.Pop()
	valueStack.Set(frame.sp, o.localIdx, v)
	frame.NextStep()
	return nil
}

type opLocalTee struct {
	localIdx int
}

func (o *opLocalTee) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	v, _ := valueStack.Top()
	valueStack.Set(frame.sp, o.localIdx, *v)
	frame.NextStep()
	return nil
}

type opGlobalGet struct {
	globalIdx int
}

func (o *opGlobalGet) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	globalAddr := frame.mod.globalAddrs[o.globalIdx]
	global := store.globals[globalAddr]
	valueStack.Push(global.value)
	frame.NextStep()
	return nil
}

type opGlobalSet struct {
	globalIdx int
}

func (o *opGlobalSet) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	frame, _ := frameStack.Top()
	globalAddr := frame.mod.globalAddrs[o.globalIdx]
	global := store.globals[globalAddr]
	if global.globalType.mut == const_ {
		return fmt.Errorf("global[%d] is a const value", o.globalIdx)
	}
	v, _ := valueStack.Top()
	if global.globalType.valueType != v.ValType {
		return fmt.Errorf("global[%d] and value types do not match ", o.globalIdx)
	}

	global.value = *v
	frame.NextStep()
	return nil
}
