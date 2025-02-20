package request

type Adapter interface {
	GetRequests() []Request
	CreateRequest(name string) Request
}

