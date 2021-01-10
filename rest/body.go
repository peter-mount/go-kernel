package rest

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"io"
)

// BodyReader() returns Request().Body unless the Content-Encoding header is set
// to gzip in which case the body is decompressed before calling the reader
func (r *Rest) BodyReader() (io.ReadCloser, error) {
	reader := r.request.Body

	if r.request.Header.Get("Content-Encoding") == "gzip" {
		if gr, err := gzip.NewReader(reader); err != nil {
			return nil, err
		} else {
			reader = gr
		}
	}

	return reader, nil
}

// Body decodes the request body into an interface.
// If the body is compressed with Content-Encoding header set to "gzip" then the
// body is decoded first.
// If the Content-Type is "text/xml" or "application/xml" then the body is presumed
// to be in XML, otherwise json is presumed.
func (r *Rest) Body(v interface{}) error {

	contentType := r.GetHeader("Content-Type")
	isXml := contentType == TEXT_XML || contentType == APPLICATION_XML
	isJson := contentType == TEXT_JSON || contentType == APPLICATION_JSON

	// Ensure we have a valid contentType default to APPLICATION_JSON if not
	if !isXml && !isJson {
		contentType = APPLICATION_JSON
		isJson = true
	}

	if reader, err := r.BodyReader(); err != nil {
		return err
	} else if isXml {
		if err := xml.NewDecoder(reader).Decode(v); err != nil {
			return err
		}
	} else {
		if err := json.NewDecoder(reader).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
