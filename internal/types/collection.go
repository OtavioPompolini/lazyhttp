package types

type Collection struct {
	Name   string
	IsOpen bool

	head *Request
	tail *Request
}

func NewCollection(requests []*Request) *Collection {
	var head *Request
	var tail *Request
	var prev *Request

	for i, req := range requests {
		if i == 0 {
			head = req
		}

		if i == len(requests)-1 {
			tail = req
		}

		req.Prev = prev
		if prev != nil {
			prev.Next = req
		}

		prev = req
	}

	return &Collection{
		head: head,
		tail: tail,
	}
}
