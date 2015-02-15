

// Package github.comeid/omeid/httptest extends net/http/httptest with useful methods for HTTP testing.
package httptest

import (
	"bytes"
	"encoding/json"
	"log"
	"mime"
	"net/http/httptest"
	"reflect"
)

// NewRecorder returns an initialized ResponseRecorder, it's compatiable with the official
// httptest.ResponseRecoder by embedding it.
func NewRecorder() *ResponseRecoder {
	return &ResponseRecoder{httptest.NewRecorder()}
}


// ResponseRecorder is an extension of httptest.ResponseRecoder which  implementation of http.ResponseWriter that records its mutations for later inspection in tests. 
type ResponseRecoder struct {
	*httptest.ResponseRecorder
}

// CheckCode checks returns the response status code and if it matches the expected status code.
// This is mostly provided for consistency with other method and http's package refering to Status Code as Status, you can directly acess the Code as `ResponseRecoder.Code` if you may.
func (r *ResponseRecoder) ExpectStatus(expect int) (int, bool) {
	return r.Code, r.Code == expect
}

// ExpectBytes check the response body against the provided []byte, if strict is set to false,
// It will ignore a possible trialing '\n' in response body.
func (r *ResponseRecoder) ExpectBytes(expect []byte, strict bool) ([]byte, bool) {
	body := r.Body.Bytes()
	if bytes.Compare(body, expect) == 0 {
		return body, true
	}
	if !strict && bytes.Compare(body, append(expect, '\n')) == 0 {
		return body, true
	}
	return body, false
}

// ExpectJSON checks if decoding the body to `model` will match the `expect` object.
// Providing mismatching Model and Expect will result into a fatal error.
func (r *ResponseRecoder) ExpectJSON(model, expect interface{}) ([]byte, bool) {

	mt := reflect.TypeOf(model)
	me := reflect.TypeOf(expect)
	if me != mt {
		log.Fatal("Model and Expect mismatch: Model: %s. Expect: %s.", mt, me)
	}
	err := json.Unmarshal(r.Body.Bytes(), model)
	return r.Body.Bytes(), err == nil && reflect.DeepEqual(model, expect)
}


// Checks if the Response has the expected content type, returns the content type
// And if it is matching.
func (r *ResponseRecoder) ExpectContentType(expect string) (string, bool) {
	v := r.HeaderMap.Get("Content-Type")
	if v == "" {
		return v, false
	}
	t, _, err := mime.ParseMediaType(v)
	if err != nil {
		return err.Error(), false
	}
	return t, t == expect
}
