package banners

import (
	"context"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/mijgona/http/pkg/types"
)

//Service представляет собой сервис по управлению баннерами
type Service struct {
	mu 		sync.RWMutex
	items	[]*types.Banner
}

//NewService создаёт сервис
func NewService() *Service {
	log.Print("Banners.NewService(): start")
	return &Service{items: make([]*types.Banner, 0)}
}

//All Возвращает все существующие баннеры
func (s *Service) All(ctx context.Context) ([]*types.Banner, error)  {
	log.Print("Banners.All(): start")
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil	
}



//ByID Возвращает баннер по ID
func (s *Service) ByID(ctx context.Context, id int64) (*types.Banner, error)  {
	log.Print("Banners.ByID(): start")
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}


//Save Возвращает сохроннённый\обновлённый баннер
func (s *Service) Save(ctx context.Context, item *types.Banner) (*types.Banner, error)  {	
	log.Print("Banners.Save(): start")
	log.Print("item: ",item)
	for i, _ := range s.items {
		if s.items[i].ID == item.ID {
			image := s.items[i].Image
			if item.Image!="" {
				image=strconv.Itoa(int(s.items[i].ID))+item.Image
			}
			s.items[i] = &types.Banner{
				ID:		s.items[i].ID,
				Title:   item.Title,
				Content: item.Content,
				Button:  item.Button,
				Link:    item.Link,
				Image: 	 image,
			}
			return s.items[i], nil
		}
	}

	if item.ID==0 {
		id := int64(1)
		if len(s.items) != 0{
			id = s.items[len(s.items)-1].ID+1
		}
		if item.Image==""{
			return nil, errors.New("banner does not contain image")
		}
		newBanner := &types.Banner{
			ID:      id,
			Title:   item.Title,
			Content: item.Content,
			Button:  item.Button,
			Link:    item.Link,
			Image: 	 strconv.Itoa(int(id))+item.Image,
		}
		s.items = append(s.items, newBanner)
		return newBanner, nil
	}
	
	

	
	return nil, errors.New("item not found")
}


//RemoveByID Возвращает удалённый баннер
func (s *Service) RemoveByID(ctx context.Context, id int64) (*types.Banner, error)  {
	log.Print("Banners.RemoveByID(): start")
	s.mu.RLock()
	defer s.mu.RUnlock()
	//находим удаляемый элемент
	RemBanner := &types.Banner{}
	for _, banner := range s.items {
		if banner.ID == id {
			RemBanner= banner
		}
	}
	//Если не нашли возвращаем ошибку
	if RemBanner.ID==0{
		return nil, errors.New("item not found")
	}

	//создаём новый слайс без удаляеимого элемента
	newItems := []*types.Banner{}
	for _, banner := range s.items {
		if banner.ID==RemBanner.ID {
			continue
		}
		newItems=append(newItems,banner)
	}
	s.items = newItems
	return RemBanner, nil
}