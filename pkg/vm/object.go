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

type ObjectFunction struct {
	next Obj

	arity int
	chunk *Chunk
	name  *ObjectString
}

func (of *ObjectFunction) Type() ObjType {
	return ObjFunction
}
func (of *ObjectFunction) SetNext(next Obj) {
	of.next = next
}
func (of *ObjectFunction) GetNext() Obj {
	return of.next
}

type ObjectString struct {
	String string
	Hash   uint32
	next   Obj
}

var _ Obj = (*ObjectString)(nil)
var _ Obj = (*ObjectFunction)(nil)

func copyString(chars string) *ObjectString {
	hash := hashString(chars)
	if interned := vm.Strings.FindString(chars, hash); interned != nil {
		return interned
	}
	// We don't do much with strings because Go already stores strings in a
	// convenient manner. We do however try to stay true to the book and create
	// an actual copy of the string. This is just so we can implement our own
	// Garbage collector later for practice.
	return allocateString(strings.Clone(chars), hash)
}

func printFunction(function *ObjectFunction) {
	fmt.Printf("<fn %s>", function.name.String)
}

func takeString(s string) *ObjectString {
	hash := hashString(s)
	if interned := vm.Strings.FindString(s, hash); interned != nil {
		return interned
	}
	return allocateString(s, hash)
}

func allocateObject(obj Obj) {
	obj.SetNext(vm.Objects)
	vm.Objects = obj
}

func newFunction() *ObjectFunction {
	fn := &ObjectFunction{
		arity: 0,
		name:  nil,
		chunk: InitChunk(),
	}
	allocateObject(fn)
	return fn
}

func allocateString(s string, hash uint32) *ObjectString {
	os := &ObjectString{
		String: s,
		Hash:   hash,
	}
	allocateObject(os)
	vm.Strings.Set(os, NilValue())
	return os
}

func hashString(s string) uint32 {
	var hash uint32 = 2166136261
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		hash = hash ^ uint32(b[i])
		hash = hash * 16777619
	}
	return hash
}

func printObject(value Value) {
	switch value.Obj.Type() {
	case ObjString:
		fmt.Print(AsGoString(value))
	case ObjFunction:
		printFunction(AsFunction(value))
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

func IsFunction(v Value) bool {
	return v.IsObject() && v.Obj.Type() == ObjFunction
}

func IsString(v Value) bool {
	return v.IsObject() && v.Obj.Type() == ObjString
}

func AsFunction(value Value) *ObjectFunction {
	if o, ok := value.Obj.(*ObjectFunction); ok {
		return o
	}
	return nil
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
