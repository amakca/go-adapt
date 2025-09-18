package adapt

import (
	"reflect"
)

type stackEditedCopies struct {
	stack    []func()
	addrCopy reflect.Value
}

func newStack() *stackEditedCopies {
	return &stackEditedCopies{
		stack: make([]func(), 0),
	}
}

func (slf *stackEditedCopies) applyChanges() {
	for i := len(slf.stack) - 1; i >= 0; i-- {
		slf.stack[i]()
	}
}

func (slf *stackEditedCopies) add(input, copyInput reflect.Value) {
	slf.stack = append(slf.stack, func() {
		input.Set(copyInput)
	})
}

func (slf *stackEditedCopies) updateCopy(copyInput reflect.Value) {
	slf.addrCopy = copyInput
}
