package server

import (
	"bytes"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
)


func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc)  {
	s.mu.Lock()
	defer s.mu.Unlock()	
	log.Println("Registered new path: ", path)
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

		go s.handle(Request{
			Conn: conn,
		})
		
	}

}



func (s *Server) handle(req Request){
	var err error
	defer func ()  {
		if cerr := req.Conn.Close(); cerr != nil {
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
		n, err :=req.Conn.Read(buf)
		
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
		path, err =url.PathUnescape(path)		
		if err!=nil {
			log.Print("cant decoding path")
			continue
		}

		uri, _ :=url.ParseRequestURI(path)
		queryValues :=uri.Query()
		log.Print(uri.Query())
		
		pathParams, handlerPath:=s.convertPath(uri.Path)
		if len(pathParams)!=0{
			log.Print((pathParams))	
			s.mu.RLock()		
			handler, ok :=s.handlers[handlerPath]
			s.mu.RUnlock()
			if !ok {
				log.Print("cant find path:", uri.Path)
				return
			}
			req := Request{
				Conn:        req.Conn,
				QueryParams: queryValues,
				PathParams:  pathParams,
			}
			handler(&req)
		}else{
			log.Print((pathParams))	
			
			s.mu.RLock()		
			handler, ok :=s.handlers[handlerPath]
			s.mu.RUnlock()
			if !ok {
				log.Print("cant find path:", uri.Path)
				return
			}
			req := Request{
				Conn:        req.Conn,
				QueryParams: queryValues,
				PathParams:  pathParams,
			}
			handler(&req)
		}

	}	
	
	
}


func (s *Server) AddPath(path string, body string, conType string) {	
	s.Register(path, func(req *Request) {
	log.Println("QueryParams:", req.QueryParams)
	log.Println("PathParams:", req.PathParams)
		_, err := req.Conn.Write([]byte(
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
