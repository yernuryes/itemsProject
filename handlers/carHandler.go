package handlers

import (
	"encoding/json"
	"februaryMVCProject/entities"
	"februaryMVCProject/services"
	"log"
	"net/http"
	"strconv"
)

type CarHandler interface {
	HandleCars(w http.ResponseWriter, r *http.Request)
}

type carHandler struct {
	carService services.CarService
}

func NewCarHandler(carService services.CarService) CarHandler {
	return &carHandler{carService: carService}
}

func (carHandler *carHandler) HandleCars(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var idStr string = r.URL.Query().Get("id")
		if idStr == "" {
			err := json.NewEncoder(w).Encode(carHandler.carService.GetAllCars())
			if err != nil {
				log.Fatal("Не удалось получить json всех товаров")
			}
		} else if idStr != "" {
			id, errTwo := strconv.Atoi(idStr)
			if errTwo != nil {
				log.Fatal("Ошибка конвертации")
			}
			errThree := json.NewEncoder(w).Encode(carHandler.carService.GetCar(id))
			if errThree != nil {
				log.Fatal("Ошибка получения json")
			}
		}
	} else if r.Method == http.MethodPost {
		var car entities.Car
		json.NewDecoder(r.Body).Decode(&car)
		carHandler.carService.AddCar(car)
	} else if r.Method == http.MethodPut {
		var updateCar entities.Car
		json.NewDecoder(r.Body).Decode(&updateCar)
		carHandler.carService.UpdateCar(updateCar)
	} else if r.Method == http.MethodDelete {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal("Ошибка парсинга")
		}
		carHandler.carService.DeleteCar(id)
	}
}
