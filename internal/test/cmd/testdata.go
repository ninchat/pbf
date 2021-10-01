package main

import (
	"io/ioutil"
	"math"

	"github.com/ninchat/pbf/internal/test"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func main() {
	m := &test.Test{
		A: 1,
		B: 2,
		C: 3,
		D: -4,
		E: -5,
		F: -6,
		G: 0x7f000007,
		H: 0x8000000000000008,
		I: math.Pi,
		J: math.Pi,
		K: []byte("PBF"),
		L: "Hello, world!",
		M: []uint64{1, 2, 3, 4, 5},
		N: &test.Sub{X: 1234, Y: 56789},
		O: &test.List{
			Z: []int32{-3, -2, -1, 0, 1, 2},
		},
		P: []int32{10, 20, 30, 40},
		Q: []*test.Sub{
			&test.Sub{X: 100, Y: 200},
			&test.Sub{X: 101, Y: 201},
			&test.Sub{X: 102, Y: 202},
		},
		R: 18,
		// Don't set S.
		T: []float32{1.1, 2.2, math.Pi, 4.4, 5.5},
		U: []float64{1.1, 2.2, math.Pi, 4.4, 5.5},
	}

	buf, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("internal/testdata/test.buf", buf, 0o644); err != nil {
		panic(err)
	}

	print(prototext.Format(m))
}
