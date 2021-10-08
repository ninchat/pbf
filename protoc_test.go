package pbf_test

import (
	"testing"

	"github.com/ninchat/pbf/internal/test"
)

func BenchmarkProtoc(b *testing.B) {
	state := test.ProtocTester{}
	buf := getTestData()

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ok, err := state.Filter(buf)
		if err != nil {
			b.Fatal(err)
		}
		if !ok {
			b.Fatal(ok)
		}
	}
}
