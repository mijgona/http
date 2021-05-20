package server

import (
	"bytes"
	"io"
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

	buf := make([]byte, 8*1024)
	data := make([]byte,0)
	body := ""
	headers:=""
	requestLine := ""
	completeRequestLine := false
	complateHeader := false
	complateBody := false
	EOF := false
	for {
		if complateBody {
			break
		}

		if !EOF {
			n, err := req.Conn.Read(buf)
			if err == io.EOF {
				EOF = true
			}
			if err != nil && !EOF {
				return
			}
			data = append(data, buf[:n]...)
		}

		//Делим запрос на части
		log.Print(string(data))
		requestLineDeLim := []byte{'\r','\n'}
		requestLineEnd := bytes.Index(data, requestLineDeLim)
		if !completeRequestLine {
			requestLine = string(data[:requestLineEnd])
			data = data[requestLineEnd+2:] // delete requestLine from data. (+2 because "\r\n")
			completeRequestLine = true
		}
		
		if completeRequestLine && !complateHeader {			
			partRequestLine := bytes.Split(data, requestLineDeLim)
			if len(partRequestLine)<2 {
				continue
			}

			for i := 0; i < len(partRequestLine); i++ {
				if i>0 && len(partRequestLine[i])==0{
					complateHeader = true
					break
				}
				ind := bytes.Index(data, requestLineDeLim)
				if ind == -1 {
					log.Println("error:handle(): can't find '\\r\\n':")
					continue
				}
				headers += string(data[:ind+2])
				data = data[ind+2:]
			}
		}
		if complateHeader && !complateBody {
			if EOF {
				complateBody = true
			}
			if len(data) < 2 {
				continue
			}
			if string(data[:2]) == string(requestLineDeLim) {
				data = data[2:]
			}
			partRequestBody := bytes.Split(data, requestLineDeLim)
			for i := 0; i < len(partRequestBody); i++ {
				if i == (len(partRequestBody)-1) && !EOF {
					continue
				}
				body += string(partRequestBody[i]) + "\r\n"
			}
		}

		if requestLineEnd == -1{
			continue
		}
		
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

		//Форматируем пути
		pathParams, handlerPath:=s.convertPath(uri.Path)

		//Добавляем заголовки
		tempHead := make(map[string]string)
		header := strings.Split(headers,"\r\n")
		for i := 0; i < len(header)-1; i++ {
			header[i]=strings.Replace(header[i]," ","",-1)
			head:=strings.Split(header[i],":")
			h :=""
			for j := 1; j < len(head); j++ {
				if len(head[j])==0{
					continue
				}
				if j==len(head)-1{
					h+=head[len(head)-1]
				}else{
					h+=head[j]+":"
				}
			}
			tempHead[head[0]]=h					
						
		}
		
		log.Print("HEADERS: ",tempHead)

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
				Headers:     tempHead,
				Body:        []byte(body),
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
				Headers:     tempHead,
				Body:        []byte(body),
			}
			handler(&req)
		}

	}	
	
	
}


func (s *Server) AddPath(path string, body string, conType string) {	
	s.Register(path, func(req *Request) {
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
