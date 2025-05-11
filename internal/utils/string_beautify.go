package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/go-xmlfmt/xmlfmt"
)

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func IsXml(str string) bool {
	var js string
	return xml.Unmarshal([]byte(str), &js) == nil
}

func StringBeautify(str string) func(io.Writer) error {
	return func(wr io.Writer) error {
		// log.Printf("Beautifying string = %s", str)
		isJson := IsJSON(str)
		isXml := IsXml(str)
		var lexer chroma.Lexer

		if isJson {
			log.Print("String is a json")
			lexer = lexers.Get("json")
			var pretty bytes.Buffer
			err := json.Indent(&pretty, []byte(str), "", "  ")
			str = pretty.String()
			if err != nil {
				log.Printf("Error identing json. json = %s", str)
			}
		} else if isXml {
			lexer = lexers.Get("xml")
			log.Print("String is a xml")
			str = xmlfmt.FormatXML(str, "", "  ")
		}

		if lexer == nil {
			lexer = lexers.Analyse(str)
		}

		if lexer == nil {
			lexer = lexers.Fallback
		}
		lexer = chroma.Coalesce(lexer)

		style := styles.Get("tokyonight-night")
		if style == nil {
			log.Printf("Failed to get style. Setting fallback. Style = %s", styles.Fallback.Name)
			style = styles.Fallback
		}

		formatter := formatters.Get("terminal16m")
		if formatter == nil {
			log.Printf("Failed to get formatter. Setting fallback. Fomatter = %s", formatters.Fallback)
			formatter = formatters.Fallback
		}

		iterator, err := lexer.Tokenise(nil, str)
		if err != nil {
			log.Print("Failed to tokenize.", err)
			return err
		}

		err = formatter.Format(wr, style, iterator)
		if err != nil {
			log.Printf("Error while beautifying string. str = %s. err: %s", str, err)
			return err
		}

		return nil
	}
}

// NO. THIS WAS NOT AI GENERATED. KAPPA
func indentXML(rawXML string, prefix string, indent string) (string, error) {
	// Create input reader and output buffer
	reader := strings.NewReader(rawXML)
	var output bytes.Buffer

	// Create XML decoder and encoder
	decoder := xml.NewDecoder(reader)
	encoder := xml.NewEncoder(&output)
	encoder.Indent(prefix, indent)

	// Process tokens
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("xml decode error: %w", err)
		}

		if err := encoder.EncodeToken(token); err != nil {
			return "", fmt.Errorf("xml encode error: %w", err)
		}
	}

	// Finalize encoding
	if err := encoder.Flush(); err != nil {
		return "", fmt.Errorf("xml flush error: %w", err)
	}

	return output.String(), nil
}
