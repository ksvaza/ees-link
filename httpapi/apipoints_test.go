package httpapi_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/httpapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestAPIMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.TestPointTestAPI(w, req, ps)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestTestAPISuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.TestPointTestAPI(w, req, ps)

	require.Equal(t, http.StatusOK, w.Code)

	expected := "<h1>Hello, this is a test API!</h1>\n"
	assert.Equal(t, expected, w.Body.String())
}

func TestTestAPIWithBody(t *testing.T) {
	body := []byte("some data")
	req := httptest.NewRequest(http.MethodGet, "/api/test", bytes.NewReader(body))
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.TestPointTestAPI(w, req, ps)

	require.Equal(t, http.StatusOK, w.Code)
	expected := "<h1>Hello, this is a test API!</h1>\n"
	assert.Equal(t, expected, w.Body.String())
}

func TestPointReceiveFormMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/submit-form", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.TestPointReceiveForm(w, req, ps)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestPointReceiveFormSuccess(t *testing.T) {
	body := []byte(`{"name":"Andris Bērziņš","email":"tests@example.com","message":"Sveika, pasaule!"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/submit-form", bytes.NewReader(body))
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.TestPointReceiveForm(w, req, ps)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	expected := `{"status":"success"}`
	assert.Equal(t, expected, w.Body.String())
}

func TestPointReceiveFormWithEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/submit-form", io.NopCloser(bytes.NewReader([]byte{})))
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.TestPointReceiveForm(w, req, ps)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestPointReceiveFormWithNonJsonBody(t *testing.T) {
	body := []byte("plain text data")
	req := httptest.NewRequest(http.MethodPost, "/api/submit-form", bytes.NewReader(body))
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.TestPointReceiveForm(w, req, ps)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	expected := `{"status":"success"}`
	assert.Equal(t, expected, w.Body.String())
}

func TestPointRegisterMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.Handler(httpapi.PointRegister)(w, req, ps)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPointRegisterEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.Handler(httpapi.PointRegister)(w, req, ps)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestPointRegisterInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.Handler(httpapi.PointRegister)(w, req, ps)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestPointRegisterMissingRequiredFields(t *testing.T) {
	body := []byte(`{"fullName":"Jānis Bērziņš"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	httpapi.Handler(httpapi.PointRegister)(w, req, ps)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
