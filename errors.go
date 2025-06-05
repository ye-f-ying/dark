package dark

type DarkError struct {
	error
}

func (m *DarkError) Error() string {
	return ""
}
