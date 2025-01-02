package pkg

type ErrorWithStatusCode struct {
	Code int
	Err  error
}

func (errCode ErrorWithStatusCode) Error() string {
	return errCode.Err.Error()
}
