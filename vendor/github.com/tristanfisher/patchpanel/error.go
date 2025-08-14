package patchpanel

// NoFieldError allows for differentiating no named field vs parsing errors
type NoFieldError struct {
	Msg string
}

func (nfe NoFieldError) Error() string {
	return nfe.Msg
}

type NoValueError struct {
	Msg string
}

func (nve NoValueError) Error() string {
	return nve.Msg
}

// UnhandledParserTypeError allows for a user to handle his/her type or just use as is (e.g. if plumbing an any)
type UnhandledParserTypeError struct {
	Msg string
}

func (u UnhandledParserTypeError) Error() string {
	return u.Msg
}
