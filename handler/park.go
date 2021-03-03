package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/labstack/echo"

	"github.com/SubrotoRoy/pre-parking/datastore"
	"github.com/SubrotoRoy/pre-parking/kafkaservice"
	"github.com/SubrotoRoy/pre-parking/model"
)

//ParkHandler will be used as receiver to associate handler functions to it
type ParkHandler struct {
	Kafka    kafkaservice.Services
	DataRepo datastore.Repository
}

//NewParkHandler return an instance of ParkHandler
func NewParkHandler(kafka kafkaservice.Services, db datastore.Repository) *ParkHandler {
	return &ParkHandler{Kafka: kafka, DataRepo: db}
}

//ParkCar handles the requests to /car/park POST call
func (p *ParkHandler) ParkCar(ctx echo.Context) error {

	decodedRequest := model.Parking{}

	defer ctx.Request().Body.Close()

	err := json.NewDecoder(ctx.Request().Body).Decode(&decodedRequest)

	if err != nil {
		log.Println("Error encountered while decoding JSON request body, ERROR:", err)
		return ctx.JSON(400, ResponseManager(nil, errors.New("Error encountered while decoding JSON request body")))
	}
	if decodedRequest.Car.CarNumber == "" {
		log.Println("Invalid Car number provided")
		return ctx.JSON(400, ResponseManager(nil, errors.New("Invalid Car number provided")))
	}
	carParked := p.DataRepo.IsCarParked(decodedRequest.Car.CarNumber)
	if carParked {
		log.Println("Car already parked")
		return ctx.JSON(400, ResponseManager(nil, errors.New("Car already parked")))
	}

	slotAvailable := p.DataRepo.IsSlotAvailable()
	if !slotAvailable {
		return ctx.JSON(404, ResponseManager(nil, errors.New("No empty parking slot available")))
	}
	c := context.Background()
	err = p.Kafka.WriteToKafka(c, "park", decodedRequest)
	if err != nil {
		log.Println("Unable to post to kafka, ERROR:", err)
		return ctx.JSON(500, ResponseManager(nil, errors.New("Unable to post to kafka")))
	}
	return ctx.JSON(201, ResponseManager("Car Will be parked", nil))
}

//UnParkCar handles the requests to /car/{id}/unpark GET call
func (p *ParkHandler) UnParkCar(ctx echo.Context) error {
	defer ctx.Request().Body.Close()

	carNumber := ctx.Param("carNumber")

	if carNumber == "" {
		log.Println("Invalid id provided:", carNumber)
		return ctx.JSON(400, ResponseManager(nil, errors.New("Invalid car number provided")))
	}

	carParked := p.DataRepo.IsCarParked(carNumber)
	if !carParked {
		return ctx.JSON(404, ResponseManager(nil, errors.New("No car parked with given number")))
	}

	c := context.Background()

	err := p.Kafka.WriteToKafka(c, "unpark", model.Parking{Car: model.Car{CarNumber: carNumber}})

	if err != nil {
		log.Println("Unable to post to kafka, ERROR:", err)
		return ctx.JSON(500, ResponseManager(nil, errors.New("Unable to post to kafka")))
	}
	return ctx.JSON(200, ResponseManager("Car will be fetched from parking", nil))
}

//ResponseManager manages stuffs
func ResponseManager(response interface{}, err error) model.APIResponse {

	apiResponse := model.APIResponse{}

	apiResponse.Response = response
	if err != nil {
		apiResponse.Error = err.Error()
	}
	return apiResponse
}
