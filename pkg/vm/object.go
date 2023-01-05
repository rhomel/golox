package vm

import (
	"fmt"
	"strings"
)

// Go struct embedding is not the same as C embedding--specifically Go does not
// allow structs to be cast to other types of structs. Go does allow interfaces
// to be cast so instead we use Go's interface feature to simulate the same
// behavior.
type Obj interface {
	Type() ObjType
	SetNext(obj Obj)
	GetNext() Obj
}

type ObjectString struct {
	String string
	next   Obj
}

var _ Obj = (*ObjectString)(nil)

func copyString(chars string) *ObjectString {
	// We don't do much with strings because Go already stores strings in a
	// convenient manner. We do however try to stay true to the book and create
	// an actual copy of the string. This is just so we can implement our own
	// Garbage collector later for practice.
	return allocateString(strings.Clone(chars))
}

func takeString(s string) *ObjectString {
	return allocateString(s)
}

func allocateObject(obj Obj) {
	obj.SetNext(vm.Objects)
	vm.Objects = obj
}

func allocateString(s string) *ObjectString {
	os := &ObjectString{
		String: s,
	}
	allocateObject(os)
	return os
}

func printObject(value Value) {
	switch value.Obj.Type() {
	case ObjString:
		fmt.Print(AsGoString(value))
	}
}

func (os *ObjectString) Type() ObjType {
	return ObjString
}

func (os *ObjectString) SetNext(next Obj) {
	os.next = next
}

func (os *ObjectString) GetNext() Obj {
	return os.next
}

func (v Value) IsObject() bool {
	return v.Type == ValObj
}

func IsString(v Value) bool {
	return v.IsObject() && v.Obj.Type() == ObjString
}

func AsString(value Value) *ObjectString {
	if o, ok := value.Obj.(*ObjectString); ok {
		return o
	}
	return nil
}

func AsGoString(value Value) string {
	return AsString(value).String
}

// Go doesn't have macros but the book uses this function as a macro. So to
// save some confusion later we have it here as a function.
func isObjType(value Value, typ ObjType) bool {
	return value.IsObject() && value.AsObject().Type() == typ
}
