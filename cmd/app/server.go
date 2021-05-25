package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mijgona/http/pkg/banners"
	"github.com/mijgona/http/pkg/types"
	// "github.com/mijgona/http/pkg/types"
)

//Представляет собой логический сервер нашего приложения
type Server struct {
	mux *http.ServeMux
	bannersSvc	*banners.Service
}

//NewServer - фунция конструктор для создания сервера
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	log.Print("server.NewServer(): start")
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	log.Print("server.ServeHTTP(): start")
	s.mux.ServeHTTP(writer, request)
}


//Init - инициализирует сервер (регистрирует все handler-ы)
func (s *Server) Init(){
	log.Println("server.Init(): start")
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveById)
}

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request)  {
	log.Print("server.GetBannerByID(): start")
	log.Print(request.Header.Get("Content-Type"))
	idParam := request.PostFormValue("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err !=nil {
		log.Print(err)
	}
}


func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request)  {	
	log.Print("server.GetAllBanners(): start")
	items, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err !=nil {
		log.Print(err)
	}

}


func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request)  {	
	log.Print("server.SaveBanner(): start")
	idParam := request.PostFormValue("id")
	titleParam := request.PostFormValue("title")
	contentParam := request.PostFormValue("content")
	buttonParam := request.PostFormValue("button")
	linkParam := request.PostFormValue("link")

	
	log.Print(request.Header.Get("Content-Type"))
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(request.Body)//тело запроса
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s", body)

	err = request.ParseMultipartForm(10*1024*1024)//изображение до 10МБ
			if err != nil {
				log.Print(err)
			}
	//импортируем из формы изображение
	file, header, err :=request.FormFile("image")	
	format := ""	
	if err !=nil && err!=http.ErrMissingFile {
		log.Print(err)
	}
	//проверяем на наличие в форме изображения
	if err!=http.ErrMissingFile{				
		i := strings.Index(header.Filename, ".")
		format = header.Filename[i:]
	}
	defer file.Close()
	//отправляем формат изображения если было загруженно изображение, и "" если изображения в форме нет
	item := &types.Banner{
		ID:      id,
		Title:   titleParam,
		Content: contentParam,
		Button:  buttonParam,
		Link:    linkParam,
		Image: 	 format,
	}
	//сохроняем
	savedItem, err := s.bannersSvc.Save(request.Context(), item)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
	log.Print("saved: ",savedItem)
	path,err := filepath.Abs("../web/banners/")
	log.Print(path)
	if err != nil {
		log.Print(err)
	}
	//если было загруженно изображение, сохроняем его
	if format !=""{
		// f, err := os.OpenFile("../../web/banners/"+savedItem.Image, os.O_CREATE|os.O_RDWR, 0777)
		dst, err := os.Create(path+"/"+savedItem.Image)
		if err != nil {
			fmt.Println(err)
		}

		io.Copy(dst, file)
		dst.Close()
	}

	data, err := json.Marshal(savedItem)
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err !=nil {
		log.Print(err)
	}
}


func (s *Server) handleRemoveById(writer http.ResponseWriter, request *http.Request)  {	
	log.Print("server.RemoveById(): start")
	idParam := request.PostFormValue("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err !=nil {
		log.Print(err)
	}
}