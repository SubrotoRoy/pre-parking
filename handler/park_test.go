package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SubrotoRoy/pre-parking/datastore"
	"github.com/SubrotoRoy/pre-parking/handler"
	"github.com/SubrotoRoy/pre-parking/kafkaservice"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var mockDBRepo *datastore.MockDBRepo
var mockKafkaSvc *kafkaservice.MockKafkaSvc

var api *handler.ParkHandler

func init() {
	mockDBRepo = &datastore.MockDBRepo{}
	//datastore.Repo = mockDBRepo
	mockKafkaSvc = &kafkaservice.MockKafkaSvc{}
	api = handler.NewParkHandler(mockKafkaSvc, mockDBRepo)
}

func TestParkCarInvalidJSON(t *testing.T) {
	var inputJSON = ``

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/car/park", strings.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	if assert.NoError(t, api.ParkCar(context)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestParkCarNoCarNumber(t *testing.T) {
	var inputJSON = `{
		"car": {
			"carModel": "Nissan5"
		},
		"person": {
			"firstName": "Roy5",
			"lastName": "Subroto5",
			"phoneNumber":79895
		}
	}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/car/park", strings.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	if assert.NoError(t, api.ParkCar(context)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestParkCarAlreadyParked(t *testing.T) {
	var inputJSON = `{
		"car": {
			"carNumber": "KA335",
			"carModel": "Nissan5"
		},
		"person": {
			"firstName": "Roy5",
			"lastName": "Subroto5",
			"phoneNumber":79895
		}
	}`

	mockDBRepo.On("IsCarParked", "KA335").Return(true)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/car/park", strings.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	if assert.NoError(t, api.ParkCar(context)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestParkCarNoSlot(t *testing.T) {
	var inputJSON = `{
		"car": {
			"carNumber": "KA336",
			"carModel": "Nissan5"
		},
		"person": {
			"firstName": "Roy5",
			"lastName": "Subroto5",
			"phoneNumber":79895
		}
	}`

	mockDBRepo.On("IsCarParked", "KA336").Return(false)
	mockDBRepo.On("IsSlotAvailable").Return(false)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/car/park", strings.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	if assert.NoError(t, api.ParkCar(context)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
