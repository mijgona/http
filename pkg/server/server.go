package server

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type HandlerFunc func(conn net.Conn)

type Server struct {
	addr	string
	mu 		sync.RWMutex
	handlers	map[string]HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc)  {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path]=handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func ()  {
		if cerr := listener.Close(); cerr != nil {
			if err == nil{
				err=cerr
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error on start function:",err)
			continue
		}

		go s.handle(conn)
		
	}

}



func (s *Server) handle(conn net.Conn){
	var err error
	defer func ()  {
		if cerr := conn.Close(); cerr != nil {
			if err != nil{
				err=cerr
				return
			}
			log.Print(err)
		}
	}()

	buf := make([]byte, 4096)
	data := make([]byte,0)
	for {
		n, err :=conn.Read(buf)
		
		if err!=nil {
			return
		}

		data = append(data, buf[:n]...)
		requestLineDeLim := []byte{'\r','\n'}
		requestLineEnd := bytes.Index(data, requestLineDeLim)
		
		if requestLineEnd == -1{
			continue
		}
		
		requestLine := string(data[:requestLineEnd])
		parts := strings.Split(requestLine, " ")	
		
		if len(parts)!=3{
			continue 
		}	
		
		_, path, version := parts[0], parts[1], parts[2]
		if version[len(version)-2:] == "\r\n" {
			version = version[:len(version)-2]
		}
		if version != "HTTP/1.1" {
			log.Print("wrong version of http, should be 1.1")
			continue
		}
		
		s.mu.RLock()
		fn, ok :=s.handlers[path]
		s.mu.RUnlock()
		if !ok {
			log.Print("cant find path:", path)
			return
		}
		fn(conn)
	}	
	
	
}


func (s *Server) AddPath(path string, body string, conType string) {
	s.Register(path, func(conn net.Conn) {
		_, err := conn.Write([]byte(
			"HTTP/1.1 200 OK \r\n" +
				"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: " + conType + "\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				body,
		))
		if err != nil {
			log.Print(err)
		}
	})
}
