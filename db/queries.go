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
)
