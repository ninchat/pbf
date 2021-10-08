package pbf_test

import (
	"encoding/binary"
	"os"
	"os/exec"
	"testing"

	"github.com/ninchat/pbf/internal/test"
)

func BenchmarkProtocGo(b *testing.B) {
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

//go:generate protoc --cpp_out=internal/test/cxx test.proto

var compiledCXX bool

func compileCXX() error {
	if compiledCXX {
		return nil
	}

	cxx := os.Getenv("CXX")
	if cxx == "" {
		cxx = "c++"
	}

	cmd := exec.Command(cxx, "-O2", "-Wall", "-o", "internal/testdata/test-cxx", "internal/test/cxx/protoc.cc", "internal/test/cxx/test.pb.cc", "-lprotobuf-lite")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	compiledCXX = true
	return nil
}

func BenchmarkProtocCXX(b *testing.B) {
	if err := compileCXX(); err != nil {
		b.Fatal(err)
	}

	cmd := exec.Command("internal/testdata/test-cxx")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		b.Fatal(err)
	}
	defer func() {
		if cmd != nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	}()

	buf := getTestData()
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(buf)))

	if _, err := stdin.Write(size); err != nil {
		b.Fatal(err)
	}
	if _, err := stdin.Write(buf); err != nil {
		b.Fatal(err)
	}

	count := make([]byte, 8)
	binary.LittleEndian.PutUint64(count, uint64(b.N))

	result := make([]byte, 1)

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()

	if _, err := stdin.Write(count); err != nil {
		b.Fatal(err)
	}

	if _, err := stdout.Read(result); err != nil {
		b.Fatal(err)
	}

	b.StopTimer()

	if result[0] != 1 {
		b.Error(result[0])
	}

	if err := cmd.Wait(); err != nil {
		b.Error(err)
	}
	cmd = nil
}
