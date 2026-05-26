package mappers

import (
	"februaryMVCProject/dtos"
	"februaryMVCProject/entities"
)

func MapToDTOCar(car entities.Car) dtos.CarDTO { // iz itema hochu sdelat' itemDTO
	var carDTO dtos.CarDTO
	carDTO.ID = car.ID
	carDTO.Name = car.Name
	carDTO.Model = car.Model
	carDTO.Year = car.Year
	return carDTO
}

func MapToDTOListCar(cars []entities.Car) []dtos.CarDTO { // iz itema hochu sdelat' itemDTO
	var carDTOs []dtos.CarDTO
	for i := 0; i < len(cars); i++ {
		var carDTO dtos.CarDTO
		carDTO.ID = cars[i].ID
		carDTO.Name = cars[i].Name
		carDTO.Model = cars[i].Model
		carDTO.Year = cars[i].Year
		carDTOs = append(carDTOs, carDTO)
	}
	return carDTOs
}
