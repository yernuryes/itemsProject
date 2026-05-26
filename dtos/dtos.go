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
