package db

const (
	ApplicantWriteRequest = `
		INSERT INTO applicants (
			id, team_name, age_group, institution, city_or_region, members,
			responsible_person, how_heard_about, comments,
			confirm_truthful, confirm_rules, confirm_media, applied_at, status
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12,
			$13, $14
		)`

	ApplicantReadRequest = `
		SELECT id, team_name, age_group, institution, city_or_region, members,
			responsible_person, how_heard_about, comments,
			confirm_truthful, confirm_rules, confirm_media, applied_at, status
		FROM applicants`

	ApplicantReadRequestByID = `
		SELECT id, team_name, age_group, institution, city_or_region, members,
			responsible_person, how_heard_about, comments,
			confirm_truthful, confirm_rules, confirm_media, applied_at, status
		FROM applicants WHERE id = $1`

	ApplicantUpdateRequest = `
		UPDATE applicants SET
			team_name = $1, age_group = $2, institution = $3, city_or_region = $4,
			members = $5, responsible_person = $6,
			how_heard_about = $7, comments = $8,
			confirm_truthful = $9, confirm_rules = $10, confirm_media = $11, applied_at = $12,
			status = $13
		WHERE id = $14`

	AccountWriteRequest = `
		INSERT INTO konti (
			id = $1, full_name = $2, date_of_birth = $3, password = $4, username = $5, email = $6, phone_number = $7
		)`

	AccountReadRequestByFullname = `
		SELECT id, full_name, date_of_birth, educational_institution, role, class_or_year, password, username, email, phone_number
		FROM konti WHERE full_name = $1`

	AccountReadRequestByDateOfBirth = `
		SELECT id, full_name, date_of_birth, educational_institution, role, class_or_year, password, username, email, phone_number
		FROM konti WHERE date_of_birth = $1`

	AccountApplicationWriteRequest = `
		INSERT INTO kontu_pieteikumi (
			full_name = $1, date_of_birth = $2, password = $3, username = $4, email = $5, phone_number = $6, team_name = $7 )`

	AccountReadRequestByUsername = `
		SELECT id, full_name, date_of_birth, educational_institution, role, class_or_year, password, username, email, phone_number
		FROM konti WHERE username = $1`

		AccountApplicationReadRequest = `
		SELECT id, full_name, date_of_birth, password, username, email, phone_number, team_name
		FROM kontu_pieteikumi`
)
