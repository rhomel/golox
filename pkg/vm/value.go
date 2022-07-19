package vm

type Value float64

type ValueArray struct {
	values []Value
}

func InitValueArray() *ValueArray {
	return &ValueArray{
		values: make([]Value, 0),
	}
}

func (va *ValueArray) Write(value Value) {
	va.values = append(va.values, value)
}

func (va *ValueArray) Count() int {
	return len(va.values)
}
