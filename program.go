package pbf

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const bytecodeHeader = uint32(0x00464250) // "PBF\0"

var (
	errBytecodeFormat  = errors.New("pbf: unknown bytecode format")
	errBytecodeInvalid = errors.New("pbf: bytecode is invalid")
	errBytecodeTooLong = errors.New("pbf: bytecode is too long")
)

// Program for filtering protobuf messages.
type Program struct {
	program
}

// NewProgram decodes and verifies a PBF bytecode program.
func NewProgram(bytecode []byte) (*Program, error) {
	if len(bytecode) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	header := binary.LittleEndian.Uint32(bytecode)
	off := 4

	if header != bytecodeHeader {
		return nil, errBytecodeFormat
	}

	if len(bytecode) > math.MaxInt32 {
		// Const offsets and lengths could overflow the field data encoding.
		return nil, errBytecodeTooLong
	}

	fieldspec, fieldcount, n, err := parseFieldSection(bytecode[off:])
	if err != nil {
		return nil, err
	}
	off += n

	p := program{
		bytecode:   bytecode,
		fieldspec:  fieldspec,
		fieldcount: fieldcount,
		insnoffset: off,
	}

	if debugging {
		debugf("prog:   Instruction offset: %d\n", p.insnoffset)
	}

	if err := verify(p); err != nil {
		return nil, err
	}

	return &Program{p}, nil
}

type program struct {
	bytecode   []byte
	fieldspec  map[int32]fieldSpec
	fieldcount uint8
	insnoffset int
}

func (p *program) insn() []byte {
	return p.bytecode[p.insnoffset:]
}
