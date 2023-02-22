package main

import "strings"

var basic = map[string]typeSettings{
	"Bool": {
		Doc: strings.TrimSpace(`
// Bool is the representation of the Python bool type. There are only two
// possible values for a Bool, True and False. Every True value refers to the
// same instance, and every False value refers to the same value.
		`),
		Check: true,
		New:   false,
		Type:  true,
	},
	"Module": {
		Doc: strings.TrimSpace(`
// Bool is the representation of the Python module type.
		`),
		Check: true,
		New:   true,
		Type:  true,
	},
	"ExceptionClass": {
		Doc: strings.TrimSpace(`
// ExceptionClass is the representation of the Python exception type. All
// exceptions are instances of this type.
//
// The Err and ErrV methods can be used to generate an error of a particular
// exception type.
		`),
		Check: false,
		New:   false,
		Type:  false,
	},
	"IteratorMethods": {
		Doc: strings.TrimSpace(`
// IteratorMethods represents an arbitrary Python iterator.
		`),
		Check: false,
		New:   false,
		Type:  false,
	},
	"AsyncIteratorMethods": {
		Doc: strings.TrimSpace(`
// AsyncIteratorMethods represents an arbitrary async Python iterator.
		`),
		Check: false,
		New:   false,
		Type:  false,
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
		Check: false,
		New:   false,
		Type:  false,
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
		Check: false,
		New:   false,
		Type:  false,
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
		Check: false,
		New:   false,
		Type:  false,
	},
	"BufferMethods": {
		Doc: strings.TrimSpace(`
// BufferMethods represents an arbitrary Python type that implements the Buffer
// Protocol.
		`),
		Check: false,
		New:   false,
		Type:  false,
	},
}
