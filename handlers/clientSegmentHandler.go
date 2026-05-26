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
