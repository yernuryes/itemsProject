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
