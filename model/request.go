package model

const (
	REQUEST_GET    string = "get"
	REQUEST_POST          = "post"
	REQUEST_DELETE        = "delete"
	REQUEST_PUT           = "put"
)

// V1 = Only name and body
type Request struct {
	Id int64
	Name string
	// Url     string
	// Method  string
	// Headers []Header
	Body string
}

type Header struct {
	key   string
	value string
}

