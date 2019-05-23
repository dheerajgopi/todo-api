package error

// UnauthorizedError is returned in case access is denied to some services
type UnauthorizedError struct {
}

func (ue *UnauthorizedError) Error() string {
	return "unauthorized"
}