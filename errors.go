package dark

type DarkErrorType string

const (
	ErrAntsNoInit DarkErrorType = "ants no init"
	ErrAntsIsNil  DarkErrorType = "ants is nil"
)

type DarkError struct {
	//error
	err DarkErrorType
}

func (m *DarkError) Error() string {
	return string(m.err)
}
