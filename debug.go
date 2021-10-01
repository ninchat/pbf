// +build pbfdebug

package pbf

import (
	"fmt"
	"os"

	"github.com/ninchat/pbf/field"
	"google.golang.org/protobuf/encoding/protowire"
)

const debugging = true

func debugf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func (f fieldSpec) String() string {
	s := fmt.Sprintf("#%d", f.index)

	if f.mod != 0 {
		s += fmt.Sprintf(" %s", f.mod)

		if f.mod == field.ModPacked {
			s += fmt.Sprintf(" %v", protowire.Type(f.subtype))
		}
	}

	return s
}
