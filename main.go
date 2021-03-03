package main

import (
	"log"
	"os"

	"github.com/SubrotoRoy/pre-parking/datastore"
	"github.com/SubrotoRoy/pre-parking/handler"
	"github.com/SubrotoRoy/pre-parking/kafkaservice"
	"github.com/SubrotoRoy/pre-parking/model"
	"github.com/labstack/echo"
)

func main() {
	err := initializeGormClient()

	service := kafkaservice.NewKafkaService()
	api := handler.NewParkHandler(service, datastore.Repo)
	if err != nil {
		log.Fatalf("Couldn't connect to database. ERROR: %s", err.Error())
	}

	e := echo.New()
	g := e.Group("/api/v1")

	g.POST("/car/park", api.ParkCar)
	g.GET("/car/:carNumber/unpark", api.UnParkCar)
	e.Logger.Fatal(e.Start(":8091"))
}

//initializeGormClient initializes the DBRepo
func initializeGormClient() error {
	log.Println("Connecting to db")
	config := model.DbConfig{
		DbUser:     os.Getenv("DBUSER"),
		DbPassword: os.Getenv("DBPASSWORD"),
		DbName:     os.Getenv("DBNAME"),
		Port:       os.Getenv("PORT"),
		Host:       os.Getenv("Host"),
	}
	datastore.Repo = &datastore.DBRepo{}
	return datastore.Repo.DBConnect(config)
}
