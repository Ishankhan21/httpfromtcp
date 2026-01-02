package response

import (
	"errors"
	"fmt"
	"httfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	Okay                StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case Okay:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case BadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case InternalServerError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		return errors.New("Invalid status code in response")
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for k, v := range headers {
		headerLine := fmt.Sprintf("%s: %s\r\n", k, v)
		w.Write([]byte(headerLine))
	}
	return nil
}
