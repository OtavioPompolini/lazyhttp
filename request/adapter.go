package request

type Adapter interface {
	GetRequests() *map[int64]Request
	CreateRequest(name string) *Request
	UpdateRequest(r *Request)
}

