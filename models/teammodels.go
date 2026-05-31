package models

type TeamData struct {
	Key               string            `json:"key"`
	CarID             int               `json:"carId"`
	TeamName          string            `json:"teamName"`
	Accounts          []Account         `json:"members"`
	AgeGroup          string            `json:"ageGroup"`
	Institution       string            `json:"institution,omitempty"`
	CityOrRegion      string            `json:"cityOrRegion"`
	ResponsiblePerson ResponsiblePerson `json:"responsiblePerson"`
	Avatar            []byte            `json:"avatar"` // tehniski base64 enkodēts binārs fails - bilde
	CarData           Car               `json:"carData,omitempty"`
}
