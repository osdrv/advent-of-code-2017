package main

type Machine struct {
	tape  map[int]bool
	pos   int
	state byte
}

func NewMachine() *Machine {
	return &Machine{
		tape:  make(map[int]bool),
		pos:   0,
		state: 'A',
	}
}

func (m *Machine) Curr() bool {
	return m.tape[m.pos]
}

func (m *Machine) Flip() *Machine {
	m.tape[m.pos] = !m.tape[m.pos]
	return m
}

func (m *Machine) Left() *Machine {
	m.pos--
	return m
}

func (m *Machine) Right() *Machine {
	m.pos++
	return m
}

func (m *Machine) ToState(state byte) *Machine {
	m.state = state
	return m
}

func (m *Machine) RunSample() {
	switch m.state {
	case 'A':
		if !m.tape[m.pos] {
			m.tape[m.pos] = true
			m.pos++
			m.state = 'B'
		} else {
			m.tape[m.pos] = false
			m.pos--
			m.state = 'B'
		}
	case 'B':
		if !m.tape[m.pos] {
			m.tape[m.pos] = true
			m.pos--
			m.state = 'A'
		} else {
			m.tape[m.pos] = true
			m.pos++
			m.state = 'A'
		}
	}
}

func (m *Machine) Run() {
	switch m.state {
	case 'A':
		if !m.Curr() {
			m.Flip().Right().ToState('B')
		} else {
			m.Flip().Left().ToState('E')
		}
	case 'B':
		if !m.Curr() {
			m.Flip().Left().ToState('C')
		} else {
			m.Flip().Right().ToState('A')
		}
	case 'C':
		if !m.Curr() {
			m.Flip().Left().ToState('D')
		} else {
			m.Flip().Right().ToState('C')
		}
	case 'D':
		if !m.Curr() {
			m.Flip().Left().ToState('E')
		} else {
			m.Flip().Left().ToState('F')
		}
	case 'E':
		if !m.Curr() {
			m.Flip().Left().ToState('A')
		} else {
			m.Left().ToState('C')
		}
	case 'F':
		if !m.Curr() {
			m.Flip().Left().ToState('E')
		} else {
			m.Right().ToState('A')
		}
	default:
		panic("Unknown state")
	}
}

func computeChecksum(m *Machine) int {
	cnt := 0
	for _, isSet := range m.tape {
		if isSet {
			cnt++
		}
	}
	return cnt
}

const (
	STEPS = 12386363
)

func main() {
	m := NewMachine()
	for i := 0; i < STEPS; i++ {
		m.Run()
	}
	printf("checksum: %d", computeChecksum(m))
}
