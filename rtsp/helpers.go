package rtsp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Header map[string]string

type Response struct {
	StatusCode int
	Headers    Header
	Body       []byte
}

func (r *Response) String() string {
	s := strconv.Itoa(r.StatusCode)
	for k, v := range r.Headers {
		s += "\n" + k + ": " + v
	}
	if r.Body != nil {
		s += "\n\n" + string(r.Body)
	}
	return s
}

func (r *Response) Bytes() []byte {
	s := fmt.Sprintf(
		"RTSP/1.0 %d %s\r\n", r.StatusCode, http.StatusText(r.StatusCode),
	)
	for k, v := range r.Headers {
		s += k + ": " + v + "\r\n"
	}
	s += "\r\n"
	if r.Body != nil {
		return append([]byte(s), r.Body...)
	} else {
		return []byte(s)
	}
}

func ReadResponse(r *bufio.Reader) (*Response, error) {
	line, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	s := strings.SplitN(string(line), " ", 3)
	if len(s) != 3 {
		return nil, errors.New("wrong response")
	}

	code, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, err
	}

	resp := &Response{StatusCode: code}
	resp.Headers, resp.Body, err = readHeader(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func readHeader(r *bufio.Reader) (Header, []byte, error) {
	header := Header{}
	var body []byte

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return nil, nil, err
		}
		if len(line) == 0 {
			break
		}
		s := strings.SplitN(string(line), ": ", 2)
		if len(s) != 2 {
			return nil, nil, errors.New("wrong header")
		}
		header[strings.ToLower(s[0])] = s[1]
	}

	if slen := header["content-length"]; slen != "" {
		ilen, err := strconv.Atoi(slen)
		body = make([]byte, ilen)
		if _, err = io.ReadFull(r, body); err != nil {
			return nil, nil, err
		}
	}
	return header, body, nil
}
