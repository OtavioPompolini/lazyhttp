package model

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/OtavioPompolini/project-postman/internal/utils"
)

// V1 = Only name and body
type Request struct {
	Id   int64
	Name string
	Body         string
	LastResponse string
}

// Refactor this pls
func (r *Request) Execute() {
	httpRequest, err := utils.ParseHttpRequest(r.Body)
	if err != nil {
		log.Panic("Parse error")
	}

	client := http.Client{}
	res, err := client.Do(httpRequest)
	if err != nil {
		log.Panic("req error")
	}

	s, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic("XISDE")
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, s, "", "  ")
	if err != nil {
		log.Panic("DUMB")
	}

	r.LastResponse = pretty.String()
}
