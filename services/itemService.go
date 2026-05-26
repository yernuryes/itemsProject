package services

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
	"februaryMVCProject/mappers"
	"februaryMVCProject/repositories"
	"github.com/google/uuid"
	"log"
)

type ItemService interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []dtos.ItemDTO
	GetItem(id int) dtos.ItemDTO
	GetItemsByName(name string) []dtos.ItemDTO
	GetItemsSortedByPrice(order string) []dtos.ItemDTO
}

type itemService struct {
	itemRepository repositories.ItemRepository
}

func NewItemService(itemRepository repositories.ItemRepository) ItemService {
	return &itemService{itemRepository: itemRepository}
}

func (itemService *itemService) AddItem(item entities.Item) {
	if item.Name == "" || item.Amount == 0 || item.Price < 0 {
		log.Fatal("Данные некорректные, добавление провально")
	} else {
		item.Promocode = uuid.New().String()
		itemService.itemRepository.AddItem(item)
	}
}

func (itemService *itemService) DeleteItem(id int) {
	itemService.itemRepository.DeleteItem(id)
}

func (itemService *itemService) UpdateItem(updateItem entities.Item) {
	var currentItem entities.Item = itemService.itemRepository.GetItem(int(updateItem.ID))
	updateItem.Promocode = currentItem.Promocode
	itemService.itemRepository.UpdateItem(updateItem)
}

func (itemService *itemService) GetAllItems() []dtos.ItemDTO {
	return mappers.MapToDTOList(itemService.itemRepository.GetAllItems())
}

func (itemService *itemService) GetItem(id int) dtos.ItemDTO {
	return mappers.MapToDTO(itemService.itemRepository.GetItem(id))
}

func (itemService *itemService) GetItemsByName(name string) []dtos.ItemDTO {
	return mappers.MapToDTOList(itemService.itemRepository.GetItemsByName(name))
}
func (itemService *itemService) GetItemsSortedByPrice(order string) []dtos.ItemDTO {
	return mappers.MapToDTOList(itemService.itemRepository.GetItemsSortedByPrice(order))
}
