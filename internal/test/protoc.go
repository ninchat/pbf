package test

import (
	"bytes"
	"math"

	"google.golang.org/protobuf/proto"
)

var options = proto.UnmarshalOptions{
	DiscardUnknown: true,
}

// ProtocTester mimicks the program in pbf_test.go using Go code generated by
// protoc and protoc-gen-go.
type ProtocTester struct {
	m Test
}

func (t *ProtocTester) Filter(b []byte) (bool, error) {
	if err := options.Unmarshal(b, &t.m); err != nil {
		return false, err
	}

	return checkFields(&t.m), nil
}

func checkFields(m *Test) bool {
	if !(m.A == 1) {
		return false
	}

	if !(m.B != 1) {
		return false
	}

	if !(m.C == 3) {
		return false
	}

	if !(m.D < 0) {
		return false
	}
	if !(m.D == m.D) {
		return false
	}

	if !(m.E == -5) {
		return false
	}

	if !(m.F < -5) {
		return false
	}

	if !(m.G >= 0x7f000007) {
		return false
	}

	if !(m.H <= 0x8000000000000008) {
		return false
	}

	if !(m.I == math.Pi) {
		return false
	}

	if !(m.J == math.Pi) {
		return false
	}

	if !bytes.Equal(m.K, []byte("PBF")) {
		return false
	}

	if !(m.L == "Hello, world!") {
		return false
	}
	if !(m.L == m.L) {
		return false
	}

	if !(len(m.M) >= 4 && m.M[3] == 4) {
		return false
	}

	if !(m.N.Y == 56789) {
		return false
	}

	if !(m.O != nil && len(m.O.Z) >= 1 && m.O.Z[0] == -3) {
		return false
	}

	if !(len(m.P) >= 2 && m.P[1] == 20) {
		return false
	}

	if !(len(m.Q) >= 3 && m.Q[2].X == 102) {
		return false
	}

	if !(m.S == nil) {
		return false
	}
	if !(len(m.S) == 0) {
		return false
	}

	if !containsUint64(m.M, 4) {
		return false
	}

	if !containsInt32(m.O.Z, 1) {
		return false
	}

	if !containsFloat32(m.T, math.Pi) {
		return false
	}

	if !containsFloat64(m.U, math.Pi) {
		return false
	}

	return true
}

func containsUint64(haystack []uint64, needle uint64) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

func containsInt32(haystack []int32, needle int32) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

func containsFloat32(haystack []float32, needle float32) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

func containsFloat64(haystack []float64, needle float64) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}
