package main

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var gormDB *gorm.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=goDB host=db port=5432 sslmode=disable"
	var err error
	gormDB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal("connection failed")
	}
	errTwo := gormDB.AutoMigrate(&entities.Item{}, &entities.Car{}, entities.Category{}, entities.ClientSegment{})
	if errTwo != nil {
		log.Fatal("auto migration failed")
	}
}
func CloseDB() {
	s, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	errTwo := s.Close()
	if errTwo != nil {
		log.Fatal("closing failed")
	}
}

func makeMigration() {
	dsn := "postgres://postgres:admin@db:5432/goDB?sslmode=disable"
	m, err := migrate.New("file:///app/db/migrations", dsn)
	if err != nil {
		panic(err)
	}
	m.Up()
}
func backMigration() {
	dsn := "postgres://postgres:admin@db:5432/goDB?sslmode=disable"
	m, err := migrate.New("file:///app/db/migrations", dsn)
	if err != nil {
		panic(err)
	}
	m.Down()
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya
	makeMigration()

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(gormDB)
	itemRepository = itemRepo

	var categoryRepository repositories.CategoryRepository
	categoryRepo := repositories.NewCategoryRepository(gormDB)
	categoryRepository = categoryRepo

	var clientSegmentRepository repositories.ClientSegmentRepository
	clientSegmentRepo := repositories.NewClientSegmentRepository(gormDB)
	clientSegmentRepository = clientSegmentRepo

	var carRepository repositories.CarRepository
	carRepo := repositories.NewCarRepository(gormDB)
	carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	var categoryService services.CategoryService
	categoryServ := services.NewCategoryService(categoryRepository)
	categoryService = categoryServ

	var clientSegmentService services.ClientSegmentService
	clientSegmentServ := services.NewClientSegmentService(clientSegmentRepository)
	clientSegmentService = clientSegmentServ

	var carService services.CarService
	carServ := services.NewCarService(carRepository)
	carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	var categoryHandler handlers.CategoryHandler
	categoryHand := handlers.NewCategoryHandler(categoryService)
	categoryHandler = categoryHand

	var clientSegmentHandler handlers.ClientSegmentHandler
	clientSegmentHand := handlers.NewClientSegmentHandler(clientSegmentService)
	clientSegmentHandler = clientSegmentHand

	var carHandler handlers.CarHandler
	carHand := handlers.NewCarHandler(carService)
	carHandler = carHand

	router := mux.NewRouter()
	router.HandleFunc("/items", itemHandler.HandleItemsGet).Methods("GET")
	router.HandleFunc("/items", itemHandler.HandleItemsPost).Methods("POST")
	router.HandleFunc("/items", itemHandler.HandleItemsPut).Methods("PUT")
	router.HandleFunc("/items", itemHandler.HandleItemsDelete).Methods("DELETE")

	router.HandleFunc("/categories", categoryHandler.HandleCategoriesGet).Methods("GET")
	router.HandleFunc("/categories", categoryHandler.HandleCategoriesPost).Methods("POST")
	router.HandleFunc("/categories", categoryHandler.HandleCategoriesPut).Methods("PUT")
	router.HandleFunc("/categories", categoryHandler.HandleCategoriesDelete).Methods("DELETE")

	router.HandleFunc("/client-segments", clientSegmentHandler.HandleClientSegmentGet).Methods("GET")
	router.HandleFunc("/client-segments", clientSegmentHandler.HandleClientSegmentPost).Methods("POST")

	http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}

/*
Lesson 4 ManyToMany - Docker
dtos.go
package dtos

import "februaryMVCProject/entities"

type ItemDTO struct {
	ID             int64
	Name           string
	Amount         int
	Price          int
	CategoryID     int
	Category       *entities.Category
	ClientSegments []entities.ClientSegment
}

type CarDTO struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


------
entities.go
package entities

type Item struct {
	ID             int64
	Name           string
	Amount         int
	Price          int
	Promocode      string
	CategoryID     int
	Category       *Category       `gorm:"foreignKey:CategoryID"`
	ClientSegments []ClientSegment `gorm:"many2many:item_client_segments"`
}

type Category struct {
	ID   int64
	Name string
}

type ClientSegment struct {
	ID    int64
	Name  string
	Items []Item `gorm:"many2many:item_client_segments"`
}

type ItemClientSegments struct {
	ItemID          int64
	ClientSegmentID int64
}

type Car struct {
	ID        int64
	Name      string
	Model     string
	Year      int
	Promocode string
}


-----
categoryHandler.go
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



------
clientSegmentHandler.go
package handlers

import (
	"encoding/json"
	"februaryMVCProject/entities"
	"februaryMVCProject/services"
	"log"
	"net/http"
)

type ClientSegmentHandler interface {
	HandleClientSegmentGet(w http.ResponseWriter, r *http.Request)
	HandleClientSegmentPost(w http.ResponseWriter, r *http.Request)
}

type clientSegmentHandler struct {
	clientSegmentService services.ClientSegmentService
}

func NewClientSegmentHandler(clientSegmentService services.ClientSegmentService) ClientSegmentHandler {
	return &clientSegmentHandler{clientSegmentService: clientSegmentService}
}

func (clientSegmentHandler *clientSegmentHandler) HandleClientSegmentGet(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(clientSegmentHandler.clientSegmentService.GetAllClientSegments())
	if err != nil {
		log.Fatal("Не удалось получить json всех товаров")
	}

}

func (clientSegmentHandler *clientSegmentHandler) HandleClientSegmentPost(w http.ResponseWriter, r *http.Request) {
	var clientSegment entities.ClientSegment
	json.NewDecoder(r.Body).Decode(&clientSegment)
	clientSegmentHandler.clientSegmentService.AddClientSegment(clientSegment)
}


------
itemHandler.go
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




-----
main.go
package main

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var gormDB *gorm.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	gormDB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal("connection failed")
	}
	errTwo := gormDB.AutoMigrate(&entities.Item{}, &entities.Car{}, entities.Category{}, entities.ClientSegment{})
	if errTwo != nil {
		log.Fatal("auto migration failed")
	}
}
func CloseDB() {
	s, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	errTwo := s.Close()
	if errTwo != nil {
		log.Fatal("closing failed")
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(gormDB)
	itemRepository = itemRepo

	var categoryRepository repositories.CategoryRepository
	categoryRepo := repositories.NewCategoryRepository(gormDB)
	categoryRepository = categoryRepo

	var clientSegmentRepository repositories.ClientSegmentRepository
	clientSegmentRepo := repositories.NewClientSegmentRepository(gormDB)
	clientSegmentRepository = clientSegmentRepo

	var carRepository repositories.CarRepository
	carRepo := repositories.NewCarRepository(gormDB)
	carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	var categoryService services.CategoryService
	categoryServ := services.NewCategoryService(categoryRepository)
	categoryService = categoryServ

	var clientSegmentService services.ClientSegmentService
	clientSegmentServ := services.NewClientSegmentService(clientSegmentRepository)
	clientSegmentService = clientSegmentServ

	var carService services.CarService
	carServ := services.NewCarService(carRepository)
	carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	var categoryHandler handlers.CategoryHandler
	categoryHand := handlers.NewCategoryHandler(categoryService)
	categoryHandler = categoryHand

	var clientSegmentHandler handlers.ClientSegmentHandler
	clientSegmentHand := handlers.NewClientSegmentHandler(clientSegmentService)
	clientSegmentHandler = clientSegmentHand

	var carHandler handlers.CarHandler
	carHand := handlers.NewCarHandler(carService)
	carHandler = carHand

	router := mux.NewRouter()
	router.HandleFunc("/items", itemHandler.HandleItemsGet).Methods("GET")
	router.HandleFunc("/items", itemHandler.HandleItemsPost).Methods("POST")
	router.HandleFunc("/items", itemHandler.HandleItemsPut).Methods("PUT")
	router.HandleFunc("/items", itemHandler.HandleItemsDelete).Methods("DELETE")

	router.HandleFunc("/categories", categoryHandler.HandleCategoriesGet).Methods("GET")
	router.HandleFunc("/categories", categoryHandler.HandleCategoriesPost).Methods("POST")
	router.HandleFunc("/categories", categoryHandler.HandleCategoriesPut).Methods("PUT")
	router.HandleFunc("/categories", categoryHandler.HandleCategoriesDelete).Methods("DELETE")

	router.HandleFunc("/client-segments", clientSegmentHandler.HandleClientSegmentGet).Methods("GET")
	router.HandleFunc("/client-segments", clientSegmentHandler.HandleClientSegmentPost).Methods("POST")

	http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}



----
itemMapper.go
package mappers

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
)

func MapToDTO(item entities.Item) dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTO dtos.ItemDTO
	itemDTO.ID = item.ID
	itemDTO.Name = item.Name
	itemDTO.Amount = item.Amount
	itemDTO.Price = item.Price
	itemDTO.CategoryID = item.CategoryID
	itemDTO.Category = item.Category
	itemDTO.ClientSegments = item.ClientSegments
	return itemDTO
}

func MapToEntity(itemDTO dtos.ItemDTO) entities.Item { // iz itema hochu sdelat' itemDTO
	var item entities.Item
	item.ID = itemDTO.ID
	item.Name = itemDTO.Name
	item.Amount = itemDTO.Amount
	item.Price = itemDTO.Price
	item.CategoryID = itemDTO.CategoryID
	item.Category = itemDTO.Category
	item.ClientSegments = itemDTO.ClientSegments
	return item
}

func MapToDTOList(items []entities.Item) []dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTOs []dtos.ItemDTO
	for i := 0; i < len(items); i++ {
		var itemDTO dtos.ItemDTO
		itemDTO.ID = items[i].ID
		itemDTO.Name = items[i].Name
		itemDTO.Amount = items[i].Amount
		itemDTO.Price = items[i].Price
		itemDTO.CategoryID = items[i].CategoryID
		itemDTO.Category = items[i].Category
		itemDTO.ClientSegments = items[i].ClientSegments
		itemDTOs = append(itemDTOs, itemDTO)
	}
	return itemDTOs
}



------
categoryRepository.go
package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	AddCategory(category entities.Category)
	DeleteCategory(id int)
	UpdateCategory(updCategory entities.Category)
	GetAllCategories() []entities.Category
	GetCategory(id int) entities.Category
}

type categoryRepository struct {
	gormDB *gorm.DB
}

func NewCategoryRepository(gormDB *gorm.DB) CategoryRepository {
	return &categoryRepository{gormDB: gormDB}
}

func (categoryRepository *categoryRepository) AddCategory(category entities.Category) {
	categoryRepository.gormDB.Create(&category)
}

func (categoryRepository *categoryRepository) DeleteCategory(id int) {
	categoryRepository.gormDB.Delete(&entities.Category{}, id)
}

func (categoryRepository *categoryRepository) UpdateCategory(updCategory entities.Category) {
	categoryRepository.gormDB.Save(&updCategory)
}

func (categoryRepository *categoryRepository) GetAllCategories() []entities.Category {
	var categories []entities.Category
	categoryRepository.gormDB.Find(&categories)
	return categories
}

func (categoryRepository *categoryRepository) GetCategory(id int) entities.Category {
	var category entities.Category
	categoryRepository.gormDB.First(&category, id)
	return category
}


------
clientSegmentRepository.go
package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type ClientSegmentRepository interface {
	AddClientSegment(clientSegment entities.ClientSegment)
	GetAllClientSegments() []entities.ClientSegment
}

type clientSegmentRepository struct {
	gormDB *gorm.DB
}

func NewClientSegmentRepository(gormDB *gorm.DB) ClientSegmentRepository {
	return &clientSegmentRepository{gormDB: gormDB}
}

func (clientSegmentRepository *clientSegmentRepository) AddClientSegment(clientSegment entities.ClientSegment) {
	clientSegmentRepository.gormDB.Create(&clientSegment)
}

func (clientSegmentRepository *clientSegmentRepository) GetAllClientSegments() []entities.ClientSegment {
	var clientSegments []entities.ClientSegment
	clientSegmentRepository.gormDB.Find(&clientSegments)
	return clientSegments
}



------
itemRepository.go
package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
	GetItemsByName(name string) []entities.Item
	GetItemsSortedByPrice(order string) []entities.Item
}

type itemRepository struct {
	gormDB *gorm.DB
}

func NewItemRepository(gormDB *gorm.DB) ItemRepository {
	return &itemRepository{gormDB: gormDB}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	itemRepository.gormDB.Create(&item)
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	itemRepository.gormDB.Delete(&entities.Item{}, id)
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	itemRepository.gormDB.Save(&updateItem)
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Preload("ClientSegments").Preload("Category").Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item
	itemRepository.gormDB.Preload("Category").First(&item, id)
	return item
}

func (itemRepository *itemRepository) GetItemsByName(name string) []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Preload("Category").Where("name = ?", name).Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItemsSortedByPrice(order string) []entities.Item {
	var items []entities.Item
	if order == "asc" {
		itemRepository.gormDB.Preload("Category").Order("price asc").Find(&items)
	} else if order == "desc" {
		itemRepository.gormDB.Preload("Category").Order("price desc").Find(&items)
	}
	return items
}



----
categoryService.go
package services

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/repositories"
)

type CategoryService interface {
	AddCategory(category entities.Category)
	DeleteCategory(id int)
	UpdateCategory(updCategory entities.Category)
	GetAllCategories() []entities.Category
	GetCategory(id int) entities.Category
}

type categoryService struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryService(categoryRepository repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (categoryService *categoryService) AddCategory(category entities.Category) {
	categoryService.categoryRepository.AddCategory(category)
}

func (categoryService *categoryService) DeleteCategory(id int) {
	categoryService.categoryRepository.DeleteCategory(id)
}

func (categoryService *categoryService) UpdateCategory(updCategory entities.Category) {
	categoryService.categoryRepository.UpdateCategory(updCategory)
}

func (categoryService *categoryService) GetAllCategories() []entities.Category {
	return categoryService.categoryRepository.GetAllCategories()
}

func (categoryService *categoryService) GetCategory(id int) entities.Category {
	return categoryService.categoryRepository.GetCategory(id)
}




-----
clientSegmentService.go
package services

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/repositories"
)

type ClientSegmentService interface {
	AddClientSegment(clientSegment entities.ClientSegment)
	GetAllClientSegments() []entities.ClientSegment
}

type clientSegmentService struct {
	clientSegmentRepository repositories.ClientSegmentRepository
}

func NewClientSegmentService(clientSegmentRepository repositories.ClientSegmentRepository) ClientSegmentService {
	return &clientSegmentService{clientSegmentRepository: clientSegmentRepository}
}

func (clientSegmentService *clientSegmentService) AddClientSegment(clientSegment entities.ClientSegment) {
	clientSegmentService.clientSegmentRepository.AddClientSegment(clientSegment)
}

func (clientSegmentService *clientSegmentService) GetAllClientSegments() []entities.ClientSegment {
	return clientSegmentService.clientSegmentRepository.GetAllClientSegments()
}




------
itemService.go
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


*/

/*
Lesson 3 ManyToOne связка таблиц

main.go
package main

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var gormDB *gorm.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	gormDB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal("connection failed")
	}
	errTwo := gormDB.AutoMigrate(&entities.Item{}, &entities.Car{}, entities.Category{})
	if errTwo != nil {
		log.Fatal("auto migration failed")
	}
}
func CloseDB() {
	s, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	errTwo := s.Close()
	if errTwo != nil {
		log.Fatal("closing failed")
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(gormDB)
	itemRepository = itemRepo

	var categoryRepository repositories.CategoryRepository
	categoryRepo := repositories.NewCategoryRepository(gormDB)
	categoryRepository = categoryRepo

	var carRepository repositories.CarRepository
	carRepo := repositories.NewCarRepository(gormDB)
	carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	var categoryService services.CategoryService
	categoryServ := services.CategoryService(categoryRepository)
	categoryService = categoryServ

	var carService services.CarService
	carServ := services.NewCarService(carRepository)
	carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	var categoryHandler handlers.CategoryHandler
	categoryHand := handlers.NewCategoryHandler(categoryService)
	categoryHandler = categoryHand

	var carHandler handlers.CarHandler
	carHand := handlers.NewCarHandler(carService)
	carHandler = carHand

	router := mux.NewRouter()
	router.HandleFunc("/items", itemHandler.HandleItemsGet).Methods("GET")
	router.HandleFunc("/items", itemHandler.HandleItemsPost).Methods("POST")
	router.HandleFunc("/items", itemHandler.HandleItemsPut).Methods("PUT")
	router.HandleFunc("/items", itemHandler.HandleItemsDelete).Methods("DELETE")

	router.HandleFunc("/categories", categoryHandler.HandleCategoriesGet).Methods("GET")
	router.HandleFunc("/categories", categoryHandler.HandleCategoriesPost).Methods("POST")
	router.HandleFunc("categories", categoryHandler.HandleCategoriesPut).Methods("PUT")
	router.HandleFunc("categories", categoryHandler.HandleCategoriesDelete).Methods("DELETE")

	http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}



-----

dtos.go

package dtos

import "februaryMVCProject/entities"

type ItemDTO struct {
	ID         int64
	Name       string
	Amount     int
	Price      int
	CategoryID int
	Category   *entities.Category
}

type CarDTO struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


---
entities.go
package entities

type Item struct {
	ID         int64
	Name       string
	Amount     int
	Price      int
	Promocode  string
	CategoryID int
	Category   *Category `gorm:"foreignKey : CategoryID"`
}

type Category struct {
	ID   int64
	Name string
}

type Car struct {
	ID        int64
	Name      string
	Model     string
	Year      int
	Promocode string
}


----
categoryHandler.go
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



----
itemHandler.go
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



itemMapper.go
package mappers

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
)

func MapToDTO(item entities.Item) dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTO dtos.ItemDTO
	itemDTO.ID = item.ID
	itemDTO.Name = item.Name
	itemDTO.Amount = item.Amount
	itemDTO.Price = item.Price
	itemDTO.CategoryID = item.CategoryID
	itemDTO.Category = item.Category
	return itemDTO
}

func MapToEntity(itemDTO dtos.ItemDTO) entities.Item { // iz itema hochu sdelat' itemDTO
	var item entities.Item
	item.ID = itemDTO.ID
	item.Name = itemDTO.Name
	item.Amount = itemDTO.Amount
	item.Price = itemDTO.Price
	item.CategoryID = itemDTO.CategoryID
	item.Category = itemDTO.Category
	return item
}

func MapToDTOList(items []entities.Item) []dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTOs []dtos.ItemDTO
	for i := 0; i < len(items); i++ {
		var itemDTO dtos.ItemDTO
		itemDTO.ID = items[i].ID
		itemDTO.Name = items[i].Name
		itemDTO.Amount = items[i].Amount
		itemDTO.Price = items[i].Price
		itemDTO.CategoryID = items[i].CategoryID
		itemDTO.Category = items[i].Category
		itemDTOs = append(itemDTOs, itemDTO)
	}
	return itemDTOs
}


categoryRepository.go
package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	AddCategory(category entities.Category)
	DeleteCategory(id int)
	UpdateCategory(updCategory entities.Category)
	GetAllCategories() []entities.Category
	GetCategory(id int) entities.Category
}

type categoryRepository struct {
	gormDB *gorm.DB
}

func NewCategoryRepository(gormDB *gorm.DB) CategoryRepository {
	return &categoryRepository{gormDB: gormDB}
}

func (categoryRepository *categoryRepository) AddCategory(category entities.Category) {
	categoryRepository.gormDB.Create(&category)
}

func (categoryRepository *categoryRepository) DeleteCategory(id int) {
	categoryRepository.gormDB.Delete(&entities.Category{}, id)
}

func (categoryRepository *categoryRepository) UpdateCategory(updCategory entities.Category) {
	categoryRepository.gormDB.Save(&updCategory)
}

func (categoryRepository *categoryRepository) GetAllCategories() []entities.Category {
	var categories []entities.Category
	categoryRepository.gormDB.Find(&categories)
	return categories
}

func (categoryRepository *categoryRepository) GetCategory(id int) entities.Category {
	var category entities.Category
	categoryRepository.gormDB.First(&category, id)
	return category
}


-----
itemRepository.go
package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
	GetItemsByName(name string) []entities.Item
	GetItemsSortedByPrice(order string) []entities.Item
}

type itemRepository struct {
	gormDB *gorm.DB
}

func NewItemRepository(gormDB *gorm.DB) ItemRepository {
	return &itemRepository{gormDB: gormDB}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	itemRepository.gormDB.Create(&item)
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	itemRepository.gormDB.Delete(&entities.Item{}, id)
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	itemRepository.gormDB.Save(&updateItem)
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Preload("Category").Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item
	itemRepository.gormDB.Preload("Category").First(&item, id)
	return item
}

func (itemRepository *itemRepository) GetItemsByName(name string) []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Preload("Category").Where("name = ?", name).Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItemsSortedByPrice(order string) []entities.Item {
	var items []entities.Item
	if order == "asc" {
		itemRepository.gormDB.Preload("Category").Order("price asc").Find(&items)
	} else if order == "desc" {
		itemRepository.gormDB.Preload("Category").Order("price desc").Find(&items)
	}
	return items
}



----
categoryService.go

package services

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/repositories"
)

type CategoryService interface {
	AddCategory(category entities.Category)
	DeleteCategory(id int)
	UpdateCategory(updCategory entities.Category)
	GetAllCategories() []entities.Category
	GetCategory(id int) entities.Category
}

type categoryService struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryService(categoryRepository repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (categoryService *categoryService) AddCategory(category entities.Category) {
	categoryService.categoryRepository.AddCategory(category)
}

func (categoryService *categoryService) DeleteCategory(id int) {
	categoryService.categoryRepository.DeleteCategory(id)
}

func (categoryService *categoryService) UpdateCategory(updCategory entities.Category) {
	categoryService.categoryRepository.UpdateCategory(updCategory)
}

func (categoryService *categoryService) GetAllCategories() []entities.Category {
	return categoryService.categoryRepository.GetAllCategories()
}

func (categoryService *categoryService) GetCategory(id int) entities.Category {
	return categoryService.categoryRepository.GetCategory(id)
}



----
itemService.go
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

*/

/*
Lesson 2 Gorilla Gorm Part 2


Do Gorilly kod


main.go

package main

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var gormDB *gorm.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	gormDB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal("connection failed")
	}
	errTwo := gormDB.AutoMigrate(&entities.Item{}, &entities.Car{})
	if errTwo != nil {
		log.Fatal("auto migration failed")
	}
}
func CloseDB() {
	s, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	errTwo := s.Close()
	if errTwo != nil {
		log.Fatal("closing failed")
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(gormDB)
	itemRepository = itemRepo

	var carRepository repositories.CarRepository
	carRepo := repositories.NewCarRepository(gormDB)
	carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	var carService services.CarService
	carServ := services.NewCarService(carRepository)
	carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	var carHandler handlers.CarHandler
	carHand := handlers.NewCarHandler(carService)
	carHandler = carHand

	http.HandleFunc("/items", itemHandler.HandleItems)

	http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}


----

dtos.go


package dtos

type ItemDTO struct {
	ID     int64
	Name   string
	Amount int
	Price  int
}

type CarDTO struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


---

entities.go

package entities

type Item struct {
	ID        int64
	Name      string
	Amount    int
	Price     int
	Promocode string
}

type Car struct {
	ID        int64
	Name      string
	Model     string
	Year      int
	Promocode string
}


----


itemMapper.go


package mappers

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
)

func MapToDTO(item entities.Item) dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTO dtos.ItemDTO
	itemDTO.ID = item.ID
	itemDTO.Name = item.Name
	itemDTO.Amount = item.Amount
	itemDTO.Price = item.Price
	return itemDTO
}

func MapToEntity(itemDTO dtos.ItemDTO) entities.Item { // iz itema hochu sdelat' itemDTO
	var item entities.Item
	item.ID = itemDTO.ID
	item.Name = itemDTO.Name
	item.Amount = itemDTO.Amount
	item.Price = itemDTO.Price
	return item
}

func MapToDTOList(items []entities.Item) []dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTOs []dtos.ItemDTO
	for i := 0; i < len(items); i++ {
		var itemDTO dtos.ItemDTO
		itemDTO.ID = items[i].ID
		itemDTO.Name = items[i].Name
		itemDTO.Amount = items[i].Amount
		itemDTO.Price = items[i].Price
		itemDTOs = append(itemDTOs, itemDTO)
	}
	return itemDTOs
}



-----


itemRepository.go


package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
	GetItemsByName(name string) []entities.Item
	GetItemsSortedByPrice(order string) []entities.Item
}

type itemRepository struct {
	gormDB *gorm.DB
}

func NewItemRepository(gormDB *gorm.DB) ItemRepository {
	return &itemRepository{gormDB: gormDB}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	itemRepository.gormDB.Create(&item)
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	itemRepository.gormDB.Delete(&entities.Item{}, id)
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	itemRepository.gormDB.Save(&updateItem)
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item
	itemRepository.gormDB.First(&item, id)
	return item
}

func (itemRepository *itemRepository) GetItemsByName(name string) []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Where("name = ?", name).Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItemsSortedByPrice(order string) []entities.Item {
	var items []entities.Item
	if order == "asc" {
		itemRepository.gormDB.Order("price asc").Find(&items)
	} else if order == "desc" {
		itemRepository.gormDB.Order("price desc").Find(&items)
	}
	return items
}


----

itemService.go


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



----

itemHandler.go


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
	HandleItems(w http.ResponseWriter, r *http.Request)
}

type itemHandler struct {
	itemService services.ItemService
}

func NewItemHandler(itemService services.ItemService) ItemHandler {
	return &itemHandler{itemService: itemService}
}

func (itemHandler *itemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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

	} else if r.Method == http.MethodPost {
		var item entities.Item
		json.NewDecoder(r.Body).Decode(&item)
		itemHandler.itemService.AddItem(item)
	} else if r.Method == http.MethodPut {
		var updateItem entities.Item
		json.NewDecoder(r.Body).Decode(&updateItem)
		itemHandler.itemService.UpdateItem(updateItem)
	} else if r.Method == http.MethodDelete {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal("Ошибка парсинга")
		}
		itemHandler.itemService.DeleteItem(id)
	}
}


KOD Posle Rouming

dtos.go
package dtos

type ItemDTO struct {
	ID     int64
	Name   string
	Amount int
	Price  int
}

type CarDTO struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


----

entities.go
package entities

type Item struct {
	ID        int64
	Name      string
	Amount    int
	Price     int
	Promocode string
}

type Car struct {
	ID        int64
	Name      string
	Model     string
	Year      int
	Promocode string
}



----
main.go

package main

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var gormDB *gorm.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	gormDB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal("connection failed")
	}
	errTwo := gormDB.AutoMigrate(&entities.Item{}, &entities.Car{})
	if errTwo != nil {
		log.Fatal("auto migration failed")
	}
}
func CloseDB() {
	s, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	errTwo := s.Close()
	if errTwo != nil {
		log.Fatal("closing failed")
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(gormDB)
	itemRepository = itemRepo

	var carRepository repositories.CarRepository
	carRepo := repositories.NewCarRepository(gormDB)
	carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	var carService services.CarService
	carServ := services.NewCarService(carRepository)
	carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	var carHandler handlers.CarHandler
	carHand := handlers.NewCarHandler(carService)
	carHandler = carHand

	router := mux.NewRouter()
	router.HandleFunc("/items", itemHandler.HandleItemsGet).Methods("GET")
	router.HandleFunc("/items", itemHandler.HandleItemsPost).Methods("POST")
	router.HandleFunc("/items", itemHandler.HandleItemsPut).Methods("PUT")
	router.HandleFunc("/items", itemHandler.HandleItemsDelete).Methods("DELETE")

	http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}


----
itemHandler.go


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



----
itemMapper.go

package mappers

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
)

func MapToDTO(item entities.Item) dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTO dtos.ItemDTO
	itemDTO.ID = item.ID
	itemDTO.Name = item.Name
	itemDTO.Amount = item.Amount
	itemDTO.Price = item.Price
	return itemDTO
}

func MapToEntity(itemDTO dtos.ItemDTO) entities.Item { // iz itema hochu sdelat' itemDTO
	var item entities.Item
	item.ID = itemDTO.ID
	item.Name = itemDTO.Name
	item.Amount = itemDTO.Amount
	item.Price = itemDTO.Price
	return item
}

func MapToDTOList(items []entities.Item) []dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTOs []dtos.ItemDTO
	for i := 0; i < len(items); i++ {
		var itemDTO dtos.ItemDTO
		itemDTO.ID = items[i].ID
		itemDTO.Name = items[i].Name
		itemDTO.Amount = items[i].Amount
		itemDTO.Price = items[i].Price
		itemDTOs = append(itemDTOs, itemDTO)
	}
	return itemDTOs
}


-----
itemRepository.go
package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
	GetItemsByName(name string) []entities.Item
	GetItemsSortedByPrice(order string) []entities.Item
}

type itemRepository struct {
	gormDB *gorm.DB
}

func NewItemRepository(gormDB *gorm.DB) ItemRepository {
	return &itemRepository{gormDB: gormDB}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	itemRepository.gormDB.Create(&item)
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	itemRepository.gormDB.Delete(&entities.Item{}, id)
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	itemRepository.gormDB.Save(&updateItem)
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item
	itemRepository.gormDB.First(&item, id)
	return item
}

func (itemRepository *itemRepository) GetItemsByName(name string) []entities.Item {
	var items []entities.Item
	itemRepository.gormDB.Where("name = ?", name).Find(&items)
	return items
}

func (itemRepository *itemRepository) GetItemsSortedByPrice(order string) []entities.Item {
	var items []entities.Item
	if order == "asc" {
		itemRepository.gormDB.Order("price asc").Find(&items)
	} else if order == "desc" {
		itemRepository.gormDB.Order("price desc").Find(&items)
	}
	return items
}




-----
itemService.go
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


*/

/*
Lesson 1 Gorilla Gorm Intro

package main

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var gormDB *gorm.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	gormDB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal("connection failed")
	}
	errTwo := gormDB.AutoMigrate(&entities.Item{}, &entities.Car{})
	if errTwo != nil {
		log.Fatal("auto migration failed")
	}
}
func CloseDB() {
	s, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	errTwo := s.Close()
	if errTwo != nil {
		log.Fatal("closing failed")
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(gormDB)
	itemRepository = itemRepo

	//var carRepository repositories.CarRepositories
	//carRepo := repositories.NewCarRepositories(dbase)
	//carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	//var carService services.CarService
	//carServ := services.NewCarService(carRepository)
	//carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	//var carHandler handlers.CarHandler
	//carHand := handlers.NewCarRepository(carService)
	//carHandler = carHand

	http.HandleFunc("/items", itemHandler.HandleItems)

	//http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}

entities.go

package entities

type Item struct {
	ID        int64
	Name      string
	Amount    int
	Price     int
	Promocode string
}

type Car struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


dtos.go

package dtos

type ItemDTO struct {
	ID     int64
	Name   string
	Amount int
	Price  int
}


itemHandler.go
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
	HandleItems(w http.ResponseWriter, r *http.Request)
}

type itemHandler struct {
	itemService services.ItemService
}

func NewItemHandler(itemService services.ItemService) ItemHandler {
	return &itemHandler{itemService: itemService}
}

func (itemHandler *itemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var idStr string = r.URL.Query().Get("id")
		if idStr == "" {
			err := json.NewEncoder(w).Encode(itemHandler.itemService.GetAllItems())
			if err != nil {
				log.Fatal("Не удалось получить json всех товаров")
			}
		} else if idStr != "" {
			id, errTwo := strconv.Atoi(idStr)
			if errTwo != nil {
				log.Fatal("Ошибка конвертации")
			}
			errThree := json.NewEncoder(w).Encode(itemHandler.itemService.GetItem(id))
			if errThree != nil {
				log.Fatal("Ошибка получения json")
			}
		}
	} else if r.Method == http.MethodPost {
		var item entities.Item
		json.NewDecoder(r.Body).Decode(&item)
		itemHandler.itemService.AddItem(item)
	} else if r.Method == http.MethodPut {
		var updateItem entities.Item
		json.NewDecoder(r.Body).Decode(&updateItem)
		itemHandler.itemService.UpdateItem(updateItem)
	} else if r.Method == http.MethodDelete {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal("Ошибка парсинга")
		}
		itemHandler.itemService.DeleteItem(id)
	}
}


----

itemMapper.go

package mappers

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
)

func MapToDTO(item entities.Item) dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTO dtos.ItemDTO
	itemDTO.ID = item.ID
	itemDTO.Name = item.Name
	itemDTO.Amount = item.Amount
	itemDTO.Price = item.Price
	return itemDTO
}

func MapToEntity(itemDTO dtos.ItemDTO) entities.Item { // iz itema hochu sdelat' itemDTO
	var item entities.Item
	item.ID = itemDTO.ID
	item.Name = itemDTO.Name
	item.Amount = itemDTO.Amount
	item.Price = itemDTO.Price
	return item
}

func MapToDTOList(items []entities.Item) []dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTOs []dtos.ItemDTO
	for i := 0; i < len(items); i++ {
		var itemDTO dtos.ItemDTO
		itemDTO.ID = items[i].ID
		itemDTO.Name = items[i].Name
		itemDTO.Amount = items[i].Amount
		itemDTO.Price = items[i].Price
		itemDTOs = append(itemDTOs, itemDTO)
	}
	return itemDTOs
}


----

itemRepository.go

package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
}

type itemRepository struct {
	gormDB *gorm.DB
}

func NewItemRepository(gormDB *gorm.DB) ItemRepository {
	return &itemRepository{gormDB: gormDB}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	itemRepository.gormDB.Create(&item)
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	itemRepository.gormDB.Delete(&entities.Item{}, id)
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	itemRepository.gormDB.Save(&updateItem)
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item

	itemRepository.gormDB.Find(&items)

	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item
	itemRepository.gormDB.First(&item, id)
	return item
}

itemService.go

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

*/

/*
Lesson 8 DTO


main.go

package main

import (
	"database/sql"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"fmt"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var dbase *sql.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	dbase, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	err = dbase.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to database")
}
func CloseDB() {
	err := dbase.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(dbase)
	itemRepository = itemRepo

	var carRepository repositories.CarRepositories
	carRepo := repositories.NewCarRepositories(dbase)
	carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	var carService services.CarService
	carServ := services.NewCarService(carRepository)
	carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	var carHandler handlers.CarHandler
	carHand := handlers.NewCarRepository(carService)
	carHandler = carHand

	http.HandleFunc("/items", itemHandler.HandleItems)

	http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}


----

entities.go

package entities

type Item struct {
	ID        int64
	Name      string
	Amount    int
	Price     int
	Promocode string
}

type Car struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


----

dtos.go

package dtos

type ItemDTO struct {
	ID     int64
	Name   string
	Amount int
	Price  int
}


-----

itemHandler.go

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
	HandleItems(w http.ResponseWriter, r *http.Request)
}

type itemHandler struct {
	itemService services.ItemService
}

func NewItemHandler(itemService services.ItemService) ItemHandler {
	return &itemHandler{itemService: itemService}
}

func (itemHandler *itemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var idStr string = r.URL.Query().Get("id")
		if idStr == "" {
			err := json.NewEncoder(w).Encode(itemHandler.itemService.GetAllItems())
			if err != nil {
				log.Fatal("Не удалось получить json всех товаров")
			}
		} else if idStr != "" {
			id, errTwo := strconv.Atoi(idStr)
			if errTwo != nil {
				log.Fatal("Ошибка конвертации")
			}
			errThree := json.NewEncoder(w).Encode(itemHandler.itemService.GetItem(id))
			if errThree != nil {
				log.Fatal("Ошибка получения json")
			}
		}
	} else if r.Method == http.MethodPost {
		var item entities.Item
		json.NewDecoder(r.Body).Decode(&item)
		itemHandler.itemService.AddItem(item)
	} else if r.Method == http.MethodPut {
		var updateItem entities.Item
		json.NewDecoder(r.Body).Decode(&updateItem)
		itemHandler.itemService.UpdateItem(updateItem)
	} else if r.Method == http.MethodDelete {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal("Ошибка парсинга")
		}
		itemHandler.itemService.DeleteItem(id)
	}
}


----

itemRepository.go


package repositories

import (
	"database/sql"
	"februaryMVCProject/entities"
	"log"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
}

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	_, err := itemRepository.db.Exec("insert into items(name, amount, price, promocode) values($1,$2,$3,$4)", item.Name, item.Amount, item.Price, item.Promocode)
	if err != nil {
		log.Fatal("Добавление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	_, err := itemRepository.db.Exec("delete from items where id = $1", id)
	if err != nil {
		log.Fatal("Удаление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	_, err := itemRepository.db.Exec("update items set name = $1, amount = $2, price = $3 where id = $4", updateItem.Name, updateItem.Amount, updateItem.Price, updateItem.ID)
	if err != nil {
		log.Fatal("Обновление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item

	rows, err := itemRepository.db.Query("select * from items")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Price, &item.Promocode) //как fmt.Scan()
		if err != nil {
			log.Fatal("Ошибка конвертации", err)
		}
		items = append(items, item)
	}
	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item

	row := itemRepository.db.QueryRow("select * from items where id = $1", id)
	if row == nil {
		log.Fatal("Row is empty")
	}
	err := row.Scan(&item.ID, &item.Name, &item.Amount, &item.Price, &item.Promocode) //как fmt.Scan()
	if err != nil {
		log.Fatal("Ошибка конвертации", err)
	}
	return item
}



----

itemService.go

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
	itemService.itemRepository.UpdateItem(updateItem)
}

func (itemService *itemService) GetAllItems() []dtos.ItemDTO {
	return mappers.MapToDTOList(itemService.itemRepository.GetAllItems())
}

func (itemService *itemService) GetItem(id int) dtos.ItemDTO {
	return mappers.MapToDTO(itemService.itemRepository.GetItem(id))
}


----

itemMapper.go
package mappers

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
)

func MapToDTO(item entities.Item) dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTO dtos.ItemDTO
	itemDTO.ID = item.ID
	itemDTO.Name = item.Name
	itemDTO.Amount = item.Amount
	itemDTO.Price = item.Price
	return itemDTO
}

func MapToEntity(itemDTO dtos.ItemDTO) entities.Item { // iz itema hochu sdelat' itemDTO
	var item entities.Item
	item.ID = itemDTO.ID
	item.Name = itemDTO.Name
	item.Amount = itemDTO.Amount
	item.Price = itemDTO.Price
	return item
}

func MapToDTOList(items []entities.Item) []dtos.ItemDTO { // iz itema hochu sdelat' itemDTO
	var itemDTOs []dtos.ItemDTO
	for i := 0; i < len(items); i++ {
		var itemDTO dtos.ItemDTO
		itemDTO.ID = items[i].ID
		itemDTO.Name = items[i].Name
		itemDTO.Amount = items[i].Amount
		itemDTO.Price = items[i].Price
		itemDTOs = append(itemDTOs, itemDTO)
	}
	return itemDTOs
}

*/

/*
Lesson 7 Handlers && Services

main.go

package main

import (
	"database/sql"
	"februaryMVCProject/handlers"
	"februaryMVCProject/repositories"
	"februaryMVCProject/services"
	"fmt"
	"log"
	"net/http"
)
import _ "github.com/lib/pq"

var dbase *sql.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	dbase, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	err = dbase.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to database")
}
func CloseDB() {
	err := dbase.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(dbase)
	itemRepository = itemRepo

	var carRepository repositories.CarRepositories
	carRepo := repositories.NewCarRepositories(dbase)
	carRepository = carRepo

	var itemService services.ItemService
	itemServ := services.NewItemService(itemRepository)
	itemService = itemServ

	var carService services.CarService
	carServ := services.NewCarService(carRepository)
	carService = carServ

	var itemHandler handlers.ItemHandler
	itemHand := handlers.NewItemHandler(itemService)
	itemHandler = itemHand

	var carHandler handlers.CarHandler
	carHand := handlers.NewCarRepository(carService)
	carHandler = carHand

	http.HandleFunc("/items", itemHandler.HandleItems)

	http.HandleFunc("/cars", carHandler.HandleCars)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}

-----

itemRepository.go
package repositories

import (
	"database/sql"
	"februaryMVCProject/entities"
	"log"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
}

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	_, err := itemRepository.db.Exec("insert into items(name, amount, price) values($1,$2,$3)", item.Name, item.Amount, item.Price)
	if err != nil {
		log.Fatal("Добавление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	_, err := itemRepository.db.Exec("delete from items where id = $1", id)
	if err != nil {
		log.Fatal("Удаление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	_, err := itemRepository.db.Exec("update items set name = $1, amount = $2, price = $3 where id = $4", updateItem.Name, updateItem.Amount, updateItem.Price, updateItem.ID)
	if err != nil {
		log.Fatal("Обновление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item

	rows, err := itemRepository.db.Query("select * from items")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Price) //как fmt.Scan()
		if err != nil {
			log.Fatal("Ошибка конвертации", err)
		}
		items = append(items, item)
	}
	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item

	row := itemRepository.db.QueryRow("select * from items where id = $1", id)
	if row == nil {
		log.Fatal("Row is empty")
	}
	err := row.Scan(&item.ID, &item.Name, &item.Amount, &item.Price) //как fmt.Scan()
	if err != nil {
		log.Fatal("Ошибка конвертации", err)
	}
	return item
}



----

carRepository.go

package repositories

import (
	"database/sql"
	"februaryMVCProject/entities"
	"log"
)

type CarRepositories interface {
	AddItem(car entities.Car)
}

type carRepositories struct {
	db *sql.DB
}

func NewCarRepositories(db *sql.DB) CarRepositories {
	return &carRepositories{db: db}
}

func (carRepositories *carRepositories) AddItem(car entities.Car) {
	_, err := carRepositories.db.Exec("insert into cars(name, model, year) values($1, $2, $3)", car.Name, car.Model, car.Year)
	if err != nil {
		log.Fatal("Добавление безуспешно", err)
	}
}


----

itemHandler.go

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
	HandleItems(w http.ResponseWriter, r *http.Request)
}

type itemHandler struct {
	itemService services.ItemService
}

func NewItemHandler(itemService services.ItemService) ItemHandler {
	return &itemHandler{itemService: itemService}
}

func (itemHandler *itemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var idStr string = r.URL.Query().Get("id")
		if idStr == "" {
			err := json.NewEncoder(w).Encode(itemHandler.itemService.GetAllItems())
			if err != nil {
				log.Fatal("Не удалось получить json всех товаров")
			}
		} else if idStr != "" {
			id, errTwo := strconv.Atoi(idStr)
			if errTwo != nil {
				log.Fatal("Ошибка конвертации")
			}
			errThree := json.NewEncoder(w).Encode(itemHandler.itemService.GetItem(id))
			if errThree != nil {
				log.Fatal("Ошибка получения json")
			}
		}
	} else if r.Method == http.MethodPost {
		var item entities.Item
		json.NewDecoder(r.Body).Decode(&item)
		itemHandler.itemService.AddItem(item)
	} else if r.Method == http.MethodPut {
		var updateItem entities.Item
		json.NewDecoder(r.Body).Decode(&updateItem)
		itemHandler.itemService.UpdateItem(updateItem)
	} else if r.Method == http.MethodDelete {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal("Ошибка парсинга")
		}
		itemHandler.itemService.DeleteItem(id)
	}
}


-----

carHandler.go

package handlers

import (
	"encoding/json"
	"februaryMVCProject/entities"
	"februaryMVCProject/services"
	"net/http"
)

type CarHandler interface {
	HandleCars(w http.ResponseWriter, r *http.Request)
}

type carHandler struct {
	carService services.CarService
}

func NewCarRepository(carService services.CarService) CarHandler {
	return &carHandler{carService: carService}
}

func (carHandler *carHandler) HandleCars(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var car entities.Car
		json.NewDecoder(r.Body).Decode(&car)
		carHandler.carService.AddItem(car)
	}
}


-----

itemService

package services

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/repositories"
	"log"
)

type ItemService interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
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
		itemService.itemRepository.AddItem(item)
	}
}

func (itemService *itemService) DeleteItem(id int) {
	itemService.itemRepository.DeleteItem(id)
}

func (itemService *itemService) UpdateItem(updateItem entities.Item) {
	itemService.itemRepository.UpdateItem(updateItem)
}

func (itemService *itemService) GetAllItems() []entities.Item {
	return itemService.itemRepository.GetAllItems()
}

func (itemService *itemService) GetItem(id int) entities.Item {
	return itemService.itemRepository.GetItem(id)
}


-----

carService.go

package services

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/repositories"
	"log"
)

type CarService interface {
	AddItem(car entities.Car)
}

type carService struct {
	carRepository repositories.CarRepositories
}

func NewCarService(carRepository repositories.CarRepositories) CarService {
	return &carService{carRepository: carRepository}
}

func (carService *carService) AddItem(car entities.Car) {
	if car.Name == "" {
		log.Fatal("Без названии нельзя добавить")
	} else {
		carService.carRepository.AddItem(car)
	}
}


----

entities.go

package entities

type Item struct {
	ID     int64
	Name   string
	Amount int
	Price  int
}

type Car struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


*/

/*
Lesson 6 Repositories

main.go
package main

import (
	"database/sql"
	"encoding/json"
	"februaryMVCProject/entities"
	"februaryMVCProject/repositories"
	"fmt"
	"log"
	"net/http"
	"strconv"
)
import _ "github.com/lib/pq"

var dbase *sql.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	dbase, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	err = dbase.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to database")
}
func CloseDB() {
	err := dbase.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitDB()
	defer CloseDB() //samym poslednim vypolnitsya

	var itemRepository repositories.ItemRepository
	itemRepo := repositories.NewItemRepository(dbase)
	itemRepository = itemRepo

	var carRepository repositories.CarRepositories
	carRepo := repositories.NewCarRepositories(dbase)
	carRepository = carRepo

	http.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var car entities.Car
			json.NewDecoder(r.Body).Decode(&car)
			carRepository.AddItem(car)
		}
	})

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			var idStr string = r.URL.Query().Get("id")
			if idStr == "" {
				err := json.NewEncoder(w).Encode(itemRepository.GetAllItems())
				if err != nil {
					log.Fatal("Не удалось получить json всех товаров")
				}
			} else if idStr != "" {
				id, errTwo := strconv.Atoi(idStr)
				if errTwo != nil {
					log.Fatal("Ошибка конвертации")
				}
				errThree := json.NewEncoder(w).Encode(itemRepository.GetItem(id))
				if errThree != nil {
					log.Fatal("Ошибка получения json")
				}
			}
		} else if r.Method == http.MethodPost {
			var item entities.Item
			json.NewDecoder(r.Body).Decode(&item)
			itemRepository.AddItem(item)
		} else if r.Method == http.MethodPut {
			var updateItem entities.Item
			json.NewDecoder(r.Body).Decode(&updateItem)
			itemRepository.UpdateItem(updateItem)
		} else if r.Method == http.MethodDelete {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				log.Fatal("Ошибка парсинга")
			}
			itemRepository.DeleteItem(id)
		}
	})

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}


---

itemRepository.go
package repositories

import (
	"database/sql"
	"februaryMVCProject/entities"
	"log"
)

type ItemRepository interface {
	AddItem(item entities.Item)
	DeleteItem(id int)
	UpdateItem(updateItem entities.Item)
	GetAllItems() []entities.Item
	GetItem(id int) entities.Item
}

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (itemRepository *itemRepository) AddItem(item entities.Item) {
	_, err := itemRepository.db.Exec("insert into items(name, amount, price) values($1,$2,$3)", item.Name, item.Amount, item.Price)
	if err != nil {
		log.Fatal("Добавление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) DeleteItem(id int) {
	_, err := itemRepository.db.Exec("delete from items where id = $1", id)
	if err != nil {
		log.Fatal("Удаление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) UpdateItem(updateItem entities.Item) {
	_, err := itemRepository.db.Exec("update items set name = $1, amount = $2, price = $3 where id = $4", updateItem.Name, updateItem.Amount, updateItem.Price, updateItem.ID)
	if err != nil {
		log.Fatal("Обновление безуспешно,", err)
	}
}

func (itemRepository *itemRepository) GetAllItems() []entities.Item {
	var items []entities.Item

	rows, err := itemRepository.db.Query("select * from items")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Price) //как fmt.Scan()
		if err != nil {
			log.Fatal("Ошибка конвертации", err)
		}
		items = append(items, item)
	}
	return items
}

func (itemRepository *itemRepository) GetItem(id int) entities.Item {
	var item entities.Item

	row := itemRepository.db.QueryRow("select * from items where id = $1", id)
	if row == nil {
		log.Fatal("Row is empty")
	}
	err := row.Scan(&item.ID, &item.Name, &item.Amount, &item.Price) //как fmt.Scan()
	if err != nil {
		log.Fatal("Ошибка конвертации", err)
	}
	return item
}


----

carRepository.go

package repositories

import (
	"database/sql"
	"februaryMVCProject/entities"
	"log"
)

type CarRepositories interface {
	AddItem(car entities.Car)
}

type carRepositories struct {
	db *sql.DB
}

func NewCarRepositories(db *sql.DB) CarRepositories {
	return &carRepositories{db: db}
}

func (carRepositories *carRepositories) AddItem(car entities.Car) {
	_, err := carRepositories.db.Exec("insert into cars(name, model, year) values($1, $2, $3)", car.Name, car.Model, car.Year)
	if err != nil {
		log.Fatal("Добавление безуспешно", err)
	}
}



----

entities.go

package entities

type Item struct {
	ID     int64
	Name   string
	Amount int
	Price  int
}

type Car struct {
	ID    int64
	Name  string
	Model string
	Year  int
}


*/

/*
Lesson 5 Rest API - Intro


package main //main.go

import (
	"encoding/json"
	"februaryMVCProject/db"
	"februaryMVCProject/entities"
	"log"
	"net/http"
	"strconv"
)
import _ "github.com/lib/pq"

func main() {
	db.InitDB()
	defer db.CloseDB() //samym poslednim vypolnitsya

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			var idStr string = r.URL.Query().Get("id")
			if idStr == "" {
				err := json.NewEncoder(w).Encode(db.GetAllItems())
				if err != nil {
					log.Fatal("Не удалось получить json всех товаров")
				}
			} else if idStr != "" {
				id, errTwo := strconv.Atoi(idStr)
				if errTwo != nil {
					log.Fatal("Ошибка конвертации")
				}
				errThree := json.NewEncoder(w).Encode(db.GetItem(id))
				if errThree != nil {
					log.Fatal("Ошибка получения json")
				}
			}
		} else if r.Method == http.MethodPost {
			var item entities.Item
			json.NewDecoder(r.Body).Decode(&item)
			db.AddItem(item)
		} else if r.Method == http.MethodPut {
			var updateItem entities.Item
			json.NewDecoder(r.Body).Decode(&updateItem)
			db.UpdateItem(updateItem)
		} else if r.Method == http.MethodDelete {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				log.Fatal("Ошибка парсинга")
			}
			db.DeleteItem(id)
		}
	})

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}

----


package db //dbManager.go

import (
	"database/sql"
	"februaryMVCProject/entities"
	"fmt"
	"log"
)

var dbase *sql.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	dbase, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	err = dbase.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to database")
}
func CloseDB() {
	err := dbase.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func AddItem(item entities.Item) {
	_, err := dbase.Exec("insert into items(name, amount, price) values($1,$2,$3)", item.Name, item.Amount, item.Price)
	if err != nil {

		log.Fatal("Добавление безуспешно,", err)
	}
}

func DeleteItem(id int) {
	_, err := dbase.Exec("delete from items where id = $1", id)
	if err != nil {

		log.Fatal("Удаление безуспешно,", err)
	}
}

func UpdateItem(updateItem entities.Item) {
	_, err := dbase.Exec("update items set name = $1, amount = $2, price = $3 where id = $4", updateItem.Name, updateItem.Amount, updateItem.Price, updateItem.ID)
	if err != nil {

		log.Fatal("Обновление безуспешно,", err)
	}
}

func GetAllItems() []entities.Item {
	var items []entities.Item

	rows, err := dbase.Query("select * from items")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Price) //как fmt.Scan()
		if err != nil {
			log.Fatal("Ошибка конвертации", err)
		}
		items = append(items, item)
	}
	return items
}

func GetItem(id int) entities.Item {
	var item entities.Item

	row := dbase.QueryRow("select * from items where id = $1", id)
	if row == nil {
		log.Fatal("Row is empty")
	}
	err := row.Scan(&item.ID, &item.Name, &item.Amount, &item.Price) //как fmt.Scan()
	if err != nil {
		log.Fatal("Ошибка конвертации", err)
	}
	return item
}

*/

/*
Lesson 4 Rabota Go s bazoi dannyh

func main() { //main
	db.InitDB()
	defer db.CloseDB() //samym poslednim vypolnitsya
	files, err := template.ParseFiles("front/home.html", "front/index.html")
	if err != nil {
		log.Fatal("Файл не найден", err)
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "home.html", db.GetAllItems())
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "index.html", "")
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	http.HandleFunc("/add-item", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodPost {
			name := r.FormValue("item-name")
			amount, err := strconv.Atoi(r.FormValue("item-amount"))
			if err != nil {
				log.Fatal("Введен неправильный тип данных", err)
			}
			price, errTwo := strconv.Atoi(r.FormValue("item-price"))
			if errTwo != nil {
				log.Fatal("Введен неправильный тип данных", errTwo)
			}
			var item entities.Item
			item.Name = name
			item.Amount = amount
			item.Price = price

			db.AddItem(item)

		}
	})

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}

--

package db //dbManager.go

import (
	"database/sql"
	"februaryMVCProject/entities"
	"fmt"
	"log"
)

var dbase *sql.DB

func InitDB() {
	connection := "user=postgres password=admin dbname=FirstDB host=localhost port=5432 sslmode=disable"
	var err error
	dbase, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	err = dbase.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to database")
}
func CloseDB() {
	err := dbase.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func AddItem(item entities.Item) {
	_, err := dbase.Exec("insert into items(name, amount, price) values($1,$2,$3)", item.Name, item.Amount, item.Price)
	if err != nil {
		log.Fatal("Добавление безуспешно,", err)
	}
}

func GetAllItems() []entities.Item {
	var items []entities.Item

	rows, err := dbase.Query("select * from items")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Price) //как fmt.Scan()
		if err != nil {
			log.Fatal("Ошибка конвертации", err)
		}
		items = append(items, item)
	}
	return items
}

----

package entities //entities

type Item struct {
	ID     int64
	Name   string
	Amount int
	Price  int
}

--
index.html, home.html ozgermedi
*/

/*
Lesson 3 Atributy i SQL


func main() {
	var db entities.DB

	files, err := template.ParseFiles("front/home.html", "front/index.html")
	if err != nil {
		log.Fatal("Файл не найден", err)
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "home.html", db.Items)
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "index.html", "")
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	http.HandleFunc("/add-item", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodPost {
			name := r.FormValue("item-name")
			amount, err := strconv.Atoi(r.FormValue("item-amount"))
			if err != nil {
				log.Fatal("Введен неправильный тип данных", err)
			}
			price, errTwo := strconv.Atoi(r.FormValue("item-price"))
			if errTwo != nil {
				log.Fatal("Введен неправильный тип данных", errTwo)
			}
			var item entities.Item
			item.Name = name
			item.Amount = amount
			item.Price = price

			db.Items = append(db.Items, item)
			fmt.Println(db.Items)
		}
	})

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}


<!DOCTYPE html> //index.html
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
</head>
<body>
 <div class="container">
     <form action="/add-item" method="post">
         <div class="row">
             <label class="form-label">NAME:</label>
             <input type="text" class="form-control" name="item-name">
         </div>
         <div class="row">
             <label class="form-label">AMOUNT:</label>
             <input type="number" class="form-control" name="item-amount">
         </div>
         <div class="row">
             <label class="form-label">PRICE:</label>
             <input type="number" class="form-control" name="item-price">
         </div>
         <div class="row">
             <button class="btn btn-info">ADD ITEM</button>
         </div>
     </form>
 </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>
</html>



<!DOCTYPE html> //home.html
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
</head>
<body>
<div class="container">
    <div class="row">
        {{range .}}
        <div class="col-md-3">
            <div class="card">
                <h3 class="card-title">{{.Name}}</h3>
                <h4 class="card-text">{{.Amount}}</h4>
                <h4 class="card-text">{{.Price}}</h4>
                <button class="btn btn-info">ABOUT</button>
            </div>
        </div>
        {{end}}
    </div>
</div>
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>
</html>


create table items( //SQL commands
    id serial primary key,
    name varchar(100),
    amount int,
    price int
);

insert into items(name, amount, price)
values('Iphone17Air', 100, 900000);

insert into items(name, amount, price)
values('Iphone16Pro', 100, 600000);

update items
set name = 'Iphone17Pro', amount = 150, price = 950000
where id = 1;

delete from items
where id = 2;

select * from items

*/

/*
Lesson 2 Post zapros

func main() {
	var db entities.DB

	files, err := template.ParseFiles("front/home.html", "front/index.html")
	if err != nil {
		log.Fatal("Файл не найден", err)
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "home.html", "")
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "index.html", "")
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	http.HandleFunc("/add-item", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodPost {
			name := r.FormValue("item-name")
			amount, err := strconv.Atoi(r.FormValue("item-amount"))
			if err != nil {
				log.Fatal("Введен неправильный тип данных", err)
			}
			price, errTwo := strconv.Atoi(r.FormValue("item-price"))
			if errTwo != nil {
				log.Fatal("Введен неправильный тип данных", errTwo)
			}
			var item entities.Item
			item.Name = name
			item.Amount = amount
			item.Price = price

			db.Items = append(db.Items, item)
			fmt.Println(db.Items)
		}
	})

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}

<!DOCTYPE html> //index.html
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
</head>
<body>
 <div class="container">
     <form action="/add-item" method="post">
         <div class="row">
             <label class="form-label">NAME:</label>
             <input type="text" class="form-control" name="item-name">
         </div>
         <div class="row">
             <label class="form-label">AMOUNT:</label>
             <input type="number" class="form-control" name="item-amount">
         </div>
         <div class="row">
             <label class="form-label">PRICE:</label>
             <input type="number" class="form-control" name="item-price">
         </div>
         <div class="row">
             <button class="btn btn-info">ADD ITEM</button>
         </div>
     </form>
 </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>
</html>


<!DOCTYPE html> //home.html
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<div class="container">
    <div class="row">
        <div class="col-md-3">
            <div class="card">

            </div>
        </div>
        <div class="col-md-3">
            <div class="card">

            </div>
        </div>
        <div class="col-md-3">
            <div class="card">

            </div>
        </div>
        <div class="col-md-3">
            <div class="card">

            </div>
        </div>
    </div>
</div>
</body>
</html>




*/

/*
Lesson 1 Vvedenie v Klient Servernuyu razrabotku
func main() {
	files, err := template.ParseFiles("front/home.html", "front/index.html")
	if err != nil {
		log.Fatal("Файл не найден", err)
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "home.html", "")
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) { //w-writer, r-request
		if r.Method == http.MethodGet {
			errThree := files.ExecuteTemplate(w, "index.html", "")
			if errThree != nil {
				log.Fatal("Ошибка вовзращения ответа,", errThree)
			}
		}
	})

	server := &http.Server{
		Addr: "localhost:8080",
	}

	errTwo := server.ListenAndServe()
	if errTwo != nil {
		log.Fatal("Ошибка сервера :", errTwo)
	}
}

<!DOCTYPE html> //home.html
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<h1>Bitlab Academy</h1>
<h2>Bitlab Academy</h2>
<h3>Bitlab Academy</h3>
<h4>Bitlab Academy</h4>
<h5>Bitlab Academy</h5>
<h6>Bitlab Academy</h6>
</body>
</html>


<!DOCTYPE html> //index.html
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<h1>Almaty</h1>
<h2>Almaty</h2>
<h3>Almaty</h3>
<h4>Almaty</h4>
<h5>Almaty</h5>
<h6>Almaty</h6>
</body>
</html>
*/
