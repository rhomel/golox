package vm

func freeObjects() {
	object := vm.Objects
	for object != nil {
		next := object.GetNext()
		freeObject(object)
		object = next
	}
	vm.Objects = nil
}

func freeObject(object Obj) {
	object.SetNext(nil)
	// no-op; this is obviously pointless (we rely on Go's garbage collection
	// to deallocate). It is just here to stay similar to the book.
}
