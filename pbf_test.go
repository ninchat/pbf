package pbf_test

import (
	"io/ioutil"
	"testing"

	"github.com/ninchat/pbf"
	"github.com/ninchat/pbf/field"
	"github.com/ninchat/pbf/op"
	"google.golang.org/protobuf/encoding/protowire"
)

//go:generate protoc --go_out=internal/test --go_opt=paths=source_relative test.proto
//go:generate go run internal/test/cmd/testdata.go

var bytecode = []byte{
	'P', 'B', 'F', 0, // Bytecode header.

	// Beginning of field section.

	22, // Field count.

	// Scalar field specs:

	1, 0, 0, 0, 0,
	2, 0, 0, 0, 0,
	3, 0, 0, 0, 0,
	4, 0, 0, 0, 0,
	5, 0, 0, 0, byte(field.ModZigZag),
	6, 0, 0, 0, byte(field.ModZigZag),
	7, 0, 0, 0, 0,
	8, 0, 0, 0, 0,
	9, 0, 0, 0, byte(field.ModFloat),
	10, 0, 0, 0, byte(field.ModFloat),
	11, 0, 0, 0, 0,
	12, 0, 0, 0, 0,
	13, 0, 0, 0, byte(field.ModPacked), byte(protowire.VarintType), 3, 0, 0, 0, 0,
	14, 0, 0, 0, byte(field.ModMessage), 2, 0, 0, 0, 0,
	15, 0, 0, 0, byte(field.ModMessage), 1, 0, 0, 0, byte(field.ModPacked), byte(protowire.VarintType), 0, 0, 0, 0, byte(field.ModZigZag),
	16, 0, 0, 0, byte(field.ModRepeated), 1, 0, 0, 0, byte(field.ModZigZag),
	17, 0, 0, 0, byte(field.ModRepeated), 2, 0, 0, 0, byte(field.ModMessage), 1, 0, 0, 0, 0,
	// Don't specify tag 18.
	19, 0, 0, 0, 0,

	// Vector field specs:

	13, 0, 0, 0, 0,
	15, 0, 0, 0, byte(field.ModMessage), 1, 0, 0, 0, 0,
	20, 0, 0, 0, 0,
	21, 0, 0, 0, 0,

	// End of field section.
	// Beginning of instruction-and-constant section.

	// Constant data:

	byte(op.Skip), 13, 0, // Skip the string constant.
	'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!',

	// Scalar field tests:

	byte(op.LoadR1FieldScalar), 0, // Tag 1.
	byte(op.LoadConstScalar1),
	byte(op.CompareUnsignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 1, // Tag 2.
	byte(op.LoadConstScalar1),
	byte(op.CompareUnsignedNE),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 2, // Tag 3.
	byte(op.LoadConstScalar), 3, 0, 0, 0, 0, 0, 0, 0,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 2, 0,
	byte(op.ReturnFalse),
	0xff, // Unknown opcode.

	byte(op.LoadR1FieldScalar), 3, // Tag 4.
	byte(op.LoadConstScalar0),
	byte(op.CompareSignedLT),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR0FieldScalar), 3,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 4, // Tag 5.
	byte(op.LoadConstScalar), 0xfb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 5, // Tag 6.
	// Keep previous R0 value.
	byte(op.CompareSignedLT),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 6, // Tag 7.
	byte(op.LoadConstScalar), 0x07, 0x00, 0x00, 0x7f, 0, 0, 0, 0,
	byte(op.CompareUnsignedGE),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 7, // Tag 8.
	byte(op.LoadConstScalar), 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80,
	byte(op.CompareUnsignedLE),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 8, // Tag 9.
	byte(op.LoadConstScalar), 0x00, 0x00, 0x00, 0x60, 0xfb, 0x21, 0x09, 0x40, // float32 cast to float64.
	byte(op.CompareFloatEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 9, // Tag 10.
	byte(op.LoadConstScalar), 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40,
	byte(op.CompareFloatEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldBytes), 10, // Tag 11.
	byte(op.LoadConstBytes), 0, 0, 0, 0, 3, 0, 0, 0, // "PBF"
	byte(op.CompareBytesEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldBytes), 11, // Tag 12.
	byte(op.LoadConstBytes), 160, 0, 0, 0, 13, 0, 0, 0, // "Hello, world!"
	byte(op.CompareBytesEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR0FieldBytes), 11,
	byte(op.LoadR1FieldBytes), 11,
	byte(op.CompareBytesEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 12, // Tag 13.
	byte(op.LoadConstScalar), 4, 0, 0, 0, 0, 0, 0, 0,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 13, // Tag 14.
	byte(op.LoadConstScalar), 0xd5, 0xdd, 0x00, 0x00, 0, 0, 0, 0,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 14, // Tag 15.
	byte(op.LoadConstScalar), 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 15, // Tag 16.
	byte(op.LoadConstScalar), 20, 0, 0, 0, 0, 0, 0, 0,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.CheckField), 16, // Tag 17.
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldScalar), 16,
	byte(op.LoadConstScalar), 102, 0, 0, 0, 0, 0, 0, 0,
	byte(op.CompareSignedEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	// (Tag 18 not used.)

	byte(op.CheckField), 17, // Tag 19.
	byte(op.SkipFalse), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldBytes), 17,
	byte(op.LoadConstBytes), 0, 0, 0, 0, 0, 0, 0, 0, // Empty string.
	byte(op.CompareBytesEQ),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	// Vector field tests:

	byte(op.LoadR1FieldVector), 18, // Tag 13.
	byte(op.LoadConstScalar), 4, 0, 0, 0, 0, 0, 0, 0,
	byte(op.ContainsVarint),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldVector), 19, // Tag 15.
	byte(op.LoadConstScalar1),
	byte(op.ContainsZigZag),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldVector), 20, // Tag 20.
	byte(op.LoadConstScalar), 0xdb, 0x0f, 0x49, 0x40, 0, 0, 0, 0, // float32 as is.
	byte(op.ContainsFixed32),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.LoadR1FieldVector), 21, // Tag 21.
	byte(op.LoadConstScalar), 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40,
	byte(op.ContainsFixed64),
	byte(op.SkipTrue), 1, 0,
	byte(op.ReturnFalse),

	byte(op.ReturnTrue),

	// End of instruction-and-constant section.
}

func getTestData() []byte {
	buf, err := ioutil.ReadFile("internal/testdata/test.buf")
	if err != nil {
		panic(err)
	}
	return buf
}

func TestPBF(t *testing.T) {
	prog, err := pbf.NewProgram(bytecode)
	if err != nil {
		t.Fatal(err)
	}

	mach := pbf.NewMachine(prog)
	buf := getTestData()

	ok, err := mach.Filter(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error(ok)
	}

	for i := 0; i < 256; i++ {
		if _, found := mach.GetRawValue(uint8(i)); found != (i < 22 && i != 17) {
			t.Error(i, found)
		}
	}

	if v, found := mach.GetRawValue(7); !found || v != 0x8000000000000008 {
		t.Error(v, found)
	}
}

func BenchmarkPrepare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := pbf.NewProgram(bytecode); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFilter(b *testing.B) {
	prog, err := pbf.NewProgram(bytecode)
	if err != nil {
		b.Fatal(err)
	}

	mach := pbf.NewMachine(prog)
	buf := getTestData()

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ok, err := mach.Filter(buf)
		if err != nil {
			b.Fatal(err)
		}
		if !ok {
			b.Fatal(ok)
		}
	}
}
