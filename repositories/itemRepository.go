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
