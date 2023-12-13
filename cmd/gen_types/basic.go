package main

import "strings"

var basic = map[string]typeSettings{
	"BaseObject": {
		Doc: strings.TrimSpace(`
// BaseObject is the concrete representation of the Python "Object *".  It is
// used less than in the C API, as the Object interface is mostly used when the
// type is not fixed.  Any Object "o" can be turned into a *BaseObject using the
// Base() method (i.e. o.Base() returns a *BaseObject that refers to the same
// underlying Python object as "o").  This allows the Python functions that
// accept any type of object to be defined as methods on *BaseObject.
		`),
		Check:    false,
		New:      true,
		Type:     false,
		Immortal: false,
	},
	"Bool": {
		Doc: strings.TrimSpace(`
// Bool is the representation of the Python bool type. There are only two
// possible values for a Bool, True and False. Every True value refers to the
// same instance, and every False value refers to the same value.
		`),
		Check:    true,
		New:      false,
		Type:     true,
		Immortal: true,
	},
	"NoneObject": {
		Doc: strings.TrimSpace(`
// NoneObject is the type of the None value.  The only value of this type is
// None.
		`),
		Check:    false,
		New:      true,
		Type:     false,
		Immortal: true,
	},
	"Module": {
		Doc: strings.TrimSpace(`
// Bool is the representation of the Python module type.
		`),
		Check:    true,
		New:      true,
		Type:     true,
		Immortal: false,
	},
	"ExceptionClass": {
		Doc: strings.TrimSpace(`
// ExceptionClass is the representation of the Python exception type. All
// exceptions are instances of this type.
//
// The Err and ErrV methods can be used to generate an error of a particular
// exception type.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"IteratorMethods": {
		Doc: strings.TrimSpace(`
// IteratorMethods represents an arbitrary Python iterator.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"AsyncIteratorMethods": {
		Doc: strings.TrimSpace(`
// AsyncIteratorMethods represents an arbitrary async Python iterator.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"IterableMethods": {
		Doc: strings.TrimSpace(`
// IterableMethods represents an arbitrary Python iterable.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"MappingMethods": {
		Doc: strings.TrimSpace(`
// MappingMethods is a concrete realisation of the full set of Mapping Protocol
// methods. A type that implements the "Mapping Protocol" can be turned into a
// MappingMethods instance using AsMapping.
//
// Note that the methods not already implemented on the type itself may return
// an error, as not all methods are implemented by all types that support the
// protocol.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"SequenceMethods": {
		Doc: strings.TrimSpace(`
// SequenceMethods is a concrete realisation of the full set of Sequence
// Protocol methods. A type that implements the "Sequence Protocol" can be
// turned into a SequenceMethods instance using AsSequence.
//
// Note that the methods not already implemented on the type itself may return
// an error, as not all methods are implemented by all types that support the
// protocol.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"NumberMethods": {
		Doc: strings.TrimSpace(`
// NumberMethods is a concrete realisation of the full set of Number Protocol
// methods. A type that implements the "Number Protocol" can be turned into a
// NumberMethods instance using AsNumber.
//
// Note that the methods not already implemented on the type itself may return
// an error, as not all methods are implemented by all types that support the
// protocol.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"BufferMethods": {
		Doc: strings.TrimSpace(`
// BufferMethods represents an arbitrary Python type that implements the Buffer
// Protocol.
		`),
		Check:    false,
		New:      false,
		Type:     false,
		Immortal: false,
	},
	"FrozenSet": {
		Check:    true,
		New:      true,
		Type:     true,
		Immortal: false,
	},
	"Super": {
		Check:    false,
		New:      true,
		Type:     true,
		Immortal: false,
	},
}
