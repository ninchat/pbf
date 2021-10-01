/*

Package pbf implements Protocol Buffers Filter.


Bytecode

Bytecode contains a 4-byte header and two sections:

- Fields
- Instructions and constants

Field section:

- Field count (1 byte)
- Sequence of variable-length field specifications

Field specification formats (NUM is protobuf field number or sequence index as
fixed-width 32-bit integer, and SUBTYPE is protobuf field wire type):

- NUM 0
- NUM ModZigZag
- NUM ModFloat
- NUM ModPacked SUBTYPE ...
- NUM ModMessage ...
- NUM ModRepeated ...

Instruction-and-constant section:

- Sequence of variable-length instructions (opcodes followed by arguments)
- Constant byte sequences are interleaved with the instructions

The boundary between the sections is determined by the field count.

The Skip instruction can be used to skip over constants, or the constants may
simply reside beyond the last Return instruction.

All integers use little-endian encoding.

*/
package pbf
