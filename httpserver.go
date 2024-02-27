package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/matheush9/tcp-server/server"
)

const (
	HTTPProtocolVersion = "1.1"
	MyServerName        = "matheusserver"
)

func Start(rm Resources, sc *server.ServerConfig) {
	sc.SetupServer()
	listenRequests(sc, rm)
}

func NewServerConfig(connectionTimeout float64, ip string, port int) server.ServerConfig {
	return server.ServerConfig{ConnectionTimeout: connectionTimeout, IPV4Address: ip, Port: port}
}

func listenRequests(sc *server.ServerConfig, rm Resources) {
	clientRequestChan := make(chan server.ClientRequest)
	go sc.HandleConnections(clientRequestChan)
	for clientRequest := range clientRequestChan {
		go handleResponse(clientRequest, sc, rm)
	}
}

func handleResponse(clientReq server.ClientRequest, sc *server.ServerConfig, rm Resources) {
	var response HTTPResponse
	parsedRequest, err := parseRequest(clientReq)
	if err != nil {
		log.Printf("error while trying to parse the request: %v", err)
	}
	resourceResponse, found := rm[ResourcePath(parsedRequest.Resource)]
	if !found {
		log.Printf("the resource asked doesn't exist: %s", parsedRequest.Resource)
		response = resourceNotFound(TextHtml)
	} else {
		response = OK(resourceResponse)
	}
	sc.SendResponse(clientReq, response.responseToString())
}

func parseRequest(request server.ClientRequest) (HTTPRequest, error) {
	splittedRequest := strings.Split(request.Request, CRLF)
	requestLines := strings.Split(splittedRequest[0], " ")
	method := requestLines[0]
	for i := range HTTPRequestMethods {
		if HTTPRequestMethods[i] == method {
			break
		} else {
			return HTTPRequest{}, errors.New("malformed http request")
		}
	}
	return HTTPRequest{
		Method:   method,
		Resource: requestLines[1],
		Body:     splittedRequest[len(splittedRequest)-1],
		//I will deal with the remaining later when I need it
	}, nil
}

func OK(rs ResourceResponse) HTTPResponse {
	return newHTTPResponse(rs.Body, 200, rs.ResponseBodyType)
}

func resourceNotFound(contentType ContentType) HTTPResponse {
	return newHTTPResponse(nil, 404, contentType)
}

func newHTTPResponse(body []byte, statusCode int, contentType ContentType) HTTPResponse {
	return HTTPResponse{
		HTTPHeader: HTTPResponseHeader{
			HTTPProtocol:  HTTPProtocolVersion,
			StatusCode:    statusCode,
			Server:        MyServerName,
			Date:          time.Now(),
			ContentType:   contentType,
			ContentLength: len([]byte(body)),
		},
		Body: body,
	}
}

func (response HTTPResponse) responseToString() string {
	res := fmt.Sprintf("HTTP/%s %d"+CRLF+
		"Date: %v"+CRLF+
		"Server: %s"+CRLF,

		response.HTTPHeader.HTTPProtocol,
		response.HTTPHeader.StatusCode,
		response.HTTPHeader.Date,
		response.HTTPHeader.Server,
	)
	if len(response.Body) > 0 {
		res += fmt.Sprintf("Content-Type: %s"+CRLF+
			"Content-Length: %d"+CRLF+
			CRLF+"%s",

			response.HTTPHeader.ContentType,
			response.HTTPHeader.ContentLength,
			response.Body,
		)
	}
	return res
}
