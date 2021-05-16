package server

import (
	"net"
	"net/url"
	"sync"
)

type HandlerFunc func(req *Request)

type Server struct {
	addr	string
	mu 		sync.RWMutex
	handlers	map[string]HandlerFunc
}

type Request struct {
	Conn 		net.Conn
	QueryParams	url.Values
	PathParams	map[string]string
}