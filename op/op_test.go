package op

import (
	"testing"
)

func TestOptions(t *testing.T) {
	if LoadConstScalar0.Option() {
		t.Error(1)
	}
	if !LoadConstScalar1.Option() {
		t.Error(1)
	}
	if ReturnFalse.Option() {
		t.Error(1)
	}
	if !ReturnTrue.Option() {
		t.Error(1)
	}
	if CompareFloatInfPos.Option() {
		t.Error(1)
	}
	if !CompareFloatInfNeg.Option() {
		t.Error(1)
	}
	if SkipFalse.Option() {
		t.Error(1)
	}
	if !SkipTrue.Option() {
		t.Error(1)
	}
}

func TestRegs(t *testing.T) {
	if LoadR0FieldScalar.Reg() != R0 || LoadR0FieldScalar.Reg() != 0 {
		t.Error(1)
	}
	if LoadR1FieldScalar.Reg() != R1 || LoadR1FieldScalar.Reg() != 1 {
		t.Error(1)
	}
	if LoadR0FieldBytes.Reg() != R0 || LoadR0FieldBytes.Reg() != 0 {
		t.Error(1)
	}
	if LoadR1FieldBytes.Reg() != R1 || LoadR1FieldBytes.Reg() != 1 {
		t.Error(1)
	}
	if LoadR0FieldVector.Reg() != R0 || LoadR0FieldVector.Reg() != 0 {
		t.Error(1)
	}
	if LoadR1FieldVector.Reg() != R1 || LoadR1FieldVector.Reg() != 1 {
		t.Error(1)
	}
}

func TestCmps(t *testing.T) {
	if CompareUnsignedLT.Cmp() != CmpLT || !CompareUnsignedLT.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareUnsignedGE.Cmp() != CmpGE || !CompareUnsignedGE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareUnsignedEQ.Cmp() != CmpEQ || !CompareUnsignedEQ.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareUnsignedNE.Cmp() != CmpNE || !CompareUnsignedNE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareUnsignedLE.Cmp() != CmpLE || !CompareUnsignedLE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareUnsignedGT.Cmp() != CmpGT || !CompareUnsignedGT.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareSignedLT.Cmp() != CmpLT || !CompareSignedLT.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareSignedGE.Cmp() != CmpGE || !CompareSignedGE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareSignedEQ.Cmp() != CmpEQ || !CompareSignedEQ.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareSignedNE.Cmp() != CmpNE || !CompareSignedNE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareSignedLE.Cmp() != CmpLE || !CompareSignedLE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareSignedGT.Cmp() != CmpGT || !CompareSignedGT.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareBytesLT.Cmp() != CmpLT || !CompareBytesLT.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareBytesGE.Cmp() != CmpGE || !CompareBytesGE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareBytesEQ.Cmp() != CmpEQ || !CompareBytesEQ.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareBytesNE.Cmp() != CmpNE || !CompareBytesNE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareBytesLE.Cmp() != CmpLE || !CompareBytesLE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareBytesGT.Cmp() != CmpGT || !CompareBytesGT.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareFloatLT.Cmp() != CmpLT || !CompareFloatLT.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareFloatGE.Cmp() != CmpGE || !CompareFloatGE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareFloatEQ.Cmp() != CmpEQ || !CompareFloatEQ.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareFloatNE.Cmp() != CmpNE || !CompareFloatNE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareFloatLE.Cmp() != CmpLE || !CompareFloatLE.Cmp().IsValid() {
		t.Error(1)
	}
	if CompareFloatGT.Cmp() != CmpGT || !CompareFloatGT.Cmp().IsValid() {
		t.Error(1)
	}
}
