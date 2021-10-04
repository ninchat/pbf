"See the Go source files for documentation."

from enum import IntEnum
from struct import pack
from typing import List, Optional

BYTECODE_HEADER = b"PBF\0"


class FieldMod(IntEnum):
    Default = 0
    ZigZag = 1
    Float = 2
    Packed = 3
    Message = 4
    Repeated = 5

    @property
    def leaf(self) -> bool:
        return self <= self.Float


class FieldType(IntEnum):
    Varint = 0
    Fixed32 = 5
    Fixed64 = 1
    Bytes = 2


class FieldSpec:

    def __init__(self,
                 num: int,
                 mod: FieldMod = FieldMod.Default,
                 subtype: Optional[FieldType] = None) -> None:
        assert num in range(0, 1 << 31)
        assert (mod == FieldMod.Packed) == (subtype is not None)
        self.num = num
        self.mod = mod
        self.subtype = subtype
        self.parent = None

    def sub(self,
            num: int,
            mod: FieldMod = FieldMod.Default,
            subtype: Optional[FieldType] = None) -> 'FieldSpec':
        assert self.mod in (FieldMod.Packed, FieldMod.Message, FieldMod.Repeated)
        child = FieldSpec(num, mod, subtype)
        child.parent = self
        return child

    def encode(self) -> bytes:
        assert self.mod not in (FieldMod.Packed, FieldMod.Message, FieldMod.Repeated)
        b = b""
        f = self
        while f:
            if f.mod == FieldMod.Packed:
                b = pack("<IBB", f.num, f.mod, f.subtype) + b
            else:
                b = pack("<IB", f.num, f.mod) + b
            f = f.parent
        return b


def encode_field_section_header(field_count: int) -> bytes:
    return pack("<B", field_count)


def encode_field_section(fields: List[FieldSpec]) -> bytes:
    b = encode_field_section_header(len(fields))
    for f in fields:
        b += f.encode()
    return b


class ValueKind(IntEnum):
    Unsigned = 0
    Signed = 8
    Bytes = 16
    Float = 24


class FieldKind(IntEnum):
    Scalar = 0
    Bytes = 2
    Vector = 4


class ScalarEncoding(IntEnum):
    Varint = 0
    ZigZag = 1
    Fixed64 = 2
    Fixed32 = 3


class Reg(IntEnum):
    R0 = 0
    R1 = 1


R0 = Reg.R0
R1 = Reg.R1


class Cmp(IntEnum):
    LT = 0
    GE = 1
    EQ = 2
    NE = 3
    LE = 4
    GT = 5


class Op(IntEnum):
    CompareUnsignedLT = 0
    CompareUnsignedGE = 1
    CompareUnsignedEQ = 2
    CompareUnsignedNE = 3
    CompareUnsignedLE = 4
    CompareUnsignedGT = 5
    LoadConstScalar0 = 6
    LoadConstScalar1 = 7
    CompareSignedLT = 8
    CompareSignedGE = 9
    CompareSignedEQ = 10
    CompareSignedNE = 11
    CompareSignedLE = 12
    CompareSignedGT = 13
    ReturnFalse = 14
    ReturnTrue = 15
    CompareBytesLT = 16
    CompareBytesGE = 17
    CompareBytesEQ = 18
    CompareBytesNE = 19
    CompareBytesLE = 20
    CompareBytesGT = 21
    CompareFloatLT = 24
    CompareFloatGE = 25
    CompareFloatEQ = 26
    CompareFloatNE = 27
    CompareFloatLE = 28
    CompareFloatGT = 29
    CompareFloatInfPos = 30
    CompareFloatInfNeg = 31
    CompareFloatNaN = 32
    ContainsVarint = 36
    ContainsZigZag = 37
    ContainsFixed64 = 38
    ContainsFixed32 = 39
    LoadR0FieldScalar = 64
    LoadR1FieldScalar = 65
    LoadR0FieldBytes = 66
    LoadR1FieldBytes = 67
    LoadR0FieldVector = 68
    LoadR1FieldVector = 69
    CheckField = 70
    SkipFalse = 128
    SkipTrue = 129
    Skip = 130
    LoadConstScalar = 192
    LoadConstBytes = 194

    @classmethod
    def compare_(cls, kind: ValueKind, cmp: Cmp) -> 'Op':
        assert isinstance(kind, ValueKind)
        assert isinstance(cmp, Cmp)
        return cls(cls.CompareUnsignedLT + kind + cmp)

    @classmethod
    def load_const_scalar_(cls, bit: bool) -> 'Op':
        assert bit in (False, True)
        return cls(cls.LoadConstScalar0 + bit)

    @classmethod
    def return_(cls, status: bool) -> 'Op':
        assert status in (False, True)
        return cls(cls.ReturnFalse + status)

    @classmethod
    def compare_float_inf_(cls, negative: bool) -> 'Op':
        assert negative in (False, True)
        return cls(cls.CompareFloatInfPos + negative)

    @classmethod
    def contains_(cls, enc: ScalarEncoding) -> 'Op':
        assert isinstance(enc, ScalarEncoding)
        return cls(cls.ContainsVarint + enc)

    @classmethod
    def load_field_(cls, reg: Reg, kind: FieldKind) -> 'Op':
        assert reg in (R0, R1)
        assert isinstance(kind, FieldKind)
        return cls(cls.LoadR0FieldScalar + reg + kind)

    @classmethod
    def skip_(cls, status: bool) -> 'Op':
        assert status in (False, True)
        return cls(cls.SkipFalse + status)

    @property
    def size(self) -> int:
        "Size of the encoded instruction."
        if self < 64:
            return 1
        elif self < 128:
            return 1 + 1
        elif self < 192:
            return 1 + 2
        else:
            return 1 + 8

    def encode(self, arg: Optional[int] = None) -> bytes:
        if self < 64:
            assert arg is None
            return pack("<B", self)
        elif self < 128:
            assert isinstance(arg, int)
            return pack("<BB", self, arg)
        elif self < 192:
            assert isinstance(arg, int)
            return pack("<BH", self, arg)
        else:
            assert isinstance(arg, int)
            return pack("<BQ", self, arg)


def const_bytes_ref(offset: int, length: int) -> int:
    "Form an argument for the LoadConstBytes op."
    assert isinstance(offset, int)
    assert isinstance(length, int)
    assert offset in range(0, 1 << 31)
    assert length in range(0, 1 << 31)
    return (length << 32) | offset


if __name__ == "__main__":
    assert FieldMod.Float.leaf
    assert not FieldMod.Packed.leaf

    assert encode_field_section([
        FieldSpec(1),
        FieldSpec(2, FieldMod.Packed, FieldType.Varint).sub(9, FieldMod.Float),
        FieldSpec(4, FieldMod.Message).sub(2, FieldMod.Message).sub(65536, FieldMod.ZigZag),
    ]) == bytes([
        3,                 # Field count.

        1, 0, 0, 0, 0,     # Field at index 0.

        2, 0, 0, 0, 3, 0,
        9, 0, 0, 0, 2,     # Field at index 1.

        4, 0, 0, 0, 4,
        2, 0, 0, 0, 4,
        0, 0, 1, 0, 1,     # Field at index 2.
    ])

    assert Op.LoadConstScalar.size == 9
    assert Op.LoadConstScalar.encode(255) == b"\xc0\xff\x00\x00\x00\x00\x00\x00\x00"

    assert Op.load_field_(1, FieldKind.Bytes).size == 2
    assert Op.load_field_(1, FieldKind.Bytes).encode(42) == b"\x43\x2a"
