package db

import (
	"context"
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

//go:embed schema.sql
var schema string
var Pool *pgxpool.Pool

func MigrateUp(ctx context.Context, pool *pgxpool.Pool, dbURL string) error {
	var err error
	Pool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to create DB pool")
		return err
	}
	// Example migration: add a new column to an existing table
	fmt.Println("Running migration...")
	_, err = pool.Exec(ctx, schema)
	if err != nil {
		fmt.Println("Migration failed:")
		logrus.WithError(err).Error("Failed to run migration")
		return err
	}
	logrus.Info("Migration completed successfully")
	return nil
}

func WriteApplicantTable(ctx context.Context, pool *pgxpool.Pool, newApplicant models.Applicant) error {
	_, err := pool.Exec(ctx, ApplicantWriteRequest,
		newApplicant.ID,         // $1 - id
		newApplicant.TeamName,   // $2 - team_name
		newApplicant.School,     // $3 - school
		newApplicant.Members,    // $4 - members
		newApplicant.Supervisor, // $5 - supervisor
		newApplicant.AplliedAt,  // $6 - applied_at
		"pending")

	if err != nil {
		logrus.WithError(err).Error("Failed to write to applicants table")
	}

	return err
}

type RandomStruct struct {
	ID         int
	InstanceID int
	Name       string
	Value      float64
	Timestamp  time.Time
}

type Ahh struct {
	ID         int
	InstanceID int
	Name       string
	Value      float64
	Timestamp  time.Time
	Nuniga     string
}

var Tables map[string]any

func Init(dbURL string) error {
	logrus.Info("Connecting to DB...")
	ctx := context.Background()

	var err error
	Pool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to create DB pool")
		return err
	}

	// register tables
	Tables = map[string]any{
		"kkas":  RandomStruct{},
		"kkas2": Ahh{},
	}

	for name, model := range Tables {
		if err := CreateTableFromStruct(ctx, name, model); err != nil {
			return err
		}
	}

	return nil
}

// sql type lookup
func sqlType(t reflect.Type) (string, error) {
	if t == reflect.TypeOf(time.Time{}) {
		return "TIMESTAMPTZ", nil
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "BIGINT", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "BIGINT", nil
	case reflect.Float32, reflect.Float64:
		return "DOUBLE PRECISION", nil
	case reflect.Bool:
		return "BOOLEAN", nil
	case reflect.String:
		return "TEXT", nil
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "BYTEA", nil
		}
	}

	return "", fmt.Errorf("unsupported type: %s", t.String())
}

func toLowerCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + 32) // uppercase → lowercase in ASCII
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func ToAnySlice(slice any) []any {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil
	}
	result := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
}

func CreateTableFromStruct(ctx context.Context, tableName string, model any) error {
	t := reflect.TypeOf(model)

	// Unwrap pointer if needed (*Event → Event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("model must be a struct, got %s", t.Kind())
	}

	var columns []string
	hasPK := false

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")

		// Skip this field entirely
		if tag == "-" {
			continue
		}

		// Determine column name
		colName := toLowerCase(field.Name)
		if tag != "" && tag != "pk" {
			colName = tag // custom name from tag
		}

		// Get the Postgres type for this field
		sqlT, err := sqlType(field.Type)
		if err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}

		if tag == "pk" {
			columns = append(columns, fmt.Sprintf("%s %s PRIMARY KEY", colName, sqlT))
			hasPK = true
		} else {
			columns = append(columns, fmt.Sprintf("%s %s NOT NULL", colName, sqlT))
		}
	}

	if !hasPK {
		columns = append([]string{"id SERIAL PRIMARY KEY"}, columns...)
	}

	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n    %s\n);",
		tableName,
		strings.Join(columns, ",\n    "),
	)

	logrus.Infof("Running query:\n%s", query) // print the query

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to create table %q", tableName)
		return err
	}

	logrus.Infof("Table %q created successfully", tableName)

	return nil
}

func BatchInsertIntoTable(ctx context.Context, tableName string, models []any) error {
	if len(models) == 0 {
		return nil
	}

	t := reflect.TypeOf(models[0])
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var colNames []string
	var fieldIndexes []int

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")

		if tag == "-" || tag == "pk" {
			continue
		}

		colName := toLowerCase(field.Name)
		if tag != "" {
			colName = tag
		}

		colNames = append(colNames, colName)
		fieldIndexes = append(fieldIndexes, i)
	}

	placeholders := make([]string, len(colNames))
	for i := range colNames {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(colNames, ", "),
		strings.Join(placeholders, ", "),
	)

	batch := &pgx.Batch{}
	for _, model := range models {
		v := reflect.ValueOf(model)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		values := make([]any, len(fieldIndexes))
		for i, idx := range fieldIndexes {
			values[i] = v.Field(idx).Interface()
		}

		batch.Queue(query, values...)
	}

	br := Pool.SendBatch(ctx, batch)
	defer br.Close()

	ids := make([]int, len(models))
	for i := range models {
		if err := br.QueryRow().Scan(&ids[i]); err != nil {
			logrus.WithError(err).Errorf("Failed to scan row %d", i)
			return err
		}
	}

	logrus.Infof("Inserted %d rows into %q", len(models), tableName)

	return nil
}

func ReadFromTableWhere(ctx context.Context, tableName string, model any, colName string, colValues []any) ([]any, error) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var colNames []string
	var fieldIndexes []int

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")

		if tag == "-" {
			continue
		}

		col := toLowerCase(field.Name)
		if tag != "" && tag != "pk" {
			col = tag
		}

		colNames = append(colNames, col)
		fieldIndexes = append(fieldIndexes, i)
	}

	placeholders := make([]string, len(colValues))
	for i := range colValues {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s IN (%s)",
		strings.Join(colNames, ", "),
		tableName,
		colName,
		strings.Join(placeholders, ", "),
	)

	rows, err := Pool.Query(ctx, query, colValues...)
	if err != nil {
		//logrus.WithError(err).Errorf("Failed to query table %q", tableName)
		fmt.Println("Query failed:")
		fmt.Println(err)
		fmt.Errorf("Failed to query table %w", err)
		return nil, err
	}

	fmt.Printf("Running query:With values:")
	defer rows.Close()

	var results []any

	for rows.Next() {
		newStruct := reflect.New(t).Elem()

		valuePtrs := make([]any, len(fieldIndexes))
		for i, idx := range fieldIndexes {
			valuePtrs[i] = newStruct.Field(idx).Addr().Interface()
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			logrus.WithError(err).Error("Failed to scan row")
			return nil, err
		}

		results = append(results, newStruct.Interface())
	}

	return results, nil
}
