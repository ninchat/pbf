package pbf

import (
	"errors"
	"io"
	"math"

	"github.com/ninchat/pbf/field"
	"google.golang.org/protobuf/encoding/protowire"
)

var (
	errProtobufDeprecated = errors.New("pbf: protobuf message uses deprecated encoding")
	errProtobufFieldType  = errors.New("pbf: protobuf message has unexpected field type")
	errProtobufInvalid    = errors.New("pbf: protobuf message encoding is invalid")
	errProtobufTooLong    = errors.New("pbf: protobuf message is too long")
)

// decode protobuf fields.
func (m *Machine) decode() (err error) {
	if len(m.protobuf) > math.MaxInt32 {
		// Byte offsets and lengths could overflow the field data encoding.
		return errProtobufTooLong
	}

	if m.fieldspecarr != nil {
		err = m.decodeTopMessage(m.protobuf)
	} else {
		err = m.decodeMessage(m.fieldspecmap, 0, m.protobuf)
	}

	if debugging && err != nil {
		debugf(" ... Error: %v\n", err)
	}

	return
}

// decodeTopMessage where all referenced tag numbers are smaller than 256.
func (m *Machine) decodeTopMessage(buf []byte) error {
	if debugging {
		debugf("decode: Message{\ndecode: ")
	}

	for off := 0; off < len(buf); {
		tag, typ, n := protowire.ConsumeTag(buf[off:])
		if n < 0 {
			return protowire.ParseError(n)
		}
		off += n

		switch typ {
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := m.getTopMessageFieldSpec(tag); found {
				if err := m.decodeFieldScalar(&s, v); err != nil {
					return err
				}
			}

		case protowire.Fixed32Type:
			v, n := protowire.ConsumeFixed32(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := m.getTopMessageFieldSpec(tag); found {
				if err := m.decodeFieldScalar32(&s, v); err != nil {
					return err
				}
			}

		case protowire.Fixed64Type:
			v, n := protowire.ConsumeFixed64(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := m.getTopMessageFieldSpec(tag); found {
				if err := m.decodeFieldScalar64(&s, v); err != nil {
					return err
				}
			}

		case protowire.BytesType:
			b, taglen, err := consumeProtoBytes(buf[off:])
			if err != nil {
				return err
			}
			off += taglen

			if s, found := m.getTopMessageFieldSpec(tag); found {
				if err := m.decodeFieldBytes(&s, off, b); err != nil {
					return err
				}
			}
			off += len(b)

		case protowire.StartGroupType, protowire.EndGroupType:
			return errProtobufDeprecated

		default:
			return errProtobufInvalid
		}

		if debugging {
			debugf("\ndecode: ")
		}
	}

	if debugging {
		debugf("}\n")
	}

	return nil
}

func (m *Machine) decodeMessage(spec map[int32]fieldSpec, base int, buf []byte) error {
	if debugging {
		if base == 0 {
			debugf("decode: Message{\ndecode: ")
		} else {
			debugf("=Message{")
		}
	}

	rep := m.fieldrep.get()
	defer m.fieldrep.put(rep)

	for off := 0; off < len(buf); {
		tag, typ, n := protowire.ConsumeTag(buf[off:])
		if n < 0 {
			return protowire.ParseError(n)
		}
		off += n

		switch typ {
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := getMessageFieldSpec(spec, rep, int32(tag)); found {
				if err := m.decodeFieldScalar(&s, v); err != nil {
					return err
				}
			}

		case protowire.Fixed32Type:
			v, n := protowire.ConsumeFixed32(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := getMessageFieldSpec(spec, rep, int32(tag)); found {
				if err := m.decodeFieldScalar32(&s, v); err != nil {
					return err
				}
			}

		case protowire.Fixed64Type:
			v, n := protowire.ConsumeFixed64(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := getMessageFieldSpec(spec, rep, int32(tag)); found {
				if err := m.decodeFieldScalar64(&s, v); err != nil {
					return err
				}
			}

		case protowire.BytesType:
			b, taglen, err := consumeProtoBytes(buf[off:])
			if err != nil {
				return err
			}
			off += taglen

			if s, found := getMessageFieldSpec(spec, rep, int32(tag)); found {
				if err := m.decodeFieldBytes(&s, base+off, b); err != nil {
					return err
				}
			}
			off += len(b)

		case protowire.StartGroupType, protowire.EndGroupType:
			return errProtobufDeprecated

		default:
			return errProtobufInvalid
		}

		if debugging && base == 0 {
			debugf("\ndecode: ")
		}
	}

	if debugging {
		if base == 0 {
			debugf("}\n")
		} else {
			debugf(" }")
		}
	}

	return nil
}

func (m *Machine) decodePacked(typ uint8, spec map[int32]fieldSpec, base int, buf []byte) error {
	if debugging {
		debugf("=Packed{")
	}

	switch protowire.Type(typ) {
	case protowire.VarintType:
		for off, i := 0, int32(0); off < len(buf); i++ {
			v, n := protowire.ConsumeVarint(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := getFieldSpec(spec, i); found {
				if err := m.decodeFieldScalar(&s, v); err != nil {
					return err
				}
			}
		}

	case protowire.Fixed32Type:
		for off, i := 0, int32(0); off < len(buf); i++ {
			v, n := protowire.ConsumeFixed32(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := getFieldSpec(spec, i); found {
				if err := m.decodeFieldScalar32(&s, v); err != nil {
					return err
				}
			}
		}

	case protowire.Fixed64Type:
		for off, i := 0, int32(0); off < len(buf); i++ {
			v, n := protowire.ConsumeFixed64(buf[off:])
			if n < 0 {
				return protowire.ParseError(n)
			}
			off += n

			if s, found := getFieldSpec(spec, i); found {
				if err := m.decodeFieldScalar64(&s, v); err != nil {
					return err
				}
			}
		}

	default:
		for off, i := 0, int32(0); off < len(buf); i++ {
			b, taglen, err := consumeProtoBytes(buf[off:])
			if err != nil {
				return err
			}
			off += taglen

			if s, found := getFieldSpec(spec, i); found {
				if err := m.decodeFieldBytes(&s, base+off, b); err != nil {
					return err
				}
			}
			off += len(b)
		}
	}

	if debugging {
		debugf(" }")
	}

	return nil
}

func (m *Machine) decodeFieldBytes(s *fieldSpec, off int, buf []byte) error {
	if s.indexed {
		m.setFieldBytes(s, off, buf)
	}
	if s.mod.IsLeaf() {
		return nil
	}

	if s.mod == field.ModPacked {
		return m.decodePacked(s.subtype, s.sub, off, buf)
	}
	return m.decodeMessage(s.sub, off, buf)
}

func (m *Machine) decodeFieldScalar(s *fieldSpec, value uint64) error {
	if s.mod == field.ModZigZag {
		value = uint64(protowire.DecodeZigZag(value))
	}

	return m.setFieldScalar(s, value)
}

func (m *Machine) decodeFieldScalar32(s *fieldSpec, v uint32) error {
	var value uint64
	if s.mod == field.ModFloat {
		value = math.Float64bits(float64(math.Float32frombits(v)))
	} else {
		value = uint64(v)
	}

	return m.setFieldScalar(s, value)
}

func (m *Machine) decodeFieldScalar64(s *fieldSpec, value uint64) error {
	return m.setFieldScalar(s, value)
}

func (m *Machine) setField(s *fieldSpec, data uint64) {
	if debugging {
		debugf("(%#x)=#%d", data, s.index)
	}

	m.fielddata[s.index] = data
	m.fieldmask[s.maskslot()] |= 1 << s.maskbit()
}

func (m *Machine) setFieldBytes(s *fieldSpec, off int, buf []byte) {
	if debugging {
		debugf("=Bytes%q", buf)
	}

	m.setField(s, packBytesRef(off, len(buf)))
}

func (m *Machine) setFieldScalar(s *fieldSpec, value uint64) error {
	if !s.indexed {
		return errProtobufFieldType
	}

	if debugging {
		debugf("=Scalar")
	}

	m.setField(s, value)
	return nil
}

func (m *Machine) getTopMessageFieldSpec(tag protowire.Number) (s fieldSpec, found bool) {
	if debugging {
		debugf(" .%d", tag)
	}

	if tag > protowire.Number(m.maxarrindex) {
		return
	}

	s = m.fieldspecarr[tag]
	found = s.mod.IsValid()
	if !found {
		return
	}

	if s.mod == field.ModRepeated {
		index := m.topfieldrep[tag]
		m.topfieldrep[tag] = index + 1

		s, found = s.sub[index]
	}
	return
}

func getFieldSpec(spec map[int32]fieldSpec, num int32) (s fieldSpec, found bool) {
	if debugging {
		debugf(" .%d", num)
	}

	s, found = spec[num]
	return
}

func getMessageFieldSpec(spec map[int32]fieldSpec, rep map[int32]int32, num int32) (s fieldSpec, found bool) {
	s, found = getFieldSpec(spec, num)
	if !found {
		return
	}

	if s.mod == field.ModRepeated {
		index := rep[num]
		rep[num] = index + 1

		s, found = s.sub[index]
	}
	return
}

// consumeProtoBytes returns tag length (payload offset), not consumed length.
func consumeProtoBytes(buf []byte) (payload []byte, off int, err error) {
	size, off := protowire.ConsumeVarint(buf)
	if off < 0 {
		err = protowire.ParseError(off)
		return
	}

	if size > uint64(len(buf[off:])) {
		err = io.ErrUnexpectedEOF
		return
	}

	payload = buf[off:][:size]
	return
}

func packBytesRef(offset, length int) uint64 {
	return uint64(offset) | (uint64(length) << 32)
}

func unpackBytesRef(ref uint64) (offset, length uint32) {
	offset = uint32(ref)
	length = uint32(ref >> 32)
	return
}
