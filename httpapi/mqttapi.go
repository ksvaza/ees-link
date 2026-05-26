package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/ksvaza/ees-link/myqtt"
	"github.com/pkg/errors"
)

func PointGetMQTTLogs(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	// pagaidām piekļuve ir visiem, vēlāk būs tikai visadministratoram piekļuve pilnajiem logiem.

	prelimit := ps.ByName("limit")
	if prelimit == "" {
		return nil, errors.New("missing limit")
	}

	limit, err := strconv.Atoi(prelimit)
	if err != nil {
		return nil, errors.Wrap(err, "invalid limit")
	}

	logs, err := realDB.GetMQTTLogs(r.Context(), limit) // Adjust limit as needed
	if err != nil {
		return nil, errors.Wrap(err, "failed to get MQTT logs")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         logs,
	}, nil

}

func PointPostMQTTMessage(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var message models.MqttMessage
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &message); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	err = myqtt.SendMQTTMessage(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send MQTT message")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "MQTT message sent successfully",
	}, nil
}
