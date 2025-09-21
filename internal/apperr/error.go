package apperr

type Error string

// implemeinting error interface
func (e Error) Error() string {
	return string(e)
}
