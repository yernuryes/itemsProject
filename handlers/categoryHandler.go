package handlers

import (
	"encoding/json"
	"februaryMVCProject/entities"
	"februaryMVCProject/services"
	"log"
	"net/http"
	"strconv"
)

type CategoryHandler interface {
	HandleCategoriesGet(w http.ResponseWriter, r *http.Request)
	HandleCategoriesPost(w http.ResponseWriter, r *http.Request)
	HandleCategoriesPut(w http.ResponseWriter, r *http.Request)
	HandleCategoriesDelete(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}

func (categoryHandler *categoryHandler) HandleCategoriesGet(w http.ResponseWriter, r *http.Request) {
	var idStr string = r.URL.Query().Get("id")
	if idStr == "" {
		err := json.NewEncoder(w).Encode(categoryHandler.categoryService.GetAllCategories())
		if err != nil {
			log.Fatal("Не удалось получить json всех товаров")
		}
	} else if idStr != "" {
		id, errTwo := strconv.Atoi(idStr)
		if errTwo != nil {
			log.Fatal("Ошибка конвертации")
		}
		errThree := json.NewEncoder(w).Encode(categoryHandler.categoryService.GetCategory(id))
		if errThree != nil {
			log.Fatal("Ошибка получения json")
		}
	}
}

func (categoryHandler *categoryHandler) HandleCategoriesPost(w http.ResponseWriter, r *http.Request) {
	var category entities.Category
	json.NewDecoder(r.Body).Decode(&category)
	categoryHandler.categoryService.AddCategory(category)
}

func (categoryHandler *categoryHandler) HandleCategoriesPut(w http.ResponseWriter, r *http.Request) {
	var updCategory entities.Category
	json.NewDecoder(r.Body).Decode(&updCategory)
	categoryHandler.categoryService.UpdateCategory(updCategory)
}

func (categoryHandler *categoryHandler) HandleCategoriesDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal("Ошибка парсинга")
	}
	categoryHandler.categoryService.DeleteCategory(id)
}
