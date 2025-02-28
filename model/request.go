package model

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/OtavioPompolini/project-postman/utils"
)

const (
	GET    string = "get"
	POST          = "post"
	DELETE        = "delete"
	PUT           = "put"
)


// type Header struct {
// 	key   string
// 	value string
// }

// V1 = Only name and body
type Request struct {
	Id int64
	Name string
	// Url     string
	// Method  string
	// Headers []Header
	Body string
	LastResponse string
}

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


	// THIS WILL REMAIN IN HERE TO SHOW HOW DUMB I'AM
	// for _, v := range s {
	// 	if v == "}" {
	// 		ident -= 1
	// 	}
	//
	// 	if v == "]" {
	// 		ident -= 1
	// 	}
	//
	// 	final += strings.Repeat(spaces, ident)
	// 	final += v
	// 	if v == "{" {
	// 		ident += 1
	// 		final += "\n"
	// 	}
	//
	// 	if v == "[" {
	// 		ident += 1
	// 		final += "\n"
	// 	}
	//
	// 	if v == "}" || v == "]" {
	// 		final += "\n"
	// 	}
	//
	//
	// 	if v == "," {
	// 		final += "\n"
	// 	}
	// }

// 	for i := range len(s) {
// 	 	final += strings.Repeat(spaces, ident)
//
// 		for s[i] != "{" &&
// 		s[i] != "}" &&
// 		s[i] != "[" &&
// 		s[i] != "]" &&
// 		s[i] != "," {
// 			final += s[i]
// 		}
//
// 		if s[i] == "," {
// 			final += s[i]
// 			i += 1
// 		}
// 		final += "\n"
//
// 		if s[i] == "}" ||
// 		s[i] == "]" {
// 			ident -= 1
// 			final += strings.Repeat(spaces, ident)
// 			final += s[i]
// 			final += "\n"
// 			i+=1
// 		}
//
// 		if s[i] == "{" ||
// 		s[i] == "[" {
// 			ident += 1
// 			final += strings.Repeat(spaces, ident)
// 			final += s[i]
// 			final += "\n"
// 			i+=1
// 		}
// 	}

	// for i, v := range s {
	// 	if v != '{' &&
	// 	v != '}' &&
	// 	v != '[' &&
	// 	v != ']' &&
	// 	v != ',' {
	// 		final += string(v)
	// 	}
	//
	// 	if v == '{' || v == '[' {
	// 		final += string(v)
	// 		ident += 1
	// 		final += "\n"
	// 		final += strings.Repeat(spaces, ident)
	// 	}
	//
	// 	if v == '}' || v == ']' {
	// 		final += "\n"
	// 		ident -= 1
	// 		final += strings.Repeat(spaces, ident)
	// 		final += string(v)
	// 		if s[i-1] != '}' && s[i-1] != ']' {
	// 			final += "\n"
	// 		}
	// 		final += strings.Repeat(spaces, ident)
	// 	}
	//
	// 	if v == ',' {
	// 		final += string(v)
	// 		final += "\n"
	// 		final += strings.Repeat(spaces, ident)
	// 	}
	//
	// }
	// r.LastResponse = final
}
