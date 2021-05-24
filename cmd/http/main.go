package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {

	host := "0.0.0.0"
	port := "9999"
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error) {
	srv := &http.Server{
		Addr: net.JoinHostPort(host,port),
		Handler: http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request){
			log.Print(request.RequestURI)//полный URI
			log.Print(request.Method)//method
			log.Print(request.Header)//все заголовки
			log.Print(request.Header.Get("Content-Type"))//конкретный заголовок

			log.Print(request.FormValue("tags"))//толлько первое значение Query+post
			log.Print(request.PostFormValue("tags"))//толлько первое значение post

			body, err := ioutil.ReadAll(request.Body)//тело запроса
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s", body)

			err = request.ParseMultipartForm(10*1024*1024)
			if err != nil {
				log.Print(err)
			}

			//доступно только после ParseForm либо formValue, postFormValue
			log.Print(request.Form)
			log.Print(request.PostForm)

			log.Print(request.FormFile("image"))
		}),
	}
	return srv.ListenAndServe()
}
