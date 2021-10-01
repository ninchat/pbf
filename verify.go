package pbf

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"time"

	"github.com/ninchat/pbf/op"
)

type accessMode uint8

const (
	accessUndefined = accessMode(iota)
	accessScalar
	accessBytes
	accessVector
)

func (m accessMode) String() string {
	switch m {
	case accessUndefined:
		return "undefined"
	case accessScalar:
		return "scalar"
	case accessBytes:
		return "bytes"
	case accessVector:
		return "vector"
	default:
		return fmt.Sprintf("<invalid accessMode value %d>", m)
	}
}

type verifier struct {
	program
	fieldmode  []accessMode
	debugPaths uintptr
}

func verify(p program) (err error) {
	defer func() {
		if x := recover(); x != nil {
			e, _ := x.(error)
			if e == nil {
				panic(x)
			}
			if _, ok := e.(runtime.Error); ok {
				panic(x)
			}
			err = e
		}
	}()

	v := verifier{
		program:   p,
		fieldmode: make([]accessMode, p.fieldcount),
	}

	var debugTime time.Time
	if debugging {
		debugTime = time.Now()
	}

	v.simulate([2]accessMode{}, v.insn())

	if debugging {
		debugf("verify: Simulation time: %v\n", time.Now().Sub(debugTime))
		debugf("verify: Execution paths: %d\n", v.debugPaths)
	}

	return
}

func (v *verifier) simulate(reg [2]accessMode, insn []byte) {
	if debugging {
		v.debugPaths++
	}

	for {
		opcode := op.Code(insn[0])
		insn = insn[1:]

		switch {
		case opcode < 64: // No arguments.
			switch opcode {
			case op.CompareUnsignedLT, op.CompareUnsignedGE, op.CompareUnsignedEQ, op.CompareUnsignedNE, op.CompareUnsignedLE, op.CompareUnsignedGT:
				checkRegs(reg, accessScalar)

			case op.LoadConstScalar0, op.LoadConstScalar1:
				reg[0] = accessScalar

			case op.CompareSignedLT, op.CompareSignedGE, op.CompareSignedEQ, op.CompareSignedNE, op.CompareSignedLE, op.CompareSignedGT:
				checkRegs(reg, accessScalar)

			case op.ReturnFalse, op.ReturnTrue:
				return

			case op.CompareBytesLT, op.CompareBytesGE, op.CompareBytesEQ, op.CompareBytesNE, op.CompareBytesLE, op.CompareBytesGT:
				checkRegs(reg, accessBytes)

			case op.CompareFloatLT, op.CompareFloatGE, op.CompareFloatEQ, op.CompareFloatNE, op.CompareFloatLE, op.CompareFloatGT:
				checkRegs(reg, accessScalar)

			case op.CompareFloatInfPos, op.CompareFloatInfNeg, op.CompareFloatNaN:
				checkReg(reg, op.R0, accessScalar)

			case op.ContainsVarint, op.ContainsZigZag, op.ContainsFixed64, op.ContainsFixed32:
				checkReg(reg, op.R1, accessVector)
				checkReg(reg, op.R0, accessScalar)

			default:
				panicUnknownOpcode(opcode)
			}

		case opcode < 128: // 8-bit argument.
			arg := insn[0]
			insn = insn[1:]

			switch opcode {
			case op.LoadR0FieldScalar:
				v.markField(arg, accessScalar)
				reg[0] = accessScalar

			case op.LoadR1FieldScalar:
				v.markField(arg, accessScalar)
				reg[1] = accessScalar

			case op.LoadR0FieldBytes:
				v.markField(arg, accessBytes)
				reg[0] = accessBytes

			case op.LoadR1FieldBytes:
				v.markField(arg, accessBytes)
				reg[1] = accessBytes

			case op.LoadR0FieldVector:
				v.markField(arg, accessVector)
				reg[0] = accessVector

			case op.LoadR1FieldVector:
				v.markField(arg, accessVector)
				reg[1] = accessVector

			case op.CheckField:
				v.checkFieldIndex(arg)

			default:
				panicUnknownOpcode(opcode)
			}

		case opcode < 192: // 16-bit argument.
			arg := binary.LittleEndian.Uint16(insn)
			insn = insn[2:]

			switch opcode {
			case op.SkipFalse, op.SkipTrue:
				v.simulate(reg, insn[arg:])

			case op.Skip:
				insn = insn[arg:]

			default:
				panicUnknownOpcode(opcode)
			}

		default: // 64-bit argument.
			arg := binary.LittleEndian.Uint64(insn)
			insn = insn[8:]

			switch opcode {
			case op.LoadConstScalar:
				reg[0] = accessScalar

			case op.LoadConstBytes:
				v.checkBytesRef(arg)
				reg[0] = accessBytes

			default:
				panicUnknownOpcode(opcode)
			}
		}
	}
}

func (v *verifier) markField(index uint8, m accessMode) {
	v.checkFieldIndex(index)

	switch v.fieldmode[index] {
	case m:
	case accessUndefined:
		v.fieldmode[index] = m
	default:
		panic(fmt.Errorf("pbf: field #%d accessed as %s and %s", index, v.fieldmode[index], m))
	}
}

func (v *verifier) checkBytesRef(ref uint64) {
	off, n := unpackBytesRef(ref)
	end := uint64(off) + uint64(n)
	if end > uint64(len(v.bytecode)) {
		panic(fmt.Errorf("pbf: invalid bytes reference: %#016x", ref))
	}
}

func (v *verifier) checkFieldIndex(index uint8) {
	if index >= v.fieldcount {
		panic(fmt.Errorf("pbf: field index out of bounds: %d", index))
	}
}

func checkReg(reg [2]accessMode, r op.Reg, m accessMode) {
	if reg[r] != m {
		panic(fmt.Errorf("pbf: %s contains %s but instruction expects %s", r, reg[r], m))
	}
}

func checkRegs(reg [2]accessMode, m accessMode) {
	if reg[1] != m || reg[0] != m {
		panic(fmt.Errorf("pbf: binary %s instruction used with %s in R1 and %s in R0", m, reg[1], reg[0]))
	}
}

func panicUnknownOpcode(opcode op.Code) {
	panic(fmt.Errorf("pbf: unknown opcode: %d", opcode))
}
