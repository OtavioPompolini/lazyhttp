package request

const (
	REQUEST_GET    string = "get"
	REQUEST_POST          = "post"
	REQUEST_DELETE        = "delete"
	REQUEST_PUT           = "put"
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

func NewRequest(name string) Request {
	return Request{Name: name}
}

func (a *Request) GetRequests() []Request {
	return []Request{
		{Name: "COMINI VIADO"},
		{Name: "COMINI VIADO"},
		{Name: "COMINI VIADO"},
		{Name: "COMINI VIADO"},
		{Name: "COMINI VIADO"},
	}
}
