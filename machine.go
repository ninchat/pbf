package pbf

// Machine for program evaluation.  There can be many instances per program,
// but each instance can be used only by a single goroutine at a time.
type Machine struct {
	status      bool
	reg         [2]uint64
	protobuf    []byte    // Encoded protobuf message.
	fielddata   []uint64  // Decoded fields.
	fieldmask   [4]uint64 // Decoded field existence.
	topfieldrep *[256]int32
	fieldrep    repMapPool
	*program
}

// NewMachine creates a machine instance.
func NewMachine(p *Program) *Machine {
	m := &Machine{
		fielddata: make([]uint64, p.fieldcount),
		program:   &p.program,
	}
	if p.fieldspecarr != nil {
		m.topfieldrep = new([256]int32)
	}
	return m
}

// Filter a protobuf message, indicating whether it passes or not.  An error is
// returned if message decoding fails, but filtering is still performed against
// the partially decoded message.
func (m *Machine) Filter(message []byte) (bool, error) {
	m.reset(message)
	err := m.decode()
	ok := m.evaluate()
	return ok, err
}

// GetRawValue can be used after a Filter call to retrieve values of the
// protobuf message's fields that are referenced by the filter program.  The
// interpretation of a value depends on the field.
func (m *Machine) GetRawValue(index uint8) (value uint64, found bool) {
	slot := index >> 6
	bit := index & 63
	if m.fieldmask[slot]&(1<<bit) == 0 {
		return
	}

	value = m.fielddata[index]
	found = true
	return
}

// reset internal machine state for processing a new message.
func (m *Machine) reset(protobuf []byte) {
	m.status = false
	for i := 0; i < len(m.reg); i++ {
		m.reg[i] = 0
	}
	m.protobuf = protobuf
	for i := 0; i < len(m.fielddata); i++ {
		m.fielddata[i] = 0
	}
	for i := 0; i < len(m.fieldmask); i++ {
		m.fieldmask[i] = 0
	}
	if m.topfieldrep != nil {
		for i := uint8(0); i <= m.maxarrindex; i++ {
			m.topfieldrep[i] = 0
		}
	}
}

// repMapPool is a memory pool for use during protobuf message decoding.  It
// will theoretically grow infinitely large, but in practise its total memory
// usage is bounded by the largest or most complex protobuf message being
// filtered.  All memory is returned to the pool between each message.
type repMapPool []map[int32]int32

func (p *repMapPool) get() map[int32]int32 {
	if len(*p) > 0 {
		m := (*p)[len(*p)-1]
		*p = (*p)[:len(*p)-1]
		return m
	}
	return make(map[int32]int32)
}

func (p *repMapPool) put(m map[int32]int32) {
	for k := range m {
		m[k] = 0
	}
	*p = append(*p, m)
}
