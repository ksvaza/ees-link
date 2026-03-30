package models

import "github.com/jackc/pgx/v5/pgtype"

type Applicant struct {
	ID         string      `json:"id"`
	TeamName   string      `json:"teamName"`
	School     string      `json:"school"`
	Members    int         `json:"members"`
	Supervisor string      `json:"supervisor"`
	AplliedAt  pgtype.Date `json:"appliedAt"`
	Status     string      `json:"status"`
}
