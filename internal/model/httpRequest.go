package model

type HttpRequest struct {
	url         string
	queryParams map[string]string
	body        string //Json????
	headers     map[string]string
}
