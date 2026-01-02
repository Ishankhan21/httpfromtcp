package server

import (
	"bytes"
	"fmt"
	"httfromtcp/internal/request"
	"httfromtcp/internal/response"
	"io"
	"log"
	"net"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	con     net.Listener
	handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {

	uri := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", uri)
	if err != nil {
		log.Fatal(err)
	}
	server := &Server{con: listener, handler: handler}
	go server.listen()

	// fmt.Println("Listenting at", uri)
	return server, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.con.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.handleCon(conn)
	}
}

func (s *Server) Close() {
	s.Close()
}

func (s *Server) handleCon(con net.Conn) {
	defer con.Close()

	request, err := request.RequestFromReaderHttp(con)
	if err != nil {
		defaultHeaders := response.GetDefaultHeaders(0)
		response.WriteStatusLine(con, response.BadRequest)
		response.WriteHeaders(con, defaultHeaders)
		return
	}

	writer := bytes.NewBuffer([]byte{})

	handlerError := s.handler(writer, request)
	if handlerError != nil {
		defaultHeaders := response.GetDefaultHeaders(0)
		response.WriteHeaders(con, defaultHeaders)
		response.WriteStatusLine(con, handlerError.StatusCode)
		return
	}

	body := writer.Bytes()
	defaultHeaders := response.GetDefaultHeaders(len(body))
	response.WriteStatusLine(con, 200)
	response.WriteHeaders(con, defaultHeaders)

	con.Write([]byte("\r\n"))
	con.Write(body)
}
