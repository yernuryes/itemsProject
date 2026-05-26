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

