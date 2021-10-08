// Package op contains types and constants useful with PBF bytecode's
// instruction section.
package op

import (
	"fmt"
)

// Reg ister.
type Reg byte

// General-purpose registers.
const (
	R0 = Reg(iota)
	R1
)

func (r Reg) String() string {
	switch r {
	case R0:
		return "R0"
	case R1:
		return "R1"
	default:
		return fmt.Sprintf("<invalid op.Reg value %d>", r)
	}
}

// Cmp arison.
type Cmp byte

// Comparisons.
const (
	CmpLT = Cmp(iota)
	CmpGE
	CmpEQ
	CmpNE
	CmpLE
	CmpGT
)

func (cmp Cmp) String() string {
	switch cmp {
	case CmpLT:
		return "<"
	case CmpGE:
		return ">="
	case CmpEQ:
		return "=="
	case CmpNE:
		return "!="
	case CmpLE:
		return "<="
	case CmpGT:
		return ">"
	default:
		return fmt.Sprintf("<invalid op.Cmp value %d>", cmp)
	}
}

// IsValid value?
func (cmp Cmp) IsValid() bool {
	return cmp <= CmpGT
}

// Code for an operation (opcode).  An opcode byte and its arguments (0-8
// bytes) form an instruction.
type Code byte

// Opcodes without arguments.
const (
	CompareUnsignedLT  = Code(iota + 0) // Binary register.     [Cmp]
	CompareUnsignedGE                   // Binary register.     [Cmp]
	CompareUnsignedEQ                   // Binary register.     [Cmp]
	CompareUnsignedNE                   // Binary register.     [Cmp]
	CompareUnsignedLE                   // Binary register.     [Cmp]
	CompareUnsignedGT                   // Binary register.     [Cmp]
	LoadConstScalar0                    // Unary register (R0). [Option]
	LoadConstScalar1                    // Unary register (R0). [Option]
	CompareSignedLT                     // Binary register.     [Cmp]
	CompareSignedGE                     // Binary register.     [Cmp]
	CompareSignedEQ                     // Binary register.     [Cmp]
	CompareSignedNE                     // Binary register.     [Cmp]
	CompareSignedLE                     // Binary register.     [Cmp]
	CompareSignedGT                     // Binary register.     [Cmp]
	ReturnFalse                         // Nullary.             [Option]
	ReturnTrue                          // Nullary.             [Option]
	CompareBytesLT                      // Binary register.     [Cmp]
	CompareBytesGE                      // Binary register.     [Cmp]
	CompareBytesEQ                      // Binary register.     [Cmp]
	CompareBytesNE                      // Binary register.     [Cmp]
	CompareBytesLE                      // Binary register.     [Cmp]
	CompareBytesGT                      // Binary register.     [Cmp]
	_                                   //
	_                                   //
	CompareFloatLT                      // Binary register.     [Cmp]
	CompareFloatGE                      // Binary register.     [Cmp]
	CompareFloatEQ                      // Binary register.     [Cmp]
	CompareFloatNE                      // Binary register.     [Cmp]
	CompareFloatLE                      // Binary register.     [Cmp]
	CompareFloatGT                      // Binary register.     [Cmp]
	CompareFloatInfPos                  // Unary register (R0). [Option]
	CompareFloatInfNeg                  // Unary register (R0). [Option]
	CompareFloatNaN                     // Unary register (R0).
	_                                   //
	_                                   //
	_                                   //
	ContainsVarint                      // Binary register.
	ContainsZigZag                      // Binary register.
	ContainsFixed64                     // Binary register.
	ContainsFixed32                     // Binary register.
)

// Opcodes with a 1-byte argument.
const (
	LoadR0FieldScalar = Code(iota + 64) // Unary register; field index. [Reg]
	LoadR1FieldScalar                   // Unary register; field index. [Reg]
	LoadR0FieldBytes                    // Unary register; field index. [Reg]
	LoadR1FieldBytes                    // Unary register; field index. [Reg]
	LoadR0FieldVector                   // Unary register; field index. [Reg]
	LoadR1FieldVector                   // Unary register; field index. [Reg]
	CheckField                          // Nullary; field index.
)

// Opcodes with a 2-byte argument.
const (
	SkipFalse = Code(iota + 128) // Nullary; instruction offset. [Option]
	SkipTrue                     // Nullary; instruction offset. [Option]
	Skip                         // Nullary; instruction offset.
)

// Opcodes with 8 bytes of argument data.
const (
	LoadConstScalar = Code(iota + 192) // Unary register (R0); immediate value.
	_                                  //
	LoadConstBytes                     // Unary register (R0); bytecode address and length.
)

// Option operand.
func (op Code) Option() bool {
	if op&1 == 0 {
		return false
	}
	return true
}

// Reg of unary operation.
func (op Code) Reg() Reg {
	return Reg(op & 1)
}

// Cmp operation.
func (op Code) Cmp() Cmp {
	return Cmp(op & 7)
}
