package main

import "time"

type ContentType string

const (
	TextHtml  ContentType = "text/html"
	ImageJpeg ContentType = "image/jpeg"
	ImagePng  ContentType = "image/png"
)

type HTTPResponseHeader struct {
	HTTPProtocol  string
	StatusCode    int
	Server        string
	Date          time.Time
	ContentType   ContentType
	ContentLength int
}

type HTTPResponse struct {
	HTTPHeader HTTPResponseHeader
	Body       []byte
}

type HTTPRequest struct {
	Method   string
	Resource string
	Body     string
}

const CRLF = "\r\n"

var HTTPRequestMethods = [...]string{"GET", "POST", "DELETE", "PUT", "PATCH"}

type ResourceResponse struct {
	Body             []byte
	ResponseBodyType ContentType
}

type ResourcePath string

type Resources map[ResourcePath]ResourceResponse
