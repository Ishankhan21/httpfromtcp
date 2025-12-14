package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type ParserState int

const (
	DONE        ParserState = 1
	INITIALIZED ParserState = 0
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	ParserState ParserState
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	lineSepratorBytes := []byte("\r\n")
	if bytes.Index(data, lineSepratorBytes) == -1 {
		fmt.Println("/n not found", len(data))
		return nil, 0, nil
	}

	lines := bytes.Split(data, lineSepratorBytes)
	fmt.Println("lines +++", len(lines))
	firstLineSplit := bytes.Split(lines[0], []byte(" "))

	if len(firstLineSplit) < 3 {
		return nil, 0, errors.New("invalid first line")
	}
	// for i, v := range firstLineSplit {
	// 	// fmt.Println(string(v))
	// 	if i == 2 {
	// 		fmt.Println(string(bytes.Split(v, []byte("/"))[1]))
	// 	}
	// }
	fmt.Println(len(firstLineSplit))
	return &RequestLine{
		Method:        string(firstLineSplit[0]),
		RequestTarget: string(firstLineSplit[1]),
		HttpVersion:   (string(bytes.Split(firstLineSplit[2], []byte("/"))[1])),
	}, 1, nil
}

func (request *Request) parse(data []byte) (int, error) {
	request.ParserState = INITIALIZED

	requestLine, res, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if res == 0 {
		return 0, nil
	}
	fmt.Println("requestLine +++++", requestLine.Method, requestLine.RequestTarget, requestLine.HttpVersion)

	request.RequestLine = *requestLine
	request.ParserState = DONE
	return len(data), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request

	buf := make([]byte, 1024)
	var bufIndexToRead int
	for request.ParserState != DONE {
		n, err := reader.Read(buf[bufIndexToRead:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		_, err = request.parse(buf[:bufIndexToRead+n])
		if err != nil {
			return nil, err
		}
		bufIndexToRead += n
	}
	fmt.Println("findal request", len(request.RequestLine.Method))
	return &request, nil
}
