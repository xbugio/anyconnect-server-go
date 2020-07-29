package server

import (
	"bytes"
	"net/http"
)

type bufferCloser struct {
	bytes.Buffer
}
type response struct {
	resp *http.Response
	body *bufferCloser
}

func (buff *bufferCloser) Close() error {
	buff.Reset()
	return nil
}

func newResponse(request *http.Request) *response {
	r := new(response)
	r.resp = &http.Response{
		StatusCode:    http.StatusOK,
		Header:        make(http.Header),
		Trailer:       make(http.Header),
		ContentLength: 0,
		Request:       request,
	}
	return r
}

func (r *response) Header() http.Header {
	return r.resp.Header
}

func (r *response) WriteHeader(statuscode int) {
	r.resp.StatusCode = statuscode
}

func (r *response) Write(data []byte) (int, error) {
	if r.body == nil {
		r.body = new(bufferCloser)
		r.resp.Body = r.body
	}
	n, err := r.body.Write(data)
	r.resp.ContentLength += int64(n)
	return n, err
}
