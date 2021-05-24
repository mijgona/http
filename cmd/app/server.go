package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	
	// banner := &types.Banner{
	// 	ID:      0,
	// 	Title:   "sdfs",
	// 	Content: "sdfsd",
	// 	Button:  "sdfsd",
	// 	Link:    "sdfsdf",
	// }
	// s.bannersSvc.Save(request.Context(), banner)
	s.mux.ServeHTTP(writer, request)
}


//Init - инициализирует сервер (регистрирует все handler-ы)
func (s *Server) Init(){
	log.Println("server.Init(): start")

	

	s.mux.HandleFunc("/banners.getAll/", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById/", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save/", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById/", s.handleRemoveById)
}

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request)  {
	log.Print("server.GetBannerByID(): start")

	idParam := request.URL.Query().Get("id")

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
	idParam := request.URL.Query().Get("id")
	titleParam := request.URL.Query().Get("title")
	contentParam := request.URL.Query().Get("content")
	buttonParam := request.URL.Query().Get("button")
	linkParam := request.URL.Query().Get("link")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item := &types.Banner{
		ID:      id,
		Title:   titleParam,
		Content: contentParam,
		Button:  buttonParam,
		Link:    linkParam,
	}
	savedItem, err := s.bannersSvc.Save(request.Context(), item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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
	idParam := request.URL.Query().Get("id")

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