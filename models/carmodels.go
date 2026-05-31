package models

import "time"

type CarParameters struct {
	CarID             int       `json:"id"`       // mašīnas id
	TeamName          string    `json:"username"` // komandas nosaukums tho
	SetVoltage        float32   `json:"U"`
	CalculatedCurrent float32   `json:"I"`
	Mass              float32   `json:"m"`
	AgeGroup          string    `json:"ageGroup"`
	Avatar            []byte    `json:"avatar"` // tehniski base64 enkodēts binārs fails - bilde
	FinishedAt        time.Time `json:"finishedAt"`
}
