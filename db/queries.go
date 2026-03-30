package db

const (
	ApplicantWriteRequest = `INSERT INTO applicants (ID, team_name, school, members, supervisor, applied_at, status)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, team_name, school, members, supervisor, applied_at, status`
)
