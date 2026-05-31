package models_test

import (
	"encoding/json"
	"testing"

	"github.com/ksvaza/ees-link/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCarTelemetryStructureMarshalling(t *testing.T) {
	telementryInput := []byte(`{"ID":1,"RSSI":-96,"Accel":{"X":0,"Y":0,"Z":0},"GPS":{"Lat":569601,"Lon":233432,"Spd":0,"Sat":6},"PSU":{"Uop":4200,"Iop":1,"Pop":4,"Uip":5842,"Wh":33368},"SYS":{"Vbat":4,"DC":1,"Err":0}}`)

	var telementryStructure models.CarTelemetry

	err := json.Unmarshal(telementryInput, &telementryStructure)
	require.Nil(t, err)

	var telemetryResult []byte

	telemetryResult, err = json.Marshal(telementryStructure)
	require.Nil(t, err)

	assert.Equal(t, telementryInput, telemetryResult)
}

func TestFullCarControlStructureMarshalling(t *testing.T) {
	controlInput := []byte(`{"ID":1,"PSU":{"I":2000,"U":500,"St":1}}`)

	var controlStructure models.CarControl

	err := json.Unmarshal(controlInput, &controlStructure)
	require.Nil(t, err)

	var controlResult []byte

	controlResult, err = json.Marshal(controlStructure)
	require.Nil(t, err)

	assert.Equal(t, controlInput, controlResult)
}

func TestPartialCarControlStructureMarshalling(t *testing.T) {
	controlInput := []byte(`{"ID":1,"PSU":{"I":1000}}`)

	var controlStructure models.CarControl

	err := json.Unmarshal(controlInput, &controlStructure)
	require.Nil(t, err)

	var controlResult []byte

	controlResult, err = json.Marshal(controlStructure)
	require.Nil(t, err)

	assert.Equal(t, controlInput, controlResult)
}
