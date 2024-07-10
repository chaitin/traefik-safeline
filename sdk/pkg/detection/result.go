package detection

type ResultObjective int

const (
	RO_REQUEST  ResultObjective = 0
	RO_RESPONSE ResultObjective = 1
)

type Result struct {
	Objective   ResultObjective
	Head        byte
	Body        []byte
	Alog        []byte
	ExtraHeader []byte
	ExtraBody   []byte
	T1KContext  []byte
	Cookie      []byte
	WebLog      []byte
}

func (r *Result) Passed() bool {
	return r.Head == '.'
}

func (r *Result) Blocked() bool {
	return !r.Passed()
}
