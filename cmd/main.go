package main

import (
	"log"
	"net"

	"github.com/mijgona/http/pkg/server"
)


func main() {

	host := "0.0.0.0"
	port := "9999"
	srv:=server.NewServer(net.JoinHostPort(host,port))

	srv.AddPath("/","ok", "text/html")
	srv.AddPath("/about","About Alif Academy", "text/html")
	srv.AddPath("/api/payment{content}/{id}", "hello from hw 2", "text/plain")

	err := srv.Start()
	if err != nil {
		log.Print(err)
		return
	}
}


