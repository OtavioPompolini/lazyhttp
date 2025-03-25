package types

import "time"

// V1 = Only name and body
type Request struct {
	Id              int64
	Name            string
	Body            string
	ResponseHistory []*Response
	Next            *Request
	Prev            *Request
}

type Response struct {
	Id         int64
	RequestId  int64
	Info       string
	Body       string
	Created_at time.Time
}
