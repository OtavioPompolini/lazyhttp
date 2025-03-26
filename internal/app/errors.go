package app

type InvalidFormatError struct{}

func (ife *InvalidFormatError) Error() string {
	return "Invalid format error"
}

type InvalidRequestPosition struct{}

func (ife *InvalidRequestPosition) Error() string {
	return "Request position out of bounds"
}
