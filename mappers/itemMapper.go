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
