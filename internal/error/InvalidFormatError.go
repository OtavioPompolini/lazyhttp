package error

type InvalidFormatError struct{}

func (ife *InvalidFormatError) Error() string {
	return "Invalid format error"
}
