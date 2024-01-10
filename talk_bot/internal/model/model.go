package model

type Device struct {
	Company  string  `json:"company"`
	Model    string  `json:"model"`
	Country  string  `json:"country"`
	Price    float64 `json:"price"`
	Customer string
}
type Devices struct {
	Devices []Device `json:"devices"`
}
