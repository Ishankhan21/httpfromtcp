package request

import (
	"bytes"
	"errors"
	"fmt"
	"httfromtcp/internal/headers"
	"io"
	"strconv"
)

type ParserState string

const (
	DONE        ParserState = "DONE"
	INITIALIZED ParserState = "INIT"
	HEADERS     ParserState = "HEADERS"
	BODY        ParserState = "BODY"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	ParserState ParserState
}

func (r *Request) parseRequestBody(data []byte) (bool, error) {
	bodyLength := r.Headers.Get("content-length")
	bodyLengthInt, _ := strconv.Atoi(bodyLength)
	if len(r.Body) < bodyLengthInt {
		r.Body = append(r.Body, data...)
	}
	if len(r.Body) == bodyLengthInt {
		return true, nil
	}
	if len(r.Body) > bodyLengthInt {
		return false, errors.New("body legth more then expected")
	}
	return false, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, int, error) {

	lineSepratorBytes := []byte("\r\n")
	sepratorIndex := bytes.Index(data, lineSepratorBytes)
	if sepratorIndex == -1 {
		fmt.Println("/n not found", len(data))
		return nil, 0, 0, nil
	}

	lines := bytes.Split(data, lineSepratorBytes)
	firstLineSplit := bytes.Split(lines[0], []byte(" "))

	if len(firstLineSplit) < 3 {
		return nil, 0, 0, errors.New("invalid first line")
	}
	// for i, v := range firstLineSplit {
	// 	// fmt.Println(string(v))
	// 	if i == 2 {
	// 		fmt.Println(string(bytes.Split(v, []byte("/"))[1]))
	// 	}
	// }
	return &RequestLine{
		Method:        string(firstLineSplit[0]),
		RequestTarget: string(firstLineSplit[1]),
		HttpVersion:   (string(bytes.Split(firstLineSplit[2], []byte("/"))[1])),
	}, 1, sepratorIndex + 2, nil
}

func (request *Request) parse(data []byte) (int, error) {
	var consumedIndex int
	if request.ParserState != INITIALIZED && request.ParserState != HEADERS && request.ParserState != DONE && request.ParserState != BODY {
		request.ParserState = INITIALIZED
		requestLine, res, byteConsumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if res == 0 {
			return 0, nil
		}
		if byteConsumed > 0 {
			request.RequestLine = *requestLine
			request.ParserState = HEADERS
		}
		return byteConsumed, nil
	}

	if request.ParserState == HEADERS {
		n, done, err := request.Headers.Parse(data)
		consumedIndex = n
		if err != nil {
			return n, err
		}
		if done {
			if request.Headers.Get("content-length") != "" {
				request.ParserState = BODY
			} else {
				request.ParserState = DONE
				return n, nil
			}
		} else {
			return n, nil
		}
	}

	if request.ParserState == BODY {
		done, err := request.parseRequestBody(data[consumedIndex:])
		if err != nil {
			return 0, err
		}
		if done {
			request.ParserState = DONE
			return len(data), nil
		}
	}

	if request.ParserState == DONE {
		fmt.Println("final request AFter done ++++++", request.RequestLine, request.Headers)
	}

	return len(data), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request

	// 50 bytes a time
	buf := make([]byte, 50)
	request.Headers = headers.NewHeaders()
	var bufIndexToRead int
	var n int
	var err error
	for request.ParserState != DONE {
		n, err = reader.Read(buf[bufIndexToRead:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		bytesConsumed, err := request.parse(buf[:bufIndexToRead+n])
		if err != nil {
			return nil, err
		}
		if bytesConsumed == bufIndexToRead+n {
			tempBuf := make([]byte, 2000)
			buf = tempBuf
			bufIndexToRead = 0
		} else if bytesConsumed > 0 {
			tempBuf := make([]byte, 2000)
			copy(tempBuf, buf[bytesConsumed:])
			buf = tempBuf
			bufIndexToRead = (bufIndexToRead + n) - bytesConsumed
		}
	}
	fmt.Println("findal request", len(request.RequestLine.Method))
	return &request, nil
}

func RequestFromReaderHttp(reader io.Reader) (*Request, error) {
	var request Request
	request.Headers = headers.NewHeaders()

	var accumulated []byte
	buf := make([]byte, 1024)

	for request.ParserState != DONE {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n > 0 {
			accumulated = append(accumulated, buf[:n]...)

			for len(accumulated) > 0 && request.ParserState != DONE {
				bytesConsumed, parseErr := request.parse(accumulated)
				if parseErr != nil {
					return nil, parseErr
				}

				if bytesConsumed == 0 {
					break
				}

				accumulated = accumulated[bytesConsumed:]
			}
		}

		if err == io.EOF {
			break
		}
	}

	if request.ParserState != DONE {
		return nil, fmt.Errorf("incomplete request")
	}

	fmt.Println("final request", len(request.RequestLine.Method))
	return &request, nil
}
