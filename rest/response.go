package rest

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
)

const (
	APPLICATION_JSON string = "application/json"
	APPLICATION_XML  string = "application/xml"
	TEXT_JSON        string = "text/json"
	TEXT_HTML        string = "text/html"
	TEXT_XML         string = "text/xml"
)

var (
	responseUsed = errors.New("response already written to")
)

// Send returns data to the client.
// If the request has the Accept header of "text/xml" or "application/xml" then
// the response will be in XML, otherwise in JSON.
func (r *Rest) Send() error {
	if r.sent {
		return responseUsed
	}

	r.sent = true

	if r.status <= 0 {
		r.status = 200
	}

	// Force the Content-Type if the response contentType is not set
	if r.contentType == "" {
		r.contentType = r.GetHeader("Accept")
	}
	if r.contentType == "" {
		r.contentType = APPLICATION_JSON
	}
	r.AddHeader("Content-Type", r.contentType)

	// Until we get CORS handling correctly
	r.AddHeader("Access-Control-Allow-Origin", "*")

	// Write the headers
	h := r.writer.Header()
	for k, v := range r.headers {
		h.Add(k, v)
	}

	// Write the status
	r.writer.WriteHeader(r.status)

	// Write from a reader
	if r.reader != nil {
		if closer, ok := r.reader.(io.ReadCloser); ok {
			defer closer.Close()
		}

		_, err := io.Copy(r.writer, r.reader)
		return err
	} else if r.value != nil {
		if ba, ok := r.value.([]byte); ok {
			_, err := r.writer.Write(ba)
			return err
		} else {
			// Finally the content, encode if an object

			isXml := r.contentType == TEXT_XML || r.contentType == APPLICATION_XML
			isJson := r.contentType == TEXT_JSON || r.contentType == APPLICATION_JSON

			// Ensure we have a valid contentType default to APPLICATION_JSON if not
			if !isXml && !isJson {
				isJson = true
			}
			if isXml {
				return xml.NewEncoder(r.writer).Encode(r.value)
			} else {
				return json.NewEncoder(r.writer).Encode(r.value)
			}
		}
	}

	return nil
}
