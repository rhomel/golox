package vm

// Main is the entry point for the VM implementation of golox.
func Main() {
	chunk := InitChunk()
	constant := chunk.AddConstant(1.2)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(uint8(constant), 123)
	chunk.Write(OP_RETURN, 123)
	chunk.Disassemble("test chunk")
}
