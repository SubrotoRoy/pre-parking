package model

//APIResponse struct is to send out reponse
type APIResponse struct {
	Response interface{} `json:"response"`
	Error    string      `json:"error"`
}

//DbConfig is used for connection to database
type DbConfig struct {
	DbUser     string
	DbPassword string
	DbName     string
	Port       string
	Host       string
}

//Person struct will be used for mapping to database table person
type Person struct {
	ID          int
	FirstName   string
	LastName    string
	PhoneNumber int
}

//Car struct will be used for mapping to database table car
type Car struct {
	ID        int
	CarModel  string
	CarNumber string
}

//Slot struct will be used for mapping to databaase table slot
type Slot struct {
	ID         int
	IsOccupied bool
}

//Parking struct will be used for mapping to databaase table parking
type Parking struct {
	ID     int
	Person Person
	Car    Car
}
