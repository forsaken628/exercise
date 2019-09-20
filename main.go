package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	r := LineReader{bufio.NewReader(conn)}
	defer conn.Close()

	first, err := r.Line()
	if err != nil {
		return
	}

	spl := strings.Split(first, " ")
	if len(spl) != 3 {
		log.Print("valid first line")
		return
	}

	if spl[2] != "HTTP/1.1" {
		log.Print("require HTTP/1.1")
		return
	}

	req := &Request{
		Method: spl[0],
		Header: make(Header),
		URL:    spl[1],
		Body:   r,
	}

	for {
		l, err := r.Line()
		if err != nil {
			return
		}
		if l == "" {
			break
		}
		spl := strings.SplitN(l, ":", 2)
		if len(spl) != 2 {
			log.Println("valid header", l)
			return
		}
		//url.QueryUnescape() todo
		req.Header.Add(strings.Trim(spl[0], " "), strings.Trim(spl[1], " "))
	}

	fn, ok := httpHandle[req.URL]
	if !ok {
		return
	}

	fn(req, conn)
}

type LineReader struct {
	*bufio.Reader
}

func (r *LineReader) Line() (string, error) {
	buf := strings.Builder{}
	for {
		line, pre, err := r.ReadLine()
		if err != nil {
			return "", err
		}
		_, _ = buf.Write(line)
		if !pre {
			return buf.String(), nil
		}
	}
}

type Request struct {
	Method string
	Header Header
	URL    string
	Body   io.Reader
}

var httpHandle = map[string]func(*Request, net.Conn){
	"/": func(request *Request, conn net.Conn) {
		cookie := request.Header.Get("Cookie")
		if strings.Contains(cookie, "username=") && strings.Contains(cookie, "password=") {
			_, _ = conn.Write([]byte("HTTP/1.1 302 Found\r\n" +
				"Content-Length: 0\r\n" +
				"Connection: close\r\n" +
				"Cache-Control: no-cache\r\n" +
				"Location: /login\r\n" +
				"\r\n"))
			return
		}

		_, err := conn.Write([]byte(fmt.Sprintf(page, len(root), root)))
		if err != nil {
			return
		}
	},
	"/login": func(request *Request, conn net.Conn) {
		if request.Method == "POST" {
			i, err := strconv.Atoi(request.Header.Get("Content-Length"))
			if err != nil {
				log.Print("valid Content-Length")
				return
			}

			body, err := ioutil.ReadAll(io.LimitReader(request.Body, int64(i)))
			if err != nil {
				log.Print("read body fail")
				return
			}

			v, err := url.ParseQuery(string(body))
			if err != nil {
				log.Print("parse body fail")
				return
			}

			bd := fmt.Sprintf(login, v.Get("username"), v.Get("password"))
			_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
				"Content-Type: text/html\r\n"+
				"Content-Length: %d\r\n"+
				"Connection: close\r\n"+
				"Cache-Control: no-cache\r\n"+
				fmt.Sprintf("Set-Cookie: username=%s; Max-Age=600\r\n", v.Get("username"))+
				fmt.Sprintf("Set-Cookie: password=%s; Max-Age=600\r\n", v.Get("password"))+
				"\r\n%s", len(bd), bd)))
			if err != nil {
				return
			}
		}
		if request.Method == "GET" {
			cookie := request.Header.Get("Cookie")
			if !strings.Contains(cookie, "username=") || !strings.Contains(cookie, "password=") {
				return
			}

			username := ""
			password := ""
			for _, kv := range strings.Split(cookie, ";") {
				kv = strings.Trim(kv, " ")
				spl := strings.Split(kv, "=")
				if len(spl) != 2 {
					return
				}
				switch spl[0] {
				case "username":
					username = spl[1]
				case "password":
					password = spl[1]
				}
			}

			bd := fmt.Sprintf(login, username, password)
			_, err := conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
				"Content-Type: text/html\r\n"+
				"Content-Length: %d\r\n"+
				"Connection: close\r\n"+
				"Cache-Control: no-cache\r\n"+
				"\r\n%s", len(bd), bd)))
			if err != nil {
				return
			}
		}
	},
}

const page = "HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/html\r\n" +
	"Content-Length: %d\r\n" +
	"Connection: close\r\n" +
	"Cache-Control: no-cache\r\n" +
	"\r\n%s"

const root = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Title</title>
</head>
<body>
<form action="/login" method="post">
<input placeholder="username" name="username" >
<input placeholder="password" name="password" type="password">
<input type="submit">
</form>
</body>
</html>`

const login = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Title</title>
</head>
<body>
<p>username:%s</p>
<p>password:%s</p>
</body>
</html>`

type Header map[string]string

func (h Header) Add(k, v string) {
	h[k] = v
}

func (h Header) Get(k string) string {
	return h[k]
}
