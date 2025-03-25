package types

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/OtavioPompolini/project-postman/internal/utils"
)

// V1 = Only name and body
type Request struct {
	Id           int64
	Name         string
	Body         string
	LastResponse *Response
	Next         *Request
	Prev         *Request
}

type Response struct {
	Info string
	Body string
}

// Refactor this pls
// Didnt like this Request.Execute(). Where can I set configs in my app
// Imagine if I want to force http/1.0 instead of http/2.0
func (r *Request) Execute() error {
	httpRequest, err := utils.ParseHttpRequest(r.Body)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(httpRequest)
	if err != nil {
		return err
	}

	responseString := ""

	responseString += res.Proto + " "
	responseString += res.Status
	responseString += "\n"

	for k, v := range res.Header {
		responseString += k + ": "
		responseString += strings.Join(v, "")
		responseString += "\n"
	}

	responseString += "\n"
	s, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, s, "", "  ")
	if err != nil {
		return err
	}

	r.LastResponse = &Response{
		Info: responseString,
		Body: pretty.String(),
	}

	return nil
}
