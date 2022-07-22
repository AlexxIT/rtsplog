package rtsp

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type Conn struct {
	Out func(msg interface{})

	url  *url.URL
	conn net.Conn
	auth string
	seq  int
}

func (c *Conn) Dial(rawURL string) (err error) {
	c.url, err = url.Parse(rawURL)
	if err != nil {
		return
	}

	// remove auth from url
	user := c.url.User
	c.url.User = nil

	c.conn, err = net.DialTimeout(
		"tcp", c.url.Host, 10*time.Second,
	)
	if err != nil {
		return
	}

	var tlsConf *tls.Config
	switch c.url.Scheme {
	case "rtsps":
		tlsConf = &tls.Config{ServerName: c.url.Hostname()}
	case "rtspx":
		c.url.Scheme = "rtsps"
		tlsConf = &tls.Config{InsecureSkipVerify: true}
	}
	if tlsConf != nil {
		tlsConn := tls.Client(c.conn, tlsConf)
		if err = tlsConn.Handshake(); err != nil {
			return
		}
		c.conn = tlsConn
	}

	reader := bufio.NewReader(c.conn)

	// [RTSP] OPTIONS step
	if err = c.Request("OPTIONS " + c.url.String()); err != nil {
		return
	}

	var res *Response
	res, err = ReadResponse(reader)
	if err != nil {
		return
	}

	c.Out(res)

	if res.StatusCode != 200 {
		return errors.New("wrong response on OPTIONS")
	}

	// [RTSP] DESCRIBE step (with auth)
	for {
		if err = c.Request(
			"DESCRIBE "+c.url.String(),
			"Accept: application/sdp",
			"Require: www.onvif.org/ver20/backchannel",
		); err != nil {
			return
		}

		res, err = ReadResponse(reader)
		if err != nil {
			return
		}

		c.Out(res)

		if res.StatusCode == 401 {
			if err = c.includeAuth(res, user); err != nil {
				return
			}
		} else {
			break
		}
	}
	if res.StatusCode != 200 {
		return errors.New("wrong response on DESCRIBE")
	}

	if val := res.Headers["content-base"]; val != "" {
		c.url, err = url.Parse(val)
		if err != nil {
			return
		}
	}

	return
}

func (c *Conn) Request(header ...string) (err error) {
	c.seq++
	header[0] += " RTSP/1.0" +
		fmt.Sprintf("\r\nCSeq: %d", c.seq) +
		c.auth
	data := []byte(strings.Join(header, "\r\n") + "\r\n\r\n")
	_, err = c.conn.Write(data)
	return
}

func (c *Conn) includeAuth(res *Response, user *url.Userinfo) error {
	// second try error
	if c.auth != "" {
		return errors.New("wrong auth")
	}
	password, set := user.Password()
	// user don't have password
	if !set {
		return errors.New("password not set")
	}
	if strings.Index(res.Headers["www-authenticate"], "Basic") == 0 {
		c.auth = "\r\nAuthorization: Basic " +
			base64.StdEncoding.EncodeToString(
				[]byte(user.Username()+":"+password),
			)
	} else {
		return errors.New("unsupported auth")
	}
	return nil
}
