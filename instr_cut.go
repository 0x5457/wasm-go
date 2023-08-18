package wasm_go

// wrap ∣ extend ∣ trunc ∣ convert ∣ demote ∣ promote ∣ reinterpret
type opCut struct {
}

func (o *opCut) exec(frameStack *stack[frame], valueStack *stack[Value], store *store) error {
	return nil
}
