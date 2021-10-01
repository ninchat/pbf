package pbf

import (
	"encoding/binary"
	"io"

	"github.com/ninchat/pbf/field"
	"google.golang.org/protobuf/encoding/protowire"
)

type fieldSpec struct {
	indexed bool
	index   uint8
	mod     field.Mod
	subtype uint8 // Meaningful only if ModPacked.
	sub     map[int32]fieldSpec
}

func (f *fieldSpec) maskslot() uint8 { return f.index >> 6 }
func (f *fieldSpec) maskbit() uint8  { return f.index & 63 }

func parseFieldSection(buf []byte) (map[int32]fieldSpec, uint8, int, error) {
	if len(buf) == 0 {
		return nil, 0, 0, errBytecodeInvalid
	}
	count := buf[0]
	size := 1

	dest := make(map[int32]fieldSpec, count)

	for i := uint8(0); i < count; i++ {
		if debugging {
			debugf("field:  Proto")
		}

		n, err := parseFieldSpec(dest, buf[size:], i, "")
		size += n
		if err != nil {
			return nil, i, size, err
		}
	}

	return dest, count, size, nil
}

func parseFieldSpec(dest map[int32]fieldSpec, buf []byte, index uint8, anno string) (int, error) {
	if len(buf) < 5 {
		return 0, io.ErrUnexpectedEOF
	}
	key := int32(binary.LittleEndian.Uint32(buf))
	mod := field.Mod(buf[4])
	size := 5
	if !mod.IsValid() {
		return size, errBytecodeInvalid
	}

	if debugging {
		if anno == "" {
			debugf(".%d", key)
		} else {
			debugf(".%s[%d]", anno, key)
		}
	}

	if mod.IsLeaf() {
		s, found := dest[key]
		if found {
			// The node is already used as an intermediary.  It can be
			// referenced directly only as a vector (no mod).
			if mod != 0 || s.indexed {
				return size, errBytecodeInvalid
			}
		} else {
			s.mod = mod
		}
		s.index = index
		s.indexed = true
		dest[key] = s

		if debugging {
			debugf(" = %s\n", dest[key])
		}
	} else {
		var (
			subtype uint8
			subanno string
		)

		switch mod {
		case field.ModPacked:
			if len(buf) == size {
				return size, io.ErrUnexpectedEOF
			}
			subtype = buf[size]
			size++

			switch protowire.Type(subtype) {
			case protowire.VarintType:
				subanno = "Varint"
			case protowire.Fixed32Type:
				subanno = "Fixed32"
			case protowire.Fixed64Type:
				subanno = "Fixed64"
			case protowire.BytesType:
				subanno = "Bytes"
			default:
				return size, errBytecodeInvalid
			}

		case field.ModRepeated:
			subanno = "Repeated"
		}

		s, found := dest[key]
		if found {
			// The node is already be used as an intermediary (same specs), or
			// referenced directly as a vector (no mod).
			if s.mod != 0 && (s.mod != mod || s.subtype != subtype) {
				panic(errBytecodeInvalid) // return size, errBytecodeInvalid
			}
		}
		s.mod = mod
		s.subtype = subtype
		if s.sub == nil {
			s.sub = make(map[int32]fieldSpec)
		}
		dest[key] = s

		n, err := parseFieldSpec(s.sub, buf[size:], index, subanno)
		size += n
		if err != nil {
			return size, err
		}
	}

	return size, nil
}
