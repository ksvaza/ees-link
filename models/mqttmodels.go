package models

import (
	"encoding/json"
	"time"
)

type MqttConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MqttMessage struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

type MqttLogEntry struct {
	Message    MqttMessage `json:"message"`
	ReceivedAt time.Time   `json:"received_at"`
}

// Sacensību loģikas datu struktūras

// Mašīnas(bloku) telemetrijas datu struktūras

type TelemetryAccel struct {
	X int16 `json:"X"`
	Y int16 `json:"Y"`
	Z int16 `json:"Z"`
}

type TelemetryGPS struct {
	Latitude  uint32 `json:"Lat"`
	Longitute uint32 `json:"Lon"`
	Speed     uint32 `json:"Spd"`
	SatC      uint32 `json:"Sat"`
}

type TelemetryPSU struct {
	VoltageOut uint16 `json:"Uop"`
	CurrentOut uint16 `json:"Iop"`
	PowerOut   uint16 `json:"Pop"`
	VoltageIn  uint16 `json:"Uip"`
	WattHours  uint16 `json:"Wh"`
}

type TelemetrySYS struct {
	VoltageBat   uint16 `json:"Vbat"`
	BatConnected int8   `json:"DC"`
	ErrorCode    uint16 `json:"Err"`
}

type CarTelemetry struct {
	ID        int            `json:"ID"`
	RSSI      int            `json:"RSSI"`
	AccelData TelemetryAccel `json:"Accel"`
	GPSData   TelemetryGPS   `json:"GPS"`
	PSUData   TelemetryPSU   `json:"PSU"`
	SYSData   TelemetrySYS   `json:"SYS"`
}

// Bloku konfigurācijas/kontroles datu struktūras

type ControlPSU struct {
	Current uint16 `json:"I,omitempty"`
	Voltage uint16 `json:"U,omitempty"`
	Status  uint16 `json:"St,omitempty"`
}

type CarControl struct {
	ID         int        `json:"ID,omitempty"`
	PSUControl ControlPSU `json:"PSU"`
}
