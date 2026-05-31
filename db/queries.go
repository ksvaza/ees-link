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
			attempt, raceName
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

	// Result Entries
	/*
	   type ResultsEntry struct {
	   	// filtri
	   	RaceName  string `json:"race_name"`
	   	AttemptNr int    `json:"attempt_nr"` // Mēģinājuma Nr. 				// EkoRace, DragRace, ShuttleRun, SteeringTest, BrakingTest

	   	// mašīnas dati
	   	CarID    int    `json:"car_id"`    // ID
	   	TeamName string `json:"team_name"` // Komanda

	   	// laiki
	   	StartTime  time.Time `json:"start_time"`
	   	FinishTime time.Time `json:"finish_time"`

	   	// Manual entry
	   	Penalties       time.Duration `json:"penalties"`        // Sodi (s) 					// EkoRace, MainRace, ShuttleRun, SteeringTest
	   	UsedTime        time.Duration `json:"used_time"`        // Patērētais laiks 			// EkoRace, MainRace, ShuttleRun, SteeringTest
	   	MaxSpeed        float64       `json:"max_speed"`        // Maks. ātrums (km/h) 			// DragRace
	   	Valid           bool          `json:"valid"`            // Derīgs 						// EkoRace, DragRace, ShuttleRun, SteeringTest, BrakingTest
	   	BrakingDistance float64       `json:"braking_distance"` // Bremzēšanas distance (m) 	// BrakingTest
	   	DriveinSpeed    float64       `json:"drivein_speed"`    // Iebraukšanas ātrums (km/h)	// BrakingTest

	   	// Calculated fields
	   	AverageSpeed    float64       `json:"average_speed"`    // Vid. ātrums = distance / TotalTime 					// EkoRace, MainRace
	   	TotalTime       time.Duration `json:"total_time"`       // Kopējais laiks = patērētais laiks + sodi 			// EkoRace, MainRace, DragRace, ShuttleRun, SteeringTest
	   	Placement       int           `json:"placement"`        // Vieta = mainrace:totaltime, ekorace: skatās secību 	// EkoRace, MainRace
	   	UsedEnergy      float64       `json:"used_energy"`      // Patērētā enerģija (Wh) 								// EkoRace
	   	Efficiency      float64       `json:"efficiency"`       // Efektivitāte (Wh/kg) 								// EkoRace
	   	ShellEfficiency float64       `json:"shell_efficiency"` // "Shell" efektivitāte (km/kWh) 						// EkoRace
	   }
	*/
	ResultsEntryWriteRequest = `
		INSERT INTO results_entries (
			race_name, attempt_nr, car_id, team_name, start_time, finish_time, penalties, used_time, max_speed, valid, braking_distance, drivein_speed, average_speed, total_time, placement, used_energy, efficiency, shell_efficiency
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`

	ResultsEntryUpdateRequest = `
		UPDATE results_entries SET
			team_name = $4,
			start_time = $5,
			finish_time = $6,
			penalties = $7,
			used_time = $8,
			max_speed = $9,
			valid = $10,
			braking_distance = $11,
			drivein_speed = $12,
			average_speed = $13,
			total_time = $14,
			placement = $15,
			used_energy = $16,
			efficiency = $17,
			shell_efficiency = $18
		WHERE race_name = $1 AND attempt_nr = $2 AND car_id = $3`

	ResultsEntryUpsertRequest = `
		INSERT INTO results_entries (
			race_name, attempt_nr, car_id, team_name, start_time, finish_time, penalties, used_time, max_speed, valid, braking_distance, drivein_speed, average_speed, total_time, placement, used_energy, efficiency, shell_efficiency
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) ON CONFLICT (race_name, attempt_nr, car_id) DO UPDATE SET
			team_name = EXCLUDED.team_name,
			start_time = EXCLUDED.start_time,
			finish_time = EXCLUDED.finish_time,
			penalties = EXCLUDED.penalties,
			used_time = EXCLUDED.used_time,
			max_speed = EXCLUDED.max_speed,
			valid = EXCLUDED.valid,
			braking_distance = EXCLUDED.braking_distance,
			drivein_speed = EXCLUDED.drivein_speed,
			average_speed = EXCLUDED.average_speed,
			total_time = EXCLUDED.total_time,
			placement = EXCLUDED.placement,
			used_energy = EXCLUDED.used_energy,
			efficiency = EXCLUDED.efficiency,
			shell_efficiency = EXCLUDED.shell_efficiency`

	ResultsEntriesReadByRaceNameRequest = `
		SELECT race_name, attempt_nr, car_id, team_name, start_time, finish_time, penalties, used_time, max_speed, valid, braking_distance, drivein_speed, average_speed, total_time, placement, used_energy, efficiency, shell_efficiency
		FROM results_entries WHERE race_name = $1`

	ResultsEntryReadByCarID = `
		SELECT race_name, attempt_nr, car_id, team_name, start_time, finish_time, penalties, used_time, max_speed, valid, braking_distance, drivein_speed, average_speed, total_time, placement, used_energy, efficiency, shell_efficiency
		FROM results_entries WHERE car_id = $1`
)
