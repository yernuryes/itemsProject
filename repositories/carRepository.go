package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type CarRepository interface {
	AddCar(car entities.Car)
	DeleteCar(id int)
	UpdateCar(updateCar entities.Car)
	GetAllCars() []entities.Car
	GetCar(id int) entities.Car
}

type carRepository struct {
	gormDB *gorm.DB
}

func NewCarRepository(gormDB *gorm.DB) CarRepository {
	return &carRepository{gormDB: gormDB}
}

func (carRepository *carRepository) AddCar(car entities.Car) {
	carRepository.gormDB.Create(&car)
}

func (carRepository *carRepository) DeleteCar(id int) {
	carRepository.gormDB.Delete(entities.Car{}, id)
}

func (carRepository *carRepository) UpdateCar(updateCar entities.Car) {
	carRepository.gormDB.Save(&updateCar)
}

func (carRepository *carRepository) GetAllCars() []entities.Car {
	var cars []entities.Car

	carRepository.gormDB.Find(&cars)

	return cars
}

func (carRepository *carRepository) GetCar(id int) entities.Car {
	var car entities.Car
	carRepository.gormDB.First(&car, id)
	return car
}
