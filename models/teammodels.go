package models

type TeamData struct {
	Key               string            `json:"key"`
	CarID             string            `json:"carId"`
	TeamName          string            `json:"teamName"`
	Accounts          []Account         `json:"members"`
	AgeGroup          string            `json:"ageGroup"`
	Institution       string            `json:"institution,omitempty"`
	CityOrRegion      string            `json:"cityOrRegion"`
	ResponsiblePerson ResponsiblePerson `json:"responsiblePerson"`
	Avatar            string            `json:"avatar"`
	CarData           Car               `json:"carData,omitempty"`
}
