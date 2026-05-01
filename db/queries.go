package db

const (
	// -------------------------------------------------------------------------
	// Applicants (applicants table)
	// -------------------------------------------------------------------------

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

	// -------------------------------------------------------------------------
	// Accounts (konti table)
	// -------------------------------------------------------------------------

	AccountWriteRequest = `
		INSERT INTO konti (
			key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	AccountReadRequest = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt
		FROM konti`

	AccountReadRequestByUsername = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt
		FROM konti WHERE username = $1`

	AccountReadRequestByFullname = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt
		FROM konti WHERE fullname = $1`

	AccountReadRequestByDateOfBirth = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt
		FROM konti WHERE date_of_birth = $1`

	// -------------------------------------------------------------------------
	// Account applications (kontu_pieteikumi table)
	// -------------------------------------------------------------------------

	AccountApplicationWriteRequest = `
		INSERT INTO kontu_pieteikumi (
			fullname, date_of_birth, role, password, username, email, phone_number, team_name
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)`

	AccountApplicationReadRequest = `
		SELECT fullname, date_of_birth, role, password, username, email, phone_number, team_name
		FROM kontu_pieteikumi`

	AccountApplicationReadRequestByID = `
		SELECT fullname, date_of_birth, role, password, username, email, phone_number, team_name
		FROM kontu_pieteikumi WHERE id = $1`

	// -------------------------------------------------------------------------
	// Admins (admini table)
	// -------------------------------------------------------------------------

	AdminAccountWriteRequest = `
		INSERT INTO admini (
			username, password, salt, superadmin
		) VALUES (
			$1, $2, $3, $4
		)`

	AdminAccountReadRequest = `
		SELECT username, password, salt, superadmin
		FROM admini`

	AdminAccountReadRequestByUsername = `
		SELECT username, password, salt, superadmin
		FROM admini WHERE username = $1`
)
