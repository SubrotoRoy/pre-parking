package datastore

import (
	"fmt"
	"log"

	"github.com/SubrotoRoy/pre-parking/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Repo is a variable of type Repository interface
var Repo Repository

//Repository has all the database interactions
type Repository interface {
	DBConnect(config model.DbConfig) error
	IsSlotAvailable() bool
	IsCarParked(carNumber string) bool
}

//DBRepo satisfies the interface by implementing all the methods
type DBRepo struct {
	GormDB *gorm.DB
}

//DBConnect Method to connect to Db
func (dc *DBRepo) DBConnect(config model.DbConfig) error {
	var err error
	// Format DB configs to required format connect DB
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.DbUser, config.DbPassword, config.DbName, config.Port)

	dc.GormDB, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})

	if err != nil {
		log.Printf("Unable to connect DB %v", err)
		return err
	}
	log.Printf("Postgres started at %s PORT", config.Port)
	return err
}

//IsSlotAvailable checks if there are empty slots
func (dc *DBRepo) IsSlotAvailable() bool {
	availableRows := dc.GormDB.Debug().Where(`is_occupied != ?`, true).Find(&model.Slot{}).RowsAffected

	return availableRows > 0
}

//IsCarParked checks if there are empty slots
func (dc *DBRepo) IsCarParked(carNumber string) bool {
	car := model.Car{}
	carRows := dc.GormDB.Debug().Where(`car_number = ?`, carNumber).Find(&car).RowsAffected
	log.Println("carRows", carRows)
	if carRows == 0 {
		return false
	}
	var parkingIds []int
	parkingRows := dc.GormDB.Debug().Raw(`Select id from parkings where car_id = ? and has_exited is not true`, car.ID).Scan(&parkingIds).RowsAffected
	return parkingRows > 0
}
