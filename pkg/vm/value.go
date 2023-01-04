package vm

// Another way to do this is to use Go's type switch. But let's try to use a Go
// struct to simulate a C union even though it is memory inefficient. This way
// we can stay in sync with the style of code from the book.
type ValueType int

const (
	ValBool ValueType = iota
	ValNil
	ValNumber
)

type Value struct {
	Type    ValueType
	Boolean bool
	Number  float64
}

func BooleanValue(value bool) Value {
	return Value{
		Type:    ValBool,
		Boolean: value,
	}
}

func NilValue() Value {
	return Value{
		Type: ValNil,
	}
}

func NumberValue(value float64) Value {
	return Value{
		Type:   ValNumber,
		Number: value,
	}
}

func (v Value) AsBool() bool {
	return v.Boolean
}

func (v Value) AsNumber() float64 {
	return v.Number
}

func (v Value) IsBool() bool {
	return v.Type == ValBool
}

func (v Value) IsNil() bool {
	return v.Type == ValNil
}

func (v Value) IsNumber() bool {
	return v.Type == ValNumber
}

type ValueArray struct {
	values []Value
}

func ValuesEqual(a, b Value) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case ValBool:
		return a.AsBool() == b.AsBool()
	case ValNil:
		return true
	case ValNumber:
		return a.AsNumber() == b.AsNumber()
	default:
		return false // unreachable
	}
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
