package handlers

import (
	"encoding/json"
	"februaryMVCProject/entities"
	"februaryMVCProject/services"
	"log"
	"net/http"
	"strconv"
)

type ItemHandler interface {
	HandleItemsGet(w http.ResponseWriter, r *http.Request)
	HandleItemsPost(w http.ResponseWriter, r *http.Request)
	HandleItemsPut(w http.ResponseWriter, r *http.Request)
	HandleItemsDelete(w http.ResponseWriter, r *http.Request)
}

type itemHandler struct {
	itemService services.ItemService
}

func NewItemHandler(itemService services.ItemService) ItemHandler {
	return &itemHandler{itemService: itemService}
}

func (itemHandler *itemHandler) HandleItemsGet(w http.ResponseWriter, r *http.Request) {
	var idStr string = r.URL.Query().Get("id")
	var name string = r.URL.Query().Get("name")
	var sort string = r.URL.Query().Get("sort")
	if idStr == "" && name == "" && sort == "" {
		err := json.NewEncoder(w).Encode(itemHandler.itemService.GetAllItems())
		if err != nil {
			log.Fatal("Не удалось получить json всех товаров")
		}
	} else if idStr != "" && name == "" && sort == "" {
		id, errTwo := strconv.Atoi(idStr)
		if errTwo != nil {
			log.Fatal("Ошибка конвертации")
		}
		errThree := json.NewEncoder(w).Encode(itemHandler.itemService.GetItem(id))
		if errThree != nil {
			log.Fatal("Ошибка получения json")
		}
	} else if idStr == "" && name != "" && sort == "" {
		err := json.NewEncoder(w).Encode(itemHandler.itemService.GetItemsByName(name))
		if err != nil {
			log.Fatal("Ошибка получения json")
		}
	} else if idStr == "" && name == "" && sort != "" {
		err := json.NewEncoder(w).Encode(itemHandler.itemService.GetItemsSortedByPrice(sort))
		if err != nil {
			log.Fatal("Ошибка получения json")
		}
	}
}

func (itemHandler *itemHandler) HandleItemsPost(w http.ResponseWriter, r *http.Request) {
	var item entities.Item
	json.NewDecoder(r.Body).Decode(&item)
	itemHandler.itemService.AddItem(item)
}

func (itemHandler *itemHandler) HandleItemsPut(w http.ResponseWriter, r *http.Request) {
	var updateItem entities.Item
	json.NewDecoder(r.Body).Decode(&updateItem)
	itemHandler.itemService.UpdateItem(updateItem)
}

func (itemHandler *itemHandler) HandleItemsDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal("Ошибка парсинга")
	}
	itemHandler.itemService.DeleteItem(id)
}
