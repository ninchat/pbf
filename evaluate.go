package pbf

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/ninchat/pbf/op"
	"google.golang.org/protobuf/encoding/protowire"
)

const constBytesFieldFlag = uint64(1 << 63)

// evaluate instructions.
func (m *Machine) evaluate() bool {
	insn := m.insn()

	for {
		if debugging {
			debugf("eval: %5d ", len(m.insn())-len(insn))
		}

		opcode := op.Code(insn[0])
		insn = insn[1:]

		switch {
		case opcode < 64: // No arguments.
			switch cmp := opcode.Cmp(); {
			case opcode <= op.LoadConstScalar1:
				if cmp.IsValid() {
					m.opCompareUnsigned(cmp)
				} else {
					m.opLoadConstScalarBool(opcode.Option())
				}

			case opcode <= op.ReturnTrue:
				if cmp.IsValid() {
					m.opCompareSigned(cmp)
				} else {
					return m.opReturn(opcode.Option())
				}

			case opcode < op.CompareFloatLT:
				m.opCompareBytes(cmp)

			case opcode <= op.CompareFloatInfNeg:
				if cmp.IsValid() {
					m.opCompareFloat(cmp)
				} else {
					m.opCompareFloatInf(opcode.Option())
				}

			default:
				if opcode == op.CompareFloatNaN {
					m.opCompareFloatNaN()
				} else {
					m.opContains(opcode)
				}
			}

		case opcode < 128: // 8-bit argument.
			arg := insn[0]
			insn = insn[1:]

			switch {
			case opcode == op.CheckField:
				m.opCheckField(arg)

			default:
				m.opLoadField(arg, opcode.Reg())
			}

		case opcode < 192: // 16-bit argument.
			arg := binary.LittleEndian.Uint16(insn)
			insn = insn[2:]

			var offset uint16
			if opcode == op.Skip {
				offset = m.opSkip(arg)
			} else {
				offset = m.opSkipIf(arg, opcode.Option())
			}
			insn = insn[offset:]

		default: // 64-bit argument.
			arg := binary.LittleEndian.Uint64(insn)
			insn = insn[8:]

			switch {
			case opcode == op.LoadConstScalar:
				m.opLoadConstScalar(arg)

			default:
				m.opLoadConstBytes(arg)
			}
		}
	}
}

func (m *Machine) opCheckField(index uint8) {
	slot := index >> 6
	bit := index & 63
	m.status = m.fieldmask[slot]&(1<<bit) != 0

	if debugging {
		debugf("Status := CheckField[#%d] = %t\n", index, m.status)
	}
}

func (m *Machine) opCompareBytes(cmp op.Cmp) {
	r1 := m.getBytes(m.reg[1])
	r0 := m.getBytes(m.reg[0])
	diff := bytes.Compare(r1, r0)

	switch cmp {
	case op.CmpLT:
		m.status = diff < 0
	case op.CmpGE:
		m.status = diff >= 0
	case op.CmpEQ:
		m.status = diff == 0
	case op.CmpNE:
		m.status = diff != 0
	case op.CmpLE:
		m.status = diff <= 0
	default: // GT
		m.status = diff > 0
	}

	if debugging {
		debugf("Status := CompareBytes %q %s %q = %t\n", r1, cmp, r0, m.status)
	}
}

func (m *Machine) opCompareFloat(cmp op.Cmp) {
	r1 := math.Float64frombits(m.reg[1])
	r0 := math.Float64frombits(m.reg[0])

	switch cmp {
	case op.CmpLT:
		m.status = r1 < r0
	case op.CmpGE:
		m.status = r1 >= r0
	case op.CmpEQ:
		m.status = r1 == r0
	case op.CmpNE:
		m.status = r1 != r0
	case op.CmpLE:
		m.status = r1 <= r0
	default: // GT
		m.status = r1 > r0
	}

	if debugging {
		debugf("Status := CompareFloat %f %s %f = %t\n", r1, cmp, r0, m.status)
	}
}

func (m *Machine) opCompareFloatInf(negative bool) {
	sign := 1
	if negative {
		sign = -1
	}
	f := math.Float64frombits(m.reg[0])
	m.status = math.IsInf(f, sign)

	if debugging {
		debugf("Status := CompareFloatInf[%d] %f = %t\n", sign, f, m.status)
	}
}

func (m *Machine) opCompareFloatNaN() {
	f := math.Float64frombits(m.reg[0])
	m.status = math.IsNaN(f)

	if debugging {
		debugf("Status := CompareFloatNaN %f = %t\n", f, m.status)
	}
}

func (m *Machine) opCompareSigned(cmp op.Cmp) {
	r1 := int64(m.reg[1])
	r0 := int64(m.reg[0])

	switch cmp {
	case op.CmpLT:
		m.status = r1 < r0
	case op.CmpGE:
		m.status = r1 >= r0
	case op.CmpEQ:
		m.status = r1 == r0
	case op.CmpNE:
		m.status = r1 != r0
	case op.CmpLE:
		m.status = r1 <= r0
	default: // GT
		m.status = r1 > r0
	}

	if debugging {
		debugf("Status := CompareSigned %d %s %d = %t\n", r1, cmp, r0, m.status)
	}
}

func (m *Machine) opCompareUnsigned(cmp op.Cmp) {
	r1 := m.reg[1]
	r0 := m.reg[0]

	switch cmp {
	case op.CmpLT:
		m.status = r1 < r0
	case op.CmpGE:
		m.status = r1 >= r0
	case op.CmpEQ:
		m.status = r1 == r0
	case op.CmpNE:
		m.status = r1 != r0
	case op.CmpLE:
		m.status = r1 <= r0
	default: // GT
		m.status = r1 > r0
	}

	if debugging {
		debugf("Status := CompareUnsigned %d %s %d = %t\n", r1, cmp, r0, m.status)
	}
}

func (m *Machine) opContains(opcode op.Code) {
	r1 := m.getBytes(m.reg[1])

	switch opcode {
	case op.ContainsVarint:
		r0 := m.reg[0]
		m.status = containsVarint(r1, r0)

		if debugging {
			debugf("Status := ContainsVarint %v %d = %t\n", r1, r0, m.status)
		}

	case op.ContainsZigZag:
		r0 := int64(m.reg[0])
		m.status = containsZigZag(r1, r0)

		if debugging {
			debugf("Status := ContainsZigZag %v %d = %t\n", r1, r0, m.status)
		}

	case op.ContainsFixed64:
		r0 := m.reg[0]
		m.status = containsFixed64(r1, r0)

		if debugging {
			debugf("Status := ContainsFixed64 %v %d = %t\n", r1, r0, m.status)
		}

	default:
		r0 := uint32(m.reg[0])
		m.status = containsFixed32(r1, r0)

		if debugging {
			debugf("Status := ContainsFixed32 %v %d = %t\n", r1, r0, m.status)
		}
	}
}

func (m *Machine) opLoadConstBytes(arg uint64) {
	m.reg[0] = arg | constBytesFieldFlag

	if debugging {
		debugf("R0     := LoadConstBytes[%#016x] = %#016x\n", arg, m.reg[0])
	}
}

func (m *Machine) opLoadConstScalar(value uint64) {
	m.reg[0] = value

	if debugging {
		debugf("R0     := LoadConstScalar[%d] = %#x\n", value, m.reg[0])
	}
}

func (m *Machine) opLoadConstScalarBool(value bool) {
	var n uint64
	if value {
		n = 1
	}
	m.reg[0] = n

	if debugging {
		debugf("R0     := LoadConstScalar[%t] = %#x\n", value, m.reg[0])
	}
}

func (m *Machine) opLoadField(index uint8, r op.Reg) {
	m.reg[r] = m.fielddata[index]

	if debugging {
		debugf("%s     := LoadField[#%d] = %#x\n", r, index, m.reg[r])
	}
}

func (m *Machine) opReturn(status bool) bool {
	if debugging {
		debugf("          Return[%t]\n", status)
	}

	return status
}

func (m *Machine) opSkip(offset uint16) uint16 {
	if debugging {
		debugf("          Skip %d = %d\n", offset, offset)
	}

	return offset
}

func (m *Machine) opSkipIf(offset uint16, status bool) uint16 {
	var result uint16
	if m.status == status {
		result = offset
	}

	if debugging {
		debugf("          Skip[%t] %d = %d\n", status, offset, result)
	}

	return result
}

func (m *Machine) getBytes(ref uint64) []byte {
	proto := ref&constBytesFieldFlag == 0
	ref &^= constBytesFieldFlag

	off, n := unpackBytesRef(ref)
	if proto {
		return m.protobuf[off:][:n]
	}
	return m.bytecode[off:][:n]
}

func containsFixed32(b []byte, needle uint32) bool {
	for len(b) >= 4 {
		if binary.LittleEndian.Uint32(b) == needle {
			return true
		}
		b = b[4:]
	}
	return false
}

func containsFixed64(b []byte, needle uint64) bool {
	for len(b) >= 8 {
		if binary.LittleEndian.Uint64(b) == needle {
			return true
		}
		b = b[8:]
	}
	return false
}

func containsVarint(b []byte, needle uint64) bool {
	for len(b) > 0 {
		x, n := protowire.ConsumeVarint(b)
		if n < 0 {
			return false
		}
		if x == needle {
			return true
		}
		b = b[n:]
	}
	return false
}

func containsZigZag(b []byte, needle int64) bool {
	for len(b) > 0 {
		x, n := protowire.ConsumeVarint(b)
		if n < 0 {
			return false
		}
		if protowire.DecodeZigZag(x) == needle {
			return true
		}
		b = b[n:]
	}
	return false
}
