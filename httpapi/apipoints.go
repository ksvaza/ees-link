package httpapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Sveika pasaule tipa situācija
func PointTestAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("TestAPI called %+v", ps)

	errorHandler := func(err error, code int) {
		logrus.WithError(err).Error("Error")
		http.Error(w, err.Error(), code)
	}

	if r.Method != http.MethodGet {
		errorHandler(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	// Read the request body and return it as the response
	_, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Hello, this is a test API!</h1>\n"))
}

// tīri tests formu saņemšanai, vēlāk varētu būt noderīgi, lai saņemtu datus no frontend formas
func PointReceiveForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointReceiveForm called %+v", ps)

	fmt.Printf("Received form data: %+v", ps)

	errorHandler := func(err error, code int) {
		logrus.WithError(err).Error("Error")
		http.Error(w, err.Error(), code)
	}

	if r.Method != http.MethodPost {
		errorHandler(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logrus.Infof("Received form data: %s", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))

}
