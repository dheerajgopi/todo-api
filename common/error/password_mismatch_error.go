package error

// PasswordMismatchError is returned if password is not matching with the decrypted hash
type PasswordMismatchError struct {
}

func (pme *PasswordMismatchError) Error() string {
	return "password mismatch"
}
