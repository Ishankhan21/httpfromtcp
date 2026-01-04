# HTTP from TCP

## Summary

As backend developers, we create and update HTTP servers regularly. To build efficient and high-performance services, we must understand what actually happens when we create an HTTP server, how requests are delivered to our server, and how HTTP requests are formed.

In this project, we build HTTP servers from the ground up, starting at the TCP layer. This helps us understand the relationship between HTTP and TCP, and how HTTP is implemented on top of TCP.

## How HTTP Works Over TCP

Whenever we create an HTTP server and a client sends an HTTP request to this server, the following fundamental steps occur:

1. **TCP Connection Establishment**: First, the client and server create a TCP connection by following all the TCP handshake procedures.

2. **Data Transmission**: Once the TCP connection is established, the client sends the request in the form of data packets to the server, which are transmitted as binary/bytes at the TCP connection level.

3. **HTTP Format**: Since this is an HTTP request over TCP, there is a specific format that defines how HTTP requests should be structured, as specified in [RFC 9112](https://www.rfc-editor.org/rfc/rfc9112.html) (HTTP/1.1 Message Syntax) and [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html) (HTTP Semantics).

4. **Request Parsing**: The server reads the TCP binary data, converts it to a string, and parses the appropriate HTTP request components including the HTTP version, endpoint, headers, body, etc.

5. **Response Handling**: After properly parsing the request, the server invokes the appropriate controller/handler based on the endpoint and returns the response to the client in the specific format specified in the RFC.

## Learning Objectives

- Understand the relationship between HTTP and TCP protocols
- Learn how HTTP requests are formatted and transmitted at the TCP level
- Gain insight into request parsing and response handling
- Build a deeper understanding of network protocols for creating high-performance services

## References

- [RFC 9110 - HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html)
- [RFC 9112 - HTTP/1.1 Message Syntax and Routing](https://www.rfc-editor.org/rfc/rfc9112.html)