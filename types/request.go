package types

const (
	GET    string = "get"
	POST          = "post"
	DELETE        = "delete"
	PUT           = "put"
)

type Request struct {
	Name    string
	Url     string
	Method  string
	Headers []Header
	Body    string //??
}

type Header struct {
	key   string
	value string
}

func NewRequest(name string) *Request {
	return &Request{Name: name}
}
