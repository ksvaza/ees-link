package db

const (
	// -------------------------------------------------------------------------
	// Applicants (applicants table)
	// -------------------------------------------------------------------------

	TeamApplicationWriteRequest = `
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

	TeamApplicationReadRequest = `
		SELECT id, team_name, age_group, institution, city_or_region, members,
			responsible_person, how_heard_about, comments,
			confirm_truthful, confirm_rules, confirm_media, applied_at, status
		FROM applicants`

	TeamApplicationReadRequestByID = `
		SELECT id, team_name, age_group, institution, city_or_region, members,
			responsible_person, how_heard_about, comments,
			confirm_truthful, confirm_rules, confirm_media, applied_at, status
		FROM applicants WHERE id = $1`

	TeamApplicationUpdateRequest = `
		UPDATE applicants SET
			team_name = $1, age_group = $2, institution = $3, city_or_region = $4,
			members = $5, responsible_person = $6,
			how_heard_about = $7, comments = $8,
			confirm_truthful = $9, confirm_rules = $10, confirm_media = $11, applied_at = $12,
			status = $13
		WHERE id = $14`

	TeamApplicationReadRequestByTeamName = `
		SELECT id, team_name, age_group, institution, city_or_region, members,
			responsible_person, how_heard_about, comments,
			confirm_truthful, confirm_rules, confirm_media, applied_at, status
		FROM applicants WHERE team_name = $1`

	// -------------------------------------------------------------------------
	// Accounts (konti table)
	// -------------------------------------------------------------------------

	AccountWriteRequest = `
		INSERT INTO konti (
			key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt, educational_institution, class_or_year, pending_team_id, registered, avatar
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	AccountReadRequest = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt, educational_institution, class_or_year, pending_team_id, registered, avatar
		FROM konti`

	AccountReadRequestByUsername = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt, educational_institution, class_or_year, pending_team_id, registered, avatar
		FROM konti WHERE username = $1`

	AccountReadRequestByFullname = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt, educational_institution, class_or_year, pending_team_id, registered, avatar
		FROM konti WHERE fullname = $1`

	AccountReadRequestByDateOfBirth = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt, educational_institution, class_or_year, pending_team_id, registered, avatar
		FROM konti WHERE date_of_birth = $1`

	AccountReadRequestByKey = `
		SELECT key, fullname, date_of_birth, role, id, password, username, email, phone_number, salt, educational_institution, class_or_year, pending_team_id, registered, avatar
		FROM konti WHERE key = $1`

	AccountUpdateRequest = `
		UPDATE konti SET
			fullname = $1, date_of_birth = $2, role = $3, id = $4,
			password = $5, username = $6, email = $7, phone_number = $8, salt = $9,
			educational_institution = $10, class_or_year = $11, pending_team_id = $12, registered = $13, avatar = $14
		WHERE key = $15`
	// -------------------------------------------------------------------------
	// Account applications (kontu_pieteikumi table)
	// -------------------------------------------------------------------------

	AccountApplicationWriteRequest = `
		INSERT INTO kontu_pieteikumi (
			fullname, date_of_birth, role, password, username, email, phone_number, team_name, key
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	AccountApplicationReadRequest = `
		SELECT fullname, date_of_birth, role, password, username, email, phone_number, team_name, key
		FROM kontu_pieteikumi`

	AccountApplicationReadRequestByKey = `
		SELECT fullname, date_of_birth, role, password, username, email, phone_number, team_name, key
		FROM kontu_pieteikumi WHERE key = $1`

	AccountApplicationDeleteRequestByKey = `
		DELETE FROM kontu_pieteikumi WHERE key = $1`

	// -------------------------------------------------------------------------
	// Admins (admini table)
	// -------------------------------------------------------------------------

	AdminAccountWriteRequest = `
		INSERT INTO admini (
			username, password, salt, superadmin, key
		) VALUES (
			$1, $2, $3, $4, $5
		)`

	AdminAccountReadRequest = `
		SELECT username, password, salt, superadmin, key
		FROM admini`

	AdminAccountReadRequestByUsername = `
		SELECT username, password, salt, superadmin, key
		FROM admini WHERE username = $1`

	AdminAccountReadRequestByKey = `
		SELECT username, password, salt, superadmin, key
		FROM admini WHERE key = $1`

	AdminAccountUpdateRequest = `
		UPDATE admini SET
			username = $1, password = $2, salt = $3, superadmin = $4
		WHERE key = $5`

	// Team data, cars

	TeamDataReadRequestByKey = `SELECT
    key, car_id, team_name, team_members, age_group, institution,
    city_or_region, responsible_person, car_data, avatar
FROM komandas
WHERE key = $1;`

	TeamDataUpdateRequest = `UPDATE komandas SET 
    car_id = $1, team_name = $2, team_members = $3, age_group = $4, institution = $5,
    city_or_region = $6, responsible_person = $7, car_data = $8, avatar = $9
WHERE key = $10`

	TeamMembersWriteRequest = `UPDATE komandas SET team_members = $1::text[] WHERE key = $2`

	TeamDataWriteRequest = `INSERT INTO komandas (
	key, car_id, team_name, team_members, age_group, institution,
	city_or_region, responsible_person, car_data, avatar
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	// MQTT things

	MQTTLogSaveRequest = `
		INSERT INTO mqtt_logs (
			topic, payload, received_at
		) VALUES (
			$1, $2, $3
		)`

	MQTTLogsReadRequest = `
		SELECT topic, payload, received_at
		FROM mqtt_logs
		LIMIT $1`

	CarTelemetrySaveRequest = `
		INSERT INTO car_telemetry (
			id, rssi, 
			timestamp, 
			accel_x, accel_y, accel_z, 
			gps_lat, gps_lon, gps_spd, gps_sat, 
			psu_uop, psu_iop, psu_pop, psu_uip, psu_wh, 
			sys_vbat, sys_dc, sys_err, 
			meginajums
		)`

	// Car parameters
	GetAllCarParametersRequest = `
		SELECT car_id, team_name, set_voltage, calculated_current, mass, age_group, avatar, finished_at
		FROM car_parameters`
	GetCarParametersByAgeGroupRequest = `
		SELECT car_id, team_name, set_voltage, calculated_current, mass, age_group, avatar, finished_at
		FROM car_parameters
		WHERE age_group = $1`
	DeleteAllCarParametersRequest = `
		DELETE FROM car_parameters`
	InsertCarParametersRequest = `
		INSERT INTO car_parameters (
			car_id, team_name, set_voltage, calculated_current, mass, age_group, avatar, finished_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)`

	// Points
	PointsSaveRequest = `
		INSERT INTO points (
			race_name, points_data
		) VALUES (
			$1, $2
		)`
	PointsReadAllRequest = `
		SELECT race_name, points_data
		FROM points`
	PointsReadByRaceNameRequest = `
		SELECT race_name, points_data
		FROM points
		WHERE race_name = $1`
	PointsUpdateByRaceNameRequest = `
		UPDATE points SET
			points_data = $1
		WHERE race_name = $2`
)
