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
		fieldcount: fieldcount,
		insnoffset: off,
	}
	p.initFieldSpec(fieldspec)

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
	fieldcount uint8
	insnoffset int

	fieldspecarr *[256]fieldSpec
	maxarrindex  uint8

	fieldspecmap map[int32]fieldSpec
}

func (p *program) initFieldSpec(spec map[int32]fieldSpec) {
	var maxtag uint8
	for tag := range spec {
		if tag < 0 || tag >= 256 {
			p.fieldspecmap = spec
			return
		}
		if uint8(tag) > maxtag {
			maxtag = uint8(tag)
		}
	}

	var arr [256]fieldSpec
	for i := range arr {
		arr[i].mod = 0xff // Invalid mod value.
	}
	for tag, s := range spec {
		arr[uint8(tag)] = s
	}

	p.fieldspecarr = &arr
	p.maxarrindex = maxtag
}

func (p *program) insn() []byte {
	return p.bytecode[p.insnoffset:]
}
