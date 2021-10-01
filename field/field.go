// Package field contains types and constants useful with PBF bytecode's field
// section.
package field

import (
	"fmt"
)

// Mod informs the decoding of a protobuf field.
type Mod byte

// Field decoding modifiers.
const (
	// Leaf nodes:

	_ = Mod(iota)
	ModZigZag
	ModFloat

	// Intermediary nodes:

	ModPacked
	ModMessage
	ModRepeated
)

func (m Mod) String() string {
	switch m {
	case 0:
		return "0"
	case ModZigZag:
		return "ZigZag"
	case ModFloat:
		return "Float"
	case ModPacked:
		return "Packed"
	case ModMessage:
		return "Message"
	case ModRepeated:
		return "Repeated"
	default:
		return fmt.Sprintf("<invalid field.Mod value %d>", m)
	}
}

// IsValid value?
func (m Mod) IsValid() bool {
	return m <= ModRepeated
}

// IsLeaf node?  A non-leaf node is used as an intermediary for reaching a leaf
// node.
func (m Mod) IsLeaf() bool {
	return m <= ModFloat
}
