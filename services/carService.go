package services

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
	"februaryMVCProject/mappers"
	"februaryMVCProject/repositories"
	"github.com/google/uuid"
	"log"
)

type CarService interface {
	AddCar(car entities.Car)
	DeleteCar(id int)
	UpdateCar(updateItem entities.Car)
	GetAllCars() []dtos.CarDTO
	GetCar(id int) dtos.CarDTO
}

type carService struct {
	carRepository repositories.CarRepository
}

func NewCarService(carRepository repositories.CarRepository) CarService {
	return &carService{carRepository: carRepository}
}

func (carService *carService) AddCar(car entities.Car) {
	if car.Name == "" {
		log.Fatal("Без названии нельзя добавить")
	} else {
		car.Promocode = uuid.New().String()
		carService.carRepository.AddCar(car)
	}
}

func (carService *carService) DeleteCar(id int) {
	carService.carRepository.DeleteCar(id)
}

func (carService *carService) UpdateCar(updateCar entities.Car) {
	var currentCar entities.Car = carService.carRepository.GetCar(int(updateCar.ID))
	updateCar.Promocode = currentCar.Promocode
	carService.carRepository.UpdateCar(updateCar)
}

func (carService *carService) GetAllCars() []dtos.CarDTO {
	return mappers.MapToDTOListCar(carService.carRepository.GetAllCars())
}

func (carService *carService) GetCar(id int) dtos.CarDTO {
	return mappers.MapToDTOCar(carService.carRepository.GetCar(id))
}
